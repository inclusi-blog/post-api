package utils

//go:generate mockgen -source=post_validator.go -destination=./../mocks/mock_post_validator.go -package=mocks

import (
	"context"
	"errors"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/configuration"
	"post-api/models"
	"post-api/models/db"
)

type PostValidator interface {
	ValidateAndGetReadTime(draft db.Draft, ctx context.Context) (models.MetaData, error)
}

type postValidator struct {
	configData *configuration.ConfigData
}

func (validator postValidator) ValidateAndGetReadTime(draft db.Draft, ctx context.Context) (models.MetaData, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostValidator").WithField("method", "ValidateAndGetReadTime")

	draftID := draft.DraftID
	id := draftID
	logger.Infof("Validating draft to publish for draft id %v", id)

	config := validator.configData.ContentReadTimeConfig

	var interests []interface{}
	err := draft.Interest.Unmarshal(&interests)

	if err != nil {
		logger.Errorf("Error occurred while validating draft interests for draft id %v .%v", id, err)
		return models.MetaData{}, err
	}

	var postWordsCount int
	var imageCount int
	var readTime int
	var extractedTagline string
	var titleString string
	var previewImage string
	wordCountFetchErr := GetNumberOfWords(draft.PostData, &postWordsCount, ctx, &imageCount, &extractedTagline, &titleString, &previewImage)

	if wordCountFetchErr != nil {
		logger.Errorf("invalid post data for draft id %v .%v", id, wordCountFetchErr)
		return models.MetaData{}, wordCountFetchErr
	}

	readTime = CountContentReadTime(postWordsCount)
	CountImageReadTime(imageCount, &readTime)
	for _, value := range interests {
		interest := value.(map[string]interface{})
		if interest["name"] == "" {
			return models.MetaData{}, errors.New("interest is invalid")
		}
		configReadTime := config[interest["name"].(string)]
		if configReadTime != 0 {
			if readTime < configReadTime {
				logger.Errorf("post interest doesn't meet required read time %v .%v", draftID, readTime)
				return models.MetaData{}, errors.New("post interest doesn't meet required read time")
			}
			continue
		} else {
			if readTime < validator.configData.MinimumPostReadTime {
				logger.Errorf("post doesn't meet minimum read time %v .%v", draftID, readTime)
				return models.MetaData{}, errors.New("post doesn't meet minimum read time")
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
