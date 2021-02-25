package utils

import (
	"encoding/json"
	"github.com/gola-glitch/gola-utils/golaerror"
	"post-api/constants"
)

func GetGolaError(genericError error) *golaerror.Error {
	httpError, ok := genericError.(golaerror.HttpError)
	if !ok {
		return &constants.InternalServerError
	}
	golaError := &golaerror.Error{}
	err := json.Unmarshal(httpError.ResponseBody, golaError)
	if err != nil {
		return &constants.InternalServerError
	}
	return golaError
}
