package utils

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/story/models"
	"regexp"
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

func CountContentReadTime(contentWordsCount int) int {
	return int((0.0036 * float64(contentWordsCount)) * 60)
}

func GetNumberOfWords(content models.JSONString, wordsCount *int, ctx context.Context, imageCount *int, extractedTagline, titleString, previewImage *string) error {
	logger := logging.GetLogger(ctx).WithField("class", "StoryUtils").WithField("method", "GetNumberOfWords")
	var postData []interface{}
	err := content.Unmarshal(&postData)

	if err != nil {
		logger.Errorf("something went wrong %v", err)
		return err
	}

	for _, data := range postData {
		singleData := data.(map[string]interface{})
		value := singleData["type"]
		if value != nil {
			if value == "image" {
				if singleData["url"] != "" && *previewImage == "" {
					*previewImage = singleData["url"].(string)
				}
				*imageCount = *imageCount + 1
			}
		}
		singleChildren := singleData["children"].([]interface{})
		for _, childrenData := range singleChildren {
			data := childrenData.(map[string]interface{})
			textString := data["text"].(string)
			if textString != "" {
				individual := strings.Split(textString, " ")
				*wordsCount = len(individual) + *wordsCount
				if *titleString == "" {
					if len([]rune(textString)) > 100 {
						*titleString = string([]rune(textString)[:100])
						continue
					}
					*titleString = textString
				} else if *extractedTagline == "" {
					if len(textString) > 100 {
						tamilRunes := string([]rune(textString))
						*extractedTagline = tamilRunes[:100]
						continue
					}
					*extractedTagline = textString
				}
			}
		}
	}
	return nil
}

func GetTitleAndTaglineFromData(ctx context.Context, titleJson models.JSONString) (string, string, error) {
	logger := logging.GetLogger(ctx).WithField("class", "StoryUtils").WithField("method", "GetNumberOfWords")
	var postData []interface{}
	err := titleJson.Unmarshal(&postData)

	if err != nil {
		logger.Errorf("Error occurred while unmarshalling title text from slate json %v", err)
		return "", "", err
	}

	var tagline string
	var titleString string
	for _, data := range postData {
		if tagline != "" && titleString != "" {
			break
		}
		singleData := data.(map[string]interface{})
		singleChildren := singleData["children"].([]interface{})
		for _, childrenData := range singleChildren {
			if tagline != "" && titleString != "" {
				break
			}
			data := childrenData.(map[string]interface{})
			textString := data["text"].(string)
			logger.Info(textString)
			if titleString == "" {
				if len([]rune(textString)) > 100 {
					titleString = string([]rune(textString)[:100])
					continue
				}
				titleString = textString
				continue
			}
			if tagline == "" {
				if len([]rune(textString)) > 100 {
					tagline = string([]rune(textString)[:100])
					continue
				}
				tagline = textString
				continue
			}
		}
	}

	return titleString, tagline, nil
}

func GenerateUrl(titleString string) string {
	chars := []string{"]", "^", "\\\\", "[", ".", "(", ")", "!", "-", "@", "#", "%", "&", "*", "_", "+", "~", "`", "=", "{", "}", "\\", "/", "|", ",", ">", "<", "?", "$"}
	r := strings.Join(chars, "")
	re := regexp.MustCompile("[" + r + "]+")
	titleString = re.ReplaceAllString(titleString, "")
	lower := strings.ToLower(titleString)
	trimmedSpace := spaceFieldJoin(lower)
	return trimmedSpace
}

func spaceFieldJoin(str string) string {
	return strings.Join(strings.Fields(str), "-")
}
