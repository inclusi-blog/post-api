package utils

import (
	"context"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/assert"
	"post-api/story/models"
	"post-api/story/service/test_helper"
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
	titleString, tagline, err := GetTitleAndTaglineFromData(ctx, models.JSONString{
		JSONText: types.JSONText(test_helper.TitleTestData),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Install apps via helm in kubernetes", titleString)
	assert.Equal(t, "", tagline)
}

func TestGetTitleFromSlateJsonWhenTitleGreaterThan100Len(t *testing.T) {
	ctx := context.TODO()
	titleString, tagline, err := GetTitleAndTaglineFromData(ctx, models.JSONString{
		JSONText: types.JSONText(test_helper.TitleTestDataMoreThan100Len),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Install apps via helm in kubernetes Install apps via helm in kubernetes Install apps via helm in kub", titleString)
	assert.Equal(t, "", tagline)
}

func TestGetNumberOfWords(t *testing.T) {
	ctx := context.TODO()
	readTime := 0
	imageCount := 0
	extractedTagline := ""
	titleString := ""
	previewImage := ""
	err := GetNumberOfWords(models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}, &readTime, ctx, &imageCount, &extractedTagline, &titleString, &previewImage)
	assert.Equal(t, 711, readTime)
	assert.Nil(t, err)
}
