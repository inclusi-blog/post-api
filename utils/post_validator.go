package utils

//go:generate mockgen -source=post_validator.go -destination=./../mocks/mock_post_validator.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/configuration"
	"post-api/constants"
	"post-api/models"
	"post-api/models/db"
)

type PostValidator interface {
	ValidateAndGetMetaData(draft db.Draft, ctx context.Context) (models.MetaData, *golaerror.Error)
}

type postValidator struct {
	configData *configuration.ConfigData
}

func (validator postValidator) ValidateAndGetMetaData(draft db.Draft, ctx context.Context) (models.MetaData, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostValidator").WithField("method", "ValidateAndGetMetaData")

	draftID := draft.DraftID
	id := draftID
	logger.Infof("Validating draft to publish for draft id %v", id)

	configData := validator.configData
	contentReadTime := configData.ContentReadTimeConfig
	minimumPostReadTime := configData.MinimumPostReadTime

	var postWordsCount int
	var imageCount int
	var readTime int
	var extractedTagline string
	var titleString string
	var previewImage string
	wordCountFetchErr := GetNumberOfWords(draft.PostData, &postWordsCount, ctx, &imageCount, &extractedTagline, &titleString, &previewImage)

	if wordCountFetchErr != nil {
		logger.Errorf("invalid post data for draft id %v .%v", id, wordCountFetchErr)
		return models.MetaData{}, &constants.DraftValidationFailedError
	}

	readTime = CountContentReadTime(postWordsCount)
	CountImageReadTime(imageCount, &readTime)
	err := draft.IsValidInterest(ctx, contentReadTime, readTime, minimumPostReadTime)
	if err != nil {
		logger.Errorf("error occured while validating interest for draft %v, Error %v", draftID, err)
		return models.MetaData{}, err
	}

	logger.Infof("Successfully validated draft for id %v", draftID)

	return models.MetaData{
		Title:        titleString,
		Tagline:      extractedTagline,
		ReadTime:     readTime,
		PreviewImage: previewImage,
	}, nil
}

func NewPostValidator(data *configuration.ConfigData) PostValidator {
	return postValidator{
		configData: data,
	}
}
