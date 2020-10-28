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
	readTime := 0
	CountContentReadTime(500, &readTime)
	assert.Equal(t, 108, readTime)
}

func TestGetTitleFromSlateJson(t *testing.T) {
	ctx := context.TODO()
	titleString, err := GetTitleFromSlateJson(ctx, models.JSONString{
		JSONText: types.JSONText(test_helper.TitleTestData),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Install apps via helm in kubernetes", titleString)
}

func TestGetTitleFromSlateJsonWhenTitleGreaterThan100Len(t *testing.T) {
	ctx := context.TODO()
	titleString, err := GetTitleFromSlateJson(ctx, models.JSONString{
		JSONText: types.JSONText(test_helper.TitleTestDataMoreThan100Len),
	})
	assert.Nil(t, err)
	assert.Equal(t, "Install apps via helm in kubernetes Install apps via helm in kubernetes Install apps via helm in kub", titleString)
}
