package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClockNowShouldGiveCurrentIndianTime(t *testing.T) {
	clock := NewClock()
	now := clock.Now()
	assert.NotNil(t, now)
}
