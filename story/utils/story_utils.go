package utils

import (
	"context"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/mitchellh/mapstructure"
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
	var editor models.Editor
	err := content.Unmarshal(&editor)
	if err != nil {
		logger.Errorf("unable to marshal post %v", err)
		return err
	}

	for _, data := range editor.Blocks {
		value := data.Type
		if value.IsEqual(models.Image) {
			var image models.ImageElement
			err := mapstructure.Decode(data.Data, &image)
			if err != nil {
				logger.Errorf("unable to unmarshal image element %v", err)
				return err
			}
			if image.File.Url != "" && *previewImage == "" {
				*previewImage = image.File.Url
			}
			*imageCount = *imageCount + 1
			continue
		}

		if value.IsEqual(models.Header) || value.IsEqual(models.Paragraph) {
			text, err := data.GetText()
			if err != nil {
				logger.Errorf("unable to unmarshal %v data. Error %v", data.Type, err)
				return err
			}
			if text != "" {
				individual := strings.Split(text, " ")
				*wordsCount = len(individual) + *wordsCount
				if *titleString == "" {
					if len([]rune(text)) > 100 {
						*titleString = string([]rune(text)[:100])
						continue
					}
					*titleString = text
				} else if *extractedTagline == "" {
					if len(text) > 100 {
						tamilRunes := string([]rune(text))
						*extractedTagline = tamilRunes[:100]
						continue
					}
					*extractedTagline = text
				}
			}
		}
	}
	return nil
}

func GetTitleAndTaglineFromData(ctx context.Context, titleJson models.JSONString) (string, string, error) {
	logger := logging.GetLogger(ctx).WithField("class", "StoryUtils").WithField("method", "GetNumberOfWords")
	var editor models.Editor
	err := titleJson.Unmarshal(&editor)
	if err != nil {
		logger.Errorf("unable to marshal post %v", err)
		return "", "", err
	}

	var tagline string
	var titleString string
	for _, data := range editor.Blocks {
		if tagline != "" && titleString != "" {
			break
		}
		if data.Type.IsEqual(models.Header) || data.Type.IsEqual(models.Paragraph) {
			text, err := data.GetText()
			if err != nil {
				logger.Errorf("unable to unmarshal %v type. Error %v", data.Type, err)
				return "", "", err
			}
			logger.Info(text)
			if titleString == "" {
				if len([]rune(text)) > 100 {
					titleString = string([]rune(text)[:100])
					continue
				}
				titleString = text
				continue
			}
			if tagline == "" {
				if len([]rune(text)) > 100 {
					tagline = string([]rune(text)[:100])
					continue
				}
				tagline = text
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
