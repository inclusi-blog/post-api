package constants

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/golaerror"
)

const (
	PayloadValidationErrorCode string = "ERR_POST_PAYLOAD_INVALID"
	InternalServerErrorCode    string = "ERR_POST_INTERNAL_SERVER_ERROR"
	PostServiceFailureCode     string = "ERR_POST_SERVICE_FAILURE"
	NoInterestsFoundCode       string = "ERR_NO_INTERESTS_FOUND"
	NoDraftFoundCode           string = "ERR_NO_DRAFT_FOUND"
	ConnvertTitleToStringCode  string = "ERR_CONVERTING_TITLE_JSON_TO_STRING"
	DraftValidationFailedCode  string = "ERR_DRAFT_VALIDATION_FAILED"
)

var (
	PostServiceFailureError    = golaerror.Error{ErrorCode: PostServiceFailureCode, ErrorMessage: "Failed to communicate with post service"}
	PayloadValidationError     = golaerror.Error{ErrorCode: PayloadValidationErrorCode, ErrorMessage: "One or more of the request parameters are missing or invalid"}
	InternalServerError        = golaerror.Error{ErrorCode: InternalServerErrorCode, ErrorMessage: "something went wrong"}
	NoInterestsFoundError      = golaerror.Error{ErrorCode: NoInterestsFoundCode, ErrorMessage: "no interest tags found"}
	NoDraftFoundError          = golaerror.Error{ErrorCode: NoDraftFoundCode, ErrorMessage: "no draft found for the given draft id"}
	ConnvertTitleToStringError = golaerror.Error{ErrorCode: ConnvertTitleToStringCode, ErrorMessage: "Error Converting Title Json to String"}
	DraftValidationFailedError = golaerror.Error{ErrorCode: DraftValidationFailedCode, ErrorMessage: "some of the fields missing in draft"}
)

var ErrorCodeHttpStatusCodeMap = map[string]int{
	PayloadValidationErrorCode: http.StatusBadRequest,
	InternalServerErrorCode:    http.StatusInternalServerError,
	PostServiceFailureCode:     http.StatusInternalServerError,
	NoInterestsFoundCode:       http.StatusNotFound,
	NoDraftFoundCode:           http.StatusNotFound,
}

func GetGolaHttpCode(golaErrCode string) int {
	if httpCode, ok := ErrorCodeHttpStatusCodeMap[golaErrCode]; ok {
		return httpCode
	}
	return http.StatusInternalServerError
}

func StoryInternalServerError(message string) *golaerror.Error {
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
