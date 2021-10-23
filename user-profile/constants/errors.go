package constants

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/golaerror"
	"net/http"
)

const (
	InternalServerErrorCode    string = "ERR_PROFILE_INTERNAL_SERVER_ERROR"
	PayloadValidationErrorCode string = "ERR_PROFILE_PAYLOAD_INVALID"
	NoUserFoundErrorCode       string = "ERR_PROFILE_NO_USER_FOUND"
)

var (
	InternalServerError    = golaerror.Error{ErrorCode: InternalServerErrorCode, ErrorMessage: "something went wrong"}
	PayloadValidationError = golaerror.Error{ErrorCode: PayloadValidationErrorCode, ErrorMessage: "One or more of the request parameters are missing or invalid"}
	NoUserFoundError       = golaerror.Error{ErrorCode: NoUserFoundErrorCode, ErrorMessage: "no user found"}
)

var ErrorCodeHttpStatusCodeMap = map[string]int{
	InternalServerErrorCode:    http.StatusInternalServerError,
	PayloadValidationErrorCode: http.StatusBadRequest,
	NoUserFoundErrorCode:       http.StatusNotFound,
}

func GetGolaHttpCode(golaErrCode string) int {
	if httpCode, ok := ErrorCodeHttpStatusCodeMap[golaErrCode]; ok {
		return httpCode
	}
	return http.StatusInternalServerError
}

func UserProfileInternalServerError(message string) *golaerror.Error {
	return &golaerror.Error{
		ErrorCode:      InternalServerErrorCode,
		ErrorMessage:   "something went wrong",
		AdditionalData: message,
	}
}

func RespondWithGolaError(ctx *gin.Context, err error) {
	if golaErr, ok := err.(*golaerror.Error); ok {
		ctx.JSON(GetGolaHttpCode(golaErr.ErrorCode), golaErr)
		return
	}
	ctx.JSON(http.StatusInternalServerError, InternalServerError)
	return
}
