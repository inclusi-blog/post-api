package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDraftConvertInterests(t *testing.T) {
	value := "{Culture,Technology}"
	draft := Draft{
		Interests: &value,
	}
	draft.ConvertInterests()
	assert.Len(t, draft.InterestTags, 2)
	assert.Equal(t, []string{"Culture", "Technology"},draft.InterestTags)
}
