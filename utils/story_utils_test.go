package utils

import (
	"context"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
	"post-api/models"
	"post-api/service/test_helper"
	"testing"
)

func TestCountImageReadTime(t *testing.T) {
	readTime := 0
	CountImageReadTime(50, &readTime)
	assert.Equal(t, 195, readTime)
}

func TestCountImageReadTimeWhenImageCountLessThen10(t *testing.T) {
	readTime := 0
	CountImageReadTime(2, &readTime)
	assert.Equal(t, 23, readTime)
}

func TestCountContentReadTime(t *testing.T) {
	readTime := CountContentReadTime(500)
	assert.Equal(t, 108, readTime)
}

func TestGetTitleFromSlateJson(t *testing.T) {
	ctx := context.TODO()
	titleString, tagline, err := GetTitleAndTaglineFromSlateJson(ctx, models.JSONString{
		JSONText: types.JSONText(test_helper.TitleTestData),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Install apps via helm in kubernetes", titleString)
	assert.Equal(t, "", tagline)
}

func TestGetTitleFromSlateJsonWhenTitleGreaterThan100Len(t *testing.T) {
	ctx := context.TODO()
	titleString, tagline, err := GetTitleAndTaglineFromSlateJson(ctx, models.JSONString{
		JSONText: types.JSONText(test_helper.TitleTestDataMoreThan100Len),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Install apps via helm in kubernetes Install apps via helm in kubernetes Install apps via helm in kub", titleString)
	assert.Equal(t, "", tagline)
}

func TestGetTitleFromSlateJsonWhenLargePostData(t *testing.T) {
	ctx := context.TODO()
	titleString, tagline, err := GetTitleAndTaglineFromSlateJson(ctx, models.JSONString{
		JSONText: types.JSONText(test_helper.LargeTextData),
	})
	assert.Nil(t, err)
	assert.Equal(t, "இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சராச", titleString)
	assert.Equal(t, "படிப்படியாக உயர்ந்த எண்ணிக்கை", tagline)
}

func TestGetNumberOfWords(t *testing.T) {
	ctx := context.TODO()
	readTime := 0
	imageCount := 0
	extractedTagline := ""
	titleString := ""
	previewImage := ""
	err := GetNumberOfWords(models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}, &readTime, ctx, &imageCount, &extractedTagline, &titleString, &previewImage)
	assert.Equal(t, 890, readTime)
	assert.Nil(t, err)
}

func TestGenerateUrl(t *testing.T) {
	sampleString := "this is my first post !@@#!@$ 321212 1212 121!@_!@!@!@ !@!@! @! @! @!@ $R@$%U &* &^(* )*( )*  ^#$ @ #?> <|}"
	url := GenerateUrl(sampleString)
	assert.Equal(t, "this-is-my-first-post-ru", url)
}

func TestGenerateUrlAnotherInvalidString(t *testing.T) {
	sampleString := "~!@ ~ !@# this @#$ is my first #$% $% $%^ post %^&* %^& which ^&* will&*( &* be&*(( published first(*&^ )*%^"
	url := GenerateUrl(sampleString)
	assert.Equal(t, "this-is-my-first-post-which-will-be-published-first", url)
}
