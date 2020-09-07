package constants

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/golaerror"
	"net/http"
)

const (
	PayloadValidationErrorCode string = "ERR_POST_PAYLOAD_INVALID"
	InternalServerErrorCode    string = "ERR_POST_INTERNAL_SERVER_ERROR"
	PostServiceFailureCode     string = "ERR_POST_SERVICE_FAILURE"
)

var (
	PostServiceFailureError = golaerror.Error{ErrorCode: PostServiceFailureCode, ErrorMessage: "Failed to communicate with post service"}
	PayloadValidationError  = golaerror.Error{ErrorCode: PayloadValidationErrorCode, ErrorMessage: "One or more of the request parameters are missing or invalid"}
	InternalServerError     = golaerror.Error{ErrorCode: InternalServerErrorCode, ErrorMessage: "something went wrong"}
)

var ErrorCodeHttpStatusCodeMap = map[string]int{
	PayloadValidationErrorCode: http.StatusBadRequest,
	InternalServerErrorCode:    http.StatusInternalServerError,
	PostServiceFailureCode:     http.StatusInternalServerError,
}

func GetGolaHttpCode(golaErrCode string) int {
	if httpCode, ok := ErrorCodeHttpStatusCodeMap[golaErrCode]; ok {
		return httpCode
	}
	return http.StatusInternalServerError
}

func RespondWithGolaError(ctx *gin.Context, err error) {
	if golaErr, ok := err.(golaerror.Error); ok {
		ctx.JSON(GetGolaHttpCode(golaErr.ErrorCode), golaErr)
		return
	}
	ctx.JSON(http.StatusInternalServerError, InternalServerError)
	return
}
