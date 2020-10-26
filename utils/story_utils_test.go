package utils

import (
	"github.com/stretchr/testify/assert"
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
