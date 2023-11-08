package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	model2 "github.com/inclusi-blog/gola-utils/model"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"post-api/idp/constants"
	"post-api/idp/mocks"
	"post-api/idp/models/oauth"
	"post-api/idp/models/response"
	util "post-api/idp/utils"
	"testing"
)

var userId = "userId"

type TokenControllerTestSuite struct {
	suite.Suite
	mockCtrl     *gomock.Controller
	recorder     *httptest.ResponseRecorder
	context      *gin.Context
	oauthService *mocks.MockOauthLoginHandler
}

func TestTokenControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TokenControllerTestSuite))
}

func (suite *TokenControllerTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.oauthService = mocks.NewMockOauthLoginHandler(suite.mockCtrl)
}

func (suite TokenControllerTestSuite) TestShouldClearCSRFCookiesWhenExchangeTokenIsCalled() {
	tokenController := NewTokenController(suite.oauthService, false)
	request := oauth.TokenExchangeRequest{RedirectUri: "redirect", ClientId: "clientid", CodeVerifier: "code_Verifier", Code: "Code"}
	idToken := model2.IdToken{UserId: userId}
	requestBody, _ := util.Encode(request)
	exchangeResponse := response.TokenExchangeResponse{IdToken: "random-id-token", AccessToken: "random-access-token", EncryptedIdToken: "JWE", ExpiresIn: 10}
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))
	suite.oauthService.EXPECT().ExchangeToken(suite.context, request).Return(exchangeResponse, idToken, nil)

	tokenController.ExchangeToken(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	cookies := suite.recorder.Result().Cookies()
	suite.Equal(3, len(cookies))
	suite.Equal("oauth2_authentication_csrf", cookies[0].Name)
	suite.Equal("", cookies[0].Value)
	suite.Equal("oauth2_authentication_session", cookies[1].Name)
	suite.Equal("", cookies[1].Value)
	suite.Equal("oauth2_consent_csrf", cookies[2].Name)
	suite.Equal("", cookies[1].Value)
	bodyBytes, _ := ioutil.ReadAll(suite.recorder.Body)
	actualResponse := response.TokenExchangeResponse{}
	json.Unmarshal(bodyBytes, &actualResponse)

	suite.Equal(exchangeResponse.EncryptedIdToken, actualResponse.EncryptedIdToken)
	suite.Equal(exchangeResponse.AccessToken, actualResponse.AccessToken)
	suite.Equal(exchangeResponse.ExpiresIn, actualResponse.ExpiresIn)
	suite.Equal("dummy.jwt.value", actualResponse.IdToken)
}

func (suite TokenControllerTestSuite) TestSendIdTokenIfInsecureCookiesIsAllowed() {
	tokenController := NewTokenController(suite.oauthService, true)
	request := oauth.TokenExchangeRequest{RedirectUri: "redirect", ClientId: "clientid", CodeVerifier: "code_Verifier", Code: "Code"}
	idToken := model2.IdToken{UserId: userId}
	requestBody, _ := util.Encode(request)
	exchangeResponse := response.TokenExchangeResponse{IdToken: "random-id-token", AccessToken: "random-access-token", EncryptedIdToken: "JWE", ExpiresIn: 10}
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))
	suite.oauthService.EXPECT().ExchangeToken(suite.context, request).Return(exchangeResponse, idToken, nil)

	tokenController.ExchangeToken(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	cookies := suite.recorder.Result().Cookies()
	suite.Equal(3, len(cookies))
	suite.Equal("oauth2_authentication_csrf", cookies[0].Name)
	suite.Equal("", cookies[0].Value)
	suite.Equal("oauth2_authentication_session", cookies[1].Name)
	suite.Equal("", cookies[1].Value)
	suite.Equal("oauth2_consent_csrf", cookies[2].Name)
	suite.Equal("", cookies[1].Value)
	bodyBytes, _ := ioutil.ReadAll(suite.recorder.Body)
	actualResponse := response.TokenExchangeResponse{}
	json.Unmarshal(bodyBytes, &actualResponse)

	suite.Equal(exchangeResponse.EncryptedIdToken, actualResponse.EncryptedIdToken)
	suite.Equal(exchangeResponse.AccessToken, actualResponse.AccessToken)
	suite.Equal(exchangeResponse.ExpiresIn, actualResponse.ExpiresIn)
	suite.Equal(exchangeResponse.IdToken, actualResponse.IdToken)
}

func (suite TokenControllerTestSuite) TestShouldGiveBadRequestWhenRedirectURIIsNotPassed() {
	request := `{"code":"code","redirect_uri":"","client_id":"","code_verifier":"code-verifier"}`
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(request))
	tokenController := NewTokenController(suite.oauthService, false)

	tokenController.ExchangeToken(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite TokenControllerTestSuite) TestShouldReturnErrorReturnedByOauthService() {
	request := oauth.TokenExchangeRequest{Code: "code", RedirectUri: "redirecturi", ClientId: "client_id", CodeVerifier: "code_verifier"}
	requestBody, _ := util.Encode(request)
	exchangeResponse := response.TokenExchangeResponse{}
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))
	suite.oauthService.EXPECT().ExchangeToken(suite.context, request).Return(exchangeResponse, model2.IdToken{}, &constants.PayloadValidationError)
	tokenController := NewTokenController(suite.oauthService, false)

	tokenController.ExchangeToken(suite.context)
}

func (suite TokenControllerTestSuite) TestShouldReturnSeverErrorWhenServiceReturnedErrorWithTypeInternalServer() {
	request := oauth.TokenExchangeRequest{Code: "code", RedirectUri: "redirecturi", ClientId: "client_id", CodeVerifier: "code_verifier"}
	requestBody, _ := util.Encode(request)
	suite.context.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(requestBody))

	suite.oauthService.EXPECT().ExchangeToken(suite.context, request).Return(response.TokenExchangeResponse{}, model2.IdToken{}, &constants.InternalServerError)
	tokenController := NewTokenController(suite.oauthService, false)

	tokenController.ExchangeToken(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}
