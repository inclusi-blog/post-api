package utils

//go:generate mockgen -source=post_validator.go -destination=./../mocks/mock_post_validator.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/configuration"
	"post-api/story/constants"
	"post-api/story/models"
	"post-api/story/models/db"
)

type PostValidator interface {
	ValidateAndGetReadTime(draft db.Draft, ctx context.Context) (models.MetaData, *golaerror.Error)
}

type postValidator struct {
	configData *configuration.ConfigData
}

func (validator postValidator) ValidateAndGetReadTime(draft db.Draft, ctx context.Context) (models.MetaData, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostValidator").WithField("method", "ValidateAndGetReadTime")

	if len(draft.InterestTags) == 0 {
		logger.Error("interest is empty")
		return models.MetaData{}, &constants.DraftInterestParseError
	}

	draftID := draft.DraftID
	id := draftID
	logger.Infof("Validating draft to publish for draft id %v", id)

	config := validator.configData.ContentReadTimeConfig

	var postWordsCount int
	var imageCount int
	var readTime int
	var extractedTagline string
	var titleString string
	var previewImage string
	wordCountFetchErr := GetNumberOfWords(draft.Data, &postWordsCount, ctx, &imageCount, &extractedTagline, &titleString, &previewImage)

	if wordCountFetchErr != nil {
		logger.Errorf("invalid post data for draft id %v .%v", id, wordCountFetchErr)
		return models.MetaData{}, &constants.DraftValidationFailedError
	}

	readTime = CountContentReadTime(postWordsCount)
	CountImageReadTime(imageCount, &readTime)
	for _, value := range draft.InterestTags {
		configReadTime := config[value]
		if configReadTime != 0 {
			if readTime < configReadTime {
				logger.Errorf("post interest doesn't meet required read time %v .%v", draftID, readTime)
				return models.MetaData{}, &constants.InterestReadTimeDoesNotMeetErr
			}
			continue
		} else {
			if readTime < validator.configData.MinimumPostReadTime {
				logger.Errorf("post doesn't meet minimum read time %v .%v", draftID, readTime)
				return models.MetaData{}, &constants.ReadTimeNotMeetError
			}
		}
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
