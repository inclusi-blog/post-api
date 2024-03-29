package constants

import (
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"net/http"
)

const (
	PayloadValidationErrorCode     string = "ERR_IDP_PAYLOAD_INVALID"
	InternalServerErrorCode        string = "ERR_IDP_INTERNAL_SERVER_ERROR"
	IDPServiceFailureCode          string = "ERR_IDP_SERVICE_FAILURE"
	UserAlreadyExistsCode          string = "ERR_IDP_USER_ALREADY_EXISTS"
	RetryRegistrationCode          string = "ERR_IDP_RETRY_REGISTRATION"
	ActivationLinkExpiredCode      string = "ERR_IDP_ACTIVATION_LINK_EXPIRED"
	UnauthorisedRequestCode        string = "ERR_IDP_UNAUTHORIZED"
	UserNotFoundCode               string = "ERR_IDP_USER_NOT_FOUND"
	InvalidCredentialsCode         string = "ERR_IDP_INVALID_CREDENTIALS"
	InvalidLoginChallengeCode      string = "ERR_IDP_INVALID_LOGIN_CHALLENGE"
	InvalidConsentChallengeCode    string = "ERR_IDP_INVALID_CONSENT_CHALLENGE"
	UsernameUpdateErrorCode        string = "ERR_IDP_USER_USERNAME_UPDATE"
	NameUpdateErrorCode            string = "ERR_IDP_USER_NAME_UPDATE"
	AboutUpdateErrorCode           string = "ERR_IDP_USER_ABOUT_UPDATE"
	UsernameAlreadyPresentCode     string = "ERR_USERNAME_ALREADY_PRESENT"
	SocialURLUpdateErrorCode       string = "ERR_USER_PROFILE_SOCIAL_URL_UPDATE"
	UnableToAssignPreSignURL       string = "ERR_USER_PROFILE_UNABLE_TO_ASSIGN_PRESIGN"
	ObjectNotFoundErrorCode        string = "ERR_USER_PROFILE_OBJECT_NOT_FOUND"
	UnableToFetchObjectErrorCode   string = "ERR_USER_PROFILE_UNABLE_TO_FETCH_OBJECT"
	UnableToUpdateAvatarErrorCode  string = "ERR_IDP_UNABLE_TO_UPDATE_AVATAR"
	UnableToResetPasswordErrorCode string = "ERR_IDP_UNABLE_TO_RESET_PASSWORD"
)

var (
	IDPServiceFailureError           = golaerror.Error{ErrorCode: IDPServiceFailureCode, ErrorMessage: "Failed to communicate with idp service"}
	PayloadValidationError           = golaerror.Error{ErrorCode: PayloadValidationErrorCode, ErrorMessage: "One or more of the request parameters are missing or invalid"}
	InternalServerError              = golaerror.Error{ErrorCode: InternalServerErrorCode, ErrorMessage: "something went wrong"}
	RegistrationRetryError           = golaerror.Error{ErrorCode: RetryRegistrationCode, ErrorMessage: "Please retry again", AdditionalData: nil}
	UnableToProcessRegistrationError = golaerror.Error{ErrorCode: IDPServiceFailureCode, ErrorMessage: "Please try again later", AdditionalData: nil}
	ActivationLinkExpiredError       = golaerror.Error{ErrorCode: ActivationLinkExpiredCode, ErrorMessage: "Please try again or retry registration process", AdditionalData: nil}
	UserNotFoundError                = golaerror.Error{ErrorCode: UserNotFoundCode, ErrorMessage: "User not found"}
	InvalidCredentialsError          = golaerror.Error{ErrorCode: InvalidCredentialsCode, ErrorMessage: "invalid username or password"}
	InvalidLoginChallengeError       = golaerror.Error{ErrorCode: InvalidLoginChallengeCode, ErrorMessage: "Invalid login challenge"}
	InvalidConsentChallengeError     = golaerror.Error{ErrorCode: InvalidConsentChallengeCode, ErrorMessage: "invalid consent challenge code"}
	UsernameUpdateError              = golaerror.Error{ErrorCode: UsernameUpdateErrorCode, ErrorMessage: "unable to update username"}
	NameUpdateError                  = golaerror.Error{ErrorCode: NameUpdateErrorCode, ErrorMessage: "unable to update name"}
	AboutUpdateError                 = golaerror.Error{ErrorCode: AboutUpdateErrorCode, ErrorMessage: "unable to update about"}
	UsernameAlreadyPresentError      = golaerror.Error{ErrorCode: UsernameUpdateErrorCode, ErrorMessage: "username already available"}
	SocialUpdateError                = golaerror.Error{ErrorCode: SocialURLUpdateErrorCode, ErrorMessage: "unable to update social url"}
	UnableToAssignPreSignURLError    = golaerror.Error{ErrorCode: UnableToAssignPreSignURL, ErrorMessage: "unable to assign presign image url"}
	UnableToFetchObjectError         = golaerror.Error{ErrorCode: UnableToFetchObjectErrorCode, ErrorMessage: "unable to fetch object"}
	ObjectNotFoundError              = golaerror.Error{ErrorCode: ObjectNotFoundErrorCode, ErrorMessage: "image object not found"}
	UnableToUpdateAvatarError        = golaerror.Error{ErrorCode: UnableToUpdateAvatarErrorCode, ErrorMessage: "unable to upload avatar"}
	UnauthorisedRequestError         = golaerror.Error{ErrorCode: UnauthorisedRequestCode, ErrorMessage: "unauthorized request"}
	UnableToResetPasswordError       = golaerror.Error{ErrorCode: UnableToResetPasswordErrorCode, ErrorMessage: "unable to reset user password"}
)

var ErrorCodeHttpStatusCodeMap = map[string]int{
	PayloadValidationErrorCode:     http.StatusBadRequest,
	InternalServerErrorCode:        http.StatusInternalServerError,
	IDPServiceFailureCode:          http.StatusInternalServerError,
	UserAlreadyExistsCode:          http.StatusFound,
	RetryRegistrationCode:          http.StatusInternalServerError,
	ActivationLinkExpiredCode:      http.StatusUnauthorized,
	UserNotFoundCode:               http.StatusNotFound,
	InvalidCredentialsCode:         http.StatusUnauthorized,
	InvalidLoginChallengeCode:      http.StatusUnauthorized,
	InvalidConsentChallengeCode:    http.StatusUnauthorized,
	UsernameUpdateErrorCode:        http.StatusInternalServerError,
	NameUpdateErrorCode:            http.StatusInternalServerError,
	UsernameAlreadyPresentCode:     http.StatusConflict,
	ObjectNotFoundErrorCode:        http.StatusNotFound,
	UnauthorisedRequestCode:        http.StatusUnauthorized,
	UnableToResetPasswordErrorCode: http.StatusInternalServerError,
}

func GetGolaHttpCode(golaErrCode string) int {
	if httpCode, ok := ErrorCodeHttpStatusCodeMap[golaErrCode]; ok {
		return httpCode
	}
	return http.StatusInternalServerError
}

func RespondWithGolaError(ctx *gin.Context, err error) {
	if golaErr, ok := err.(*golaerror.Error); ok {
		ctx.JSON(GetGolaHttpCode(golaErr.ErrorCode), golaErr)
		return
	}
	ctx.JSON(http.StatusInternalServerError, InternalServerError)
	return
}
