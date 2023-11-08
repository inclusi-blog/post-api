package login

// mockgen -source=handlers/login/oauth_login_handler.go -destination=mocks/mock_oauth_login_handler.go -package=mocks
import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/http/request"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/inclusi-blog/gola-utils/mask_util"
	"github.com/inclusi-blog/gola-utils/model"
	oauth2 "github.com/inclusi-blog/gola-utils/oauth"
	"net/http"
	"post-api/configuration"
	"post-api/idp/constants"
	"post-api/idp/models/db"
	"post-api/idp/models/oauth"
	"post-api/idp/models/response"
	util "post-api/idp/utils"
	"time"
)

type OauthLoginHandler interface {
	AcceptLogin(ctx *gin.Context, loginChallenge string, profile db.UserProfile) (response.AcceptResponse, *golaerror.Error)
	AcceptConsentRequest(ctx *gin.Context, consentChallenge string) (interface{}, *golaerror.Error)
	ExchangeToken(ctx *gin.Context, exchangeRequest oauth.TokenExchangeRequest) (response.TokenExchangeResponse, model.IdToken, *golaerror.Error)
}

type oauthLoginHandler struct {
	httpRequestBuilder request.HttpRequestBuilder
	configData         *configuration.ConfigData
	oauthUtils         oauth2.Utils
	clock              util.Clock
}

func (authHandler oauthLoginHandler) AcceptLogin(ctx *gin.Context, loginChallenge string, profile db.UserProfile) (response.AcceptResponse, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "OauthLoginHandler").WithField("method", "AcceptLogin")
	maskEmail := mask_util.MaskEmail(ctx, profile.Email)
	logger.Infof("User email %v is trying login challenge", maskEmail)

	logger.Infof("Making Accept login request api call for user email %v", maskEmail)

	acceptResponse := response.AcceptResponse{}

	acceptLoginRequest := oauth.LoginAcceptRequest{
		Subject:     profile.Id.String(),
		Remember:    false,
		RememberFor: 0,
		Acr:         "1",
		Userprofile: profile,
	}

	queryParas := make(map[string]string)
	queryParas["login_challenge"] = loginChallenge

	auth := authHandler.configData.Oauth
	acceptLoginRequestEndpoint := auth.AdminUrl + auth.AcceptLoginRequestUrl
	httpError := authHandler.httpRequestBuilder.NewRequest().
		WithContext(ctx).
		WithJSONBody(acceptLoginRequest).
		AddQueryParameters(queryParas).
		ResponseAs(&acceptResponse).
		Put(acceptLoginRequestEndpoint)

	if httpError == nil {
		logger.Infof("Accept login request api call completed successfully for user email %v", maskEmail)
		return acceptResponse, nil
	}

	logger.Errorf("Accept login request api failed for user email %v .%v", maskEmail, httpError)

	genericAuthError := response.GenericAuthError{}
	unMarshalError := json.Unmarshal(httpError.(golaerror.HttpError).ResponseBody, &genericAuthError)

	if unMarshalError != nil {
		logger.Errorf("Error in unmarshalling generic error %v", unMarshalError)
		return acceptResponse, &constants.InternalServerError
	}

	logger.Errorf("Accept login request api failed with the hydra genericError %v", genericAuthError)

	if genericAuthError.StatusCode == http.StatusUnauthorized ||
		genericAuthError.StatusCode == http.StatusNotFound ||
		genericAuthError.StatusCode == http.StatusInternalServerError {
		return acceptResponse, &constants.InvalidLoginChallengeError
	}

	return acceptResponse, &constants.InternalServerError
}

func (authHandler oauthLoginHandler) AcceptConsentRequest(ctx *gin.Context, consentChallenge string) (interface{}, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "OauthLoginHandler").WithField("method", "AcceptConsentRequest")
	logger.Infof("Initiating API call to get consent.")

	queryParams := make(map[string]string)
	queryParams["consent_challenge"] = consentChallenge

	consentResponse := response.ConsentResponse{}

	auth := authHandler.configData.Oauth
	getConsentUrl := auth.AdminUrl + auth.GetConsentRequestUrl

	httpError := authHandler.httpRequestBuilder.NewRequest().
		WithContext(ctx).
		AddQueryParameters(queryParams).
		ResponseAs(&consentResponse).
		Get(getConsentUrl)

	if httpError != nil {
		logger.Errorf("Error in call to get consent response %v", httpError)
		return nil, &constants.InternalServerError
	}

	maskedEmail := mask_util.MaskEmail(ctx, consentResponse.Userprofile.Email)
	logger.Info("Get consent call finished successfully. Initiating call to accept consent for the subject ", maskedEmail)

	acceptConsentRequest := oauth.ConsentAcceptRequest{
		GrantAccessTokenAudience: consentResponse.RequestedAccessTokenAudience,
		GrantScope:               consentResponse.RequestedScope,
		Remember:                 true,
		RememberFor:              1,
		Session: oauth.Session{
			Userprofile: consentResponse.Userprofile,
		},
	}

	acceptResponse := response.AcceptResponse{}

	acceptConsentUrl := auth.AdminUrl + auth.AcceptConsentRequestUrl

	httpAcceptError := authHandler.httpRequestBuilder.
		NewRequest().WithContext(ctx).
		WithJSONBody(acceptConsentRequest).
		AddQueryParameters(queryParams).
		ResponseAs(&acceptResponse).
		Put(acceptConsentUrl)

	if httpAcceptError == nil {
		logger.Info("Accept consent call finished successfully for the subject ", maskedEmail)
		return acceptResponse, nil
	}

	logger.Errorf("Accept login request api failed for user email %v", maskedEmail)

	genericAuthError := response.GenericAuthError{}
	unMarshalError := json.Unmarshal(httpAcceptError.(golaerror.HttpError).ResponseBody, &genericAuthError)

	if unMarshalError != nil {
		logger.Errorf("Error in unmarshalling generic error %v", unMarshalError)
		return acceptResponse, &constants.InternalServerError
	}

	logger.Errorf("Accept consent request api failed with the hydra genericError %v", genericAuthError)

	if genericAuthError.StatusCode == http.StatusNotFound || genericAuthError.StatusCode == http.StatusInternalServerError {
		return acceptResponse, &constants.InvalidConsentChallengeError
	}

	return acceptResponse, &constants.InternalServerError
}

func (authHandler oauthLoginHandler) ExchangeToken(ctx *gin.Context, request oauth.TokenExchangeRequest) (response.TokenExchangeResponse, model.IdToken, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "OauthLoginHandler").WithField("method", "ExchangeToken")
	formData := generateFormData(request)

	auth := authHandler.configData.Oauth

	getTokenUrl := auth.PublicUrl + auth.GetTokenUrl

	logger.Info("Making call to hydra to exchange token")
	tokenExchangeResponse := response.TokenExchangeResponse{}
	httpError := authHandler.httpRequestBuilder.NewRequest().
		WithContext(ctx).
		WithFormURLEncoded(formData).
		ResponseAs(&tokenExchangeResponse).
		Post(getTokenUrl)

	if httpError != nil {
		logger.Errorf("Error in call to exchange token. Error: %-v", httpError)
		return handleTokenExchangeError(ctx, httpError)
	}

	tokenExchangeResponse.ExpiresAt = authHandler.clock.Now().UTC().Add(time.Second * time.Duration(tokenExchangeResponse.ExpiresIn)).
		Format(http.TimeFormat)

	logger.Info("Decoding id token from JWT")
	token, decodeError := authHandler.oauthUtils.DecodeIdTokenFromJWT(tokenExchangeResponse.IdToken)
	if decodeError != nil {
		logger.Errorf("Error in decoding id token from JWT. Error %-v", decodeError)
		return response.TokenExchangeResponse{}, model.IdToken{}, &constants.InternalServerError
	}

	logger.Info("Encrypting id token")
	jweToken, jweEncryptionError := authHandler.oauthUtils.EncryptIdToken(ctx, tokenExchangeResponse.IdToken)
	if jweEncryptionError != nil {
		logger.Errorf("Error in encrypting id token. Error %-v", jweEncryptionError)
		return response.TokenExchangeResponse{}, model.IdToken{}, &constants.InternalServerError
	}

	logger.Info("Exchange token call completed successfully")

	tokenExchangeResponse.EncryptedIdToken = jweToken
	return tokenExchangeResponse, token, nil
}

func handleTokenExchangeError(ctx *gin.Context, httpError error) (response.TokenExchangeResponse, model.IdToken, *golaerror.Error) {
	genericAuthError := response.GenericAuthError{}
	unmarshalError := json.Unmarshal(httpError.(golaerror.HttpError).ResponseBody, &genericAuthError)
	logger := logging.GetLogger(ctx).WithField("class", "OauthLoginHandler").WithField("method", "handleTokenExchangeError")
	if unmarshalError != nil {
		logger.Error("Error in unmarshalling Generic Error Obj", unmarshalError.Error())
		return response.TokenExchangeResponse{}, model.IdToken{}, &constants.InternalServerError
	}

	if genericAuthError.StatusCode == http.StatusUnauthorized {
		logger.Error("empty access token in response")
		return response.TokenExchangeResponse{}, model.IdToken{}, &constants.PayloadValidationError
	}

	logger.Error("error in call to exchange token", httpError.Error())
	return response.TokenExchangeResponse{}, model.IdToken{}, &constants.InternalServerError
}

func generateFormData(request oauth.TokenExchangeRequest) map[string]interface{} {
	formData := make(map[string]interface{})
	formData["grant_type"] = "authorization_code"
	formData["code_verifier"] = request.CodeVerifier
	formData["client_id"] = request.ClientId
	formData["redirect_uri"] = request.RedirectUri
	formData["code"] = request.Code
	return formData
}

func NewOauthLoginHandler(builder request.HttpRequestBuilder, data *configuration.ConfigData, utils oauth2.Utils, clock util.Clock) OauthLoginHandler {
	return oauthLoginHandler{
		httpRequestBuilder: builder,
		configData:         data,
		oauthUtils:         utils,
		clock:              clock,
	}
}
