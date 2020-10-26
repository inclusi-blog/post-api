package utils

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/models"
	"strings"
)

func CountImageReadTime(imageCount int, readTime *int) {
	imageReadStartsAt := 12
	if imageCount > 10 {
		for i := imageCount; i > 0; i-- {
			if i < 10 {
				*readTime = *readTime + 3
				continue
			}
			if imageReadStartsAt < 3 {
				*readTime = *readTime + 3
				continue
			}
			*readTime = *readTime + imageReadStartsAt
			imageReadStartsAt--
		}
		return
	}
	for i := imageCount; i > 0; i-- {
		*readTime = *readTime + imageReadStartsAt
		imageReadStartsAt--
	}
}

func CountContentReadTime(contentWordsCount int, readTime *int) {
	*readTime = *readTime + int((0.0036*float64(contentWordsCount))*60)
}

func GetNumberOfWords(content models.JSONString, wordsCount *int, ctx context.Context, imageCount *int, extractedTagline *string) error {
	logger := logging.GetLogger(ctx).WithField("class", "StoryUtils").WithField("method", "GetNumberOfWords")
	var postData []interface{}
	err := content.Unmarshal(&postData)

	if err != nil {
		logger.Errorf("something went wrong %v", err)
		return err
	}

	for topIndex, data := range postData {
		singleData := data.(map[string]interface{})
		value := singleData["type"]
		if value != nil {
			if value == "image" {
				*imageCount = *imageCount + 1
			}
		}
		singleChildren := singleData["children"].([]interface{})
		for innerIndex, childrenData := range singleChildren {
			data := childrenData.(map[string]interface{})
			textString := data["text"].(string)
			if textString != "" {
				individual := strings.Split(textString, " ")
				*wordsCount = len(individual) + *wordsCount
				if topIndex == 0 && innerIndex == 0 {
					if len(textString) > 100 {
						*extractedTagline = string([]rune(textString)[:100])
						continue
					}
					*extractedTagline = textString
				}
			}
		}
	}
	return nil
}
