package utils

import (
	"errors"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/stretchr/testify/assert"
	"post-api/constants"
	"post-api/service/test_helper"
	"testing"
)

func TestGetGolaError(t *testing.T) {
	golaError := `{"errorCode":"ERR_NO_INTERESTS_FOLLOWED","errorMessage": "user followed no interests"}`

	httpError := golaerror.HttpError{
		StatusCode:   400,
		ResponseBody: []byte(golaError),
	}
	actualGolaErr := GetGolaError(httpError)
	assert.Equal(t, "ERR_NO_INTERESTS_FOLLOWED", actualGolaErr.ErrorCode)
}

func TestGetGolaError_WhenDifferentError(t *testing.T) {
	actualGolaErr := GetGolaError(errors.New(test_helper.ErrSomethingWentWrong))
	assert.Equal(t, &constants.InternalServerError, actualGolaErr)
}

func TestGetGolaError_WhenDifferentErrorOtherThanHTTPError(t *testing.T) {
	golaError := `{"error":"ERR_NO_INTERESTS_FOLLOWED","message": "user followed no interests"}`

	httpError := golaerror.HttpError{
		StatusCode:   400,
		ResponseBody: []byte(golaError),
	}
	actualGolaErr := GetGolaError(httpError)
	assert.Equal(t, &golaerror.Error{}, actualGolaErr)
}
