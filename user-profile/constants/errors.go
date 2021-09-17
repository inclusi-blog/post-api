package constants

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/golaerror"
	"net/http"
)

const (
	InternalServerErrorCode string = "ERR_POST_INTERNAL_SERVER_ERROR"
)

var (
	InternalServerError = golaerror.Error{ErrorCode: InternalServerErrorCode, ErrorMessage: "something went wrong"}
)

var ErrorCodeHttpStatusCodeMap = map[string]int{
	InternalServerErrorCode: http.StatusInternalServerError,
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
