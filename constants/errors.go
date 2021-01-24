package constants

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/golaerror"
)

const (
	PayloadValidationErrorCode       string = "ERR_POST_PAYLOAD_INVALID"
	InternalServerErrorCode          string = "ERR_POST_INTERNAL_SERVER_ERROR"
	PostServiceFailureCode           string = "ERR_POST_SERVICE_FAILURE"
	NoInterestsFoundCode             string = "ERR_NO_INTERESTS_FOUND"
	NoDraftFoundCode                 string = "ERR_NO_DRAFT_FOUND"
	ConvertTitleToStringCode         string = "ERR_CONVERTING_TITLE_JSON_TO_STRING"
	DraftValidationFailedCode        string = "ERR_DRAFT_VALIDATION_FAILED"
	InterestParseErrorCode           string = "ERR_DRAFT_INTEREST_PARSE_FAILED"
	ReadTimeNotMeetCode              string = "ERR_DRAFT_READ_TIME_NOT_MEET"
	InterestDoesNotMeetReadTimeCode  string = "ERR_DRAFT_INTEREST_NOT_MEET_READ_TIME"
	PostNotFoundCode                 string = "ERR_POST_NOT_FOUND"
	MinimumInterestCountNotMatchCode string = "ERR_MINIMUM_INTERESTS_COUNT_NOT_MATCH"
	NoPostFoundCode                  string = "ERR_NO_POST_FOUND"
	PostNotInReadLaterCode           string = "ERR_POST_NOT_IN_READ_LATER"
	NoInterestsAndCategoriesCode     string = "ERR_NO_INTERESTS_AND_CATEGORIES"
)

var (
	PostServiceFailureError         = golaerror.Error{ErrorCode: PostServiceFailureCode, ErrorMessage: "Failed to communicate with post service"}
	PayloadValidationError          = golaerror.Error{ErrorCode: PayloadValidationErrorCode, ErrorMessage: "One or more of the request parameters are missing or invalid"}
	InternalServerError             = golaerror.Error{ErrorCode: InternalServerErrorCode, ErrorMessage: "something went wrong"}
	NoInterestsFoundError           = golaerror.Error{ErrorCode: NoInterestsFoundCode, ErrorMessage: "no interest tags found"}
	NoDraftFoundError               = golaerror.Error{ErrorCode: NoDraftFoundCode, ErrorMessage: "no draft found for the given draft id"}
	ConvertTitleToStringError       = golaerror.Error{ErrorCode: ConvertTitleToStringCode, ErrorMessage: "Error Converting Title Json to String"}
	DraftValidationFailedError      = golaerror.Error{ErrorCode: DraftValidationFailedCode, ErrorMessage: "some of the fields missing in draft"}
	DraftInterestParseError         = golaerror.Error{ErrorCode: InterestParseErrorCode, ErrorMessage: "please reenter the interests", AdditionalData: "Please re enter the interest for draft"}
	ReadTimeNotMeetError            = golaerror.Error{ErrorCode: ReadTimeNotMeetCode, ErrorMessage: "read time requirement not meet", AdditionalData: "Please Enter some more content to the draft before publishing"}
	InterestReadTimeDoesNotMeetErr  = golaerror.Error{ErrorCode: InterestDoesNotMeetReadTimeCode, ErrorMessage: "selected interest doesn't meet required read time", AdditionalData: "Increase the content for the draft"}
	PostNotFoundErr                 = golaerror.Error{ErrorCode: PostNotFoundCode, ErrorMessage: "no post found for the given post uid"}
	MinimumInterestCountNotMatchErr = golaerror.Error{ErrorCode: MinimumInterestCountNotMatchCode, ErrorMessage: "interest tags requirement not met", AdditionalData: "Please select 3-5 interest tags"}
	PostNotInReadLaterErr           = golaerror.Error{ErrorCode: PostNotInReadLaterCode, ErrorMessage: "post currently not in read later"}
	NoInterestsAndCategoriesErr     = golaerror.Error{ErrorCode: NoInterestsAndCategoriesCode, ErrorMessage: "no categories and interests found"}
)

var ErrorCodeHttpStatusCodeMap = map[string]int{
	PayloadValidationErrorCode:       http.StatusBadRequest,
	InternalServerErrorCode:          http.StatusInternalServerError,
	PostServiceFailureCode:           http.StatusInternalServerError,
	NoInterestsFoundCode:             http.StatusNotFound,
	DraftValidationFailedCode:        http.StatusNotAcceptable,
	NoDraftFoundCode:                 http.StatusNotFound,
	ConvertTitleToStringCode:         http.StatusBadRequest,
	ReadTimeNotMeetCode:              http.StatusBadRequest,
	InterestParseErrorCode:           http.StatusBadRequest,
	InterestDoesNotMeetReadTimeCode:  http.StatusNotAcceptable,
	PostNotFoundCode:                 http.StatusNotFound,
	MinimumInterestCountNotMatchCode: http.StatusNotAcceptable,
	PostNotInReadLaterCode:           http.StatusNotAcceptable,
	NoInterestsAndCategoriesCode:     http.StatusNotFound,
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
