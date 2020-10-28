package utils

//go:generate mockgen -source=post_validator.go -destination=./../mocks/mock_post_validator.go -package=mocks

import (
	"context"
	"errors"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/configuration"
	"post-api/models/db"
)

type PostValidator interface {
	ValidateAndGetReadTime(draft *db.Draft, ctx context.Context) (string, int, error)
}

type postValidator struct {
	configData *configuration.ConfigData
}

func (validator postValidator) ValidateAndGetReadTime(draft *db.Draft, ctx context.Context) (string, int, error){
	logger := logging.GetLogger(ctx).WithField("class", "PostValidator").WithField("method", "ValidateAndGetReadTime")

	draftID := draft.DraftID
	id := draftID
	logger.Infof("Validating draft to publish for draft id %v", id)

	config := validator.configData.ContentReadTimeConfig

	var interests []interface{}
	err := draft.Interest.Unmarshal(&interests)

	if err != nil {
		logger.Errorf("Error occurred while validating draft interests for draft id %v .%v", id, err)
		return "", 0, err
	}

	var postWordsCount int
	var titleWordsCount int
	var imageCount int
	var readTime int
	var extractedTagline string
	var titleString string
	wordCountFetchErr := GetNumberOfWords(draft.PostData, &postWordsCount, ctx, &imageCount, &extractedTagline)

	if wordCountFetchErr != nil {
		logger.Errorf("invalid post data for draft id %v .%v", id, wordCountFetchErr)
		return "", 0, wordCountFetchErr
	}

	titleWordsCountFetchErr := GetNumberOfWords(draft.TitleData, &titleWordsCount, ctx, &imageCount, &titleString)

	if titleWordsCountFetchErr != nil {
		logger.Errorf("invalid post title data for draft id %v .%v", id, titleWordsCountFetchErr)
		return "", 0, titleWordsCountFetchErr
	}

	overAllWordsCount := postWordsCount + titleWordsCount
	CountContentReadTime(overAllWordsCount, &readTime)
	CountImageReadTime(imageCount, &readTime)
	for _, value := range interests {
		interest := value.(map[string]interface{})
		if interest["name"] == "" {
			return "", 0, errors.New("interest is invalid")
		}
		configReadTime := config[interest["name"].(string)]
		if configReadTime != 0 {
			if readTime < configReadTime {
				logger.Errorf("post interest doesn't meet required read time %v", draftID)
				return "", 0, errors.New("post interest doesn't meet required read time")
			}
			continue
		} else {
			if readTime < validator.configData.MinimumPostReadTime {
				logger.Errorf("post doesn't meet minimum read time %v", draftID)
				return "", 0, errors.New("post doesn't meet minimum read time")
			}
		}
	}

	if draft.Tagline == "" {
		logger.Infof("Tagline is empty setting some text from post data for post id %v", draftID)
		draft.Tagline = extractedTagline
	}

	logger.Infof("Successfully validated draft for id %v", draftID)

	return titleString, readTime, nil
}

func NewPostValidator(data *configuration.ConfigData) PostValidator {
	return postValidator{
		configData: data,
	}
}
