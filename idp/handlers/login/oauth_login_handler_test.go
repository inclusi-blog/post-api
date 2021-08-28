package login

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/http/request"
	"github.com/gola-glitch/gola-utils/http/request/mocks"
	utilMocks "github.com/gola-glitch/gola-utils/mocks"
	"github.com/gola-glitch/gola-utils/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/configuration"
	"post-api/idp/constants"
	mocks2 "post-api/idp/mocks"
	"post-api/idp/models/db"
	"post-api/idp/models/oauth"
	"post-api/idp/models/response"
	"testing"
	"time"
)

type OauthLoginHandlerTest struct {
	suite.Suite
	mockController  *gomock.Controller
	ginContext      *gin.Context
	configData      *configuration.ConfigData
	mockOauthUtils  *utilMocks.MockUtils
	mockClock       *mocks2.MockClock
	oauthHandler    OauthLoginHandler
	mockHttpRequest *mocks.MockHttpRequest

	mockHttpRequestBuilder *mocks.MockHttpRequestBuilder
}

func TestOauthLoginHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(OauthLoginHandlerTest))
}

func (suite *OauthLoginHandlerTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.ginContext, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.ginContext.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	suite.configData = &configuration.ConfigData{
		Oauth: configuration.OAuth{
			AdminUrl:                "https://oauth2/",
			PublicUrl:               "https://oauth2/",
			AcceptLoginRequestUrl:   "oauth2/auth/login/accept",
			GetConsentRequestUrl:    "oauth2/auth/requests/consent",
			AcceptConsentRequestUrl: "oauth2/auth/requests/consent/accept",
			GetTokenUrl:             "token-url",
		},
	}
	suite.mockOauthUtils = utilMocks.NewMockUtils(suite.mockController)
	suite.mockClock = mocks2.NewMockClock(suite.mockController)
	suite.mockHttpRequest = mocks.NewMockHttpRequest(suite.mockController)
	suite.mockHttpRequestBuilder = mocks.NewMockHttpRequestBuilder(suite.mockController)
	suite.oauthHandler = NewOauthLoginHandler(suite.mockHttpRequestBuilder, suite.configData, suite.mockOauthUtils, suite.mockClock)
}

func (suite *OauthLoginHandlerTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *OauthLoginHandlerTest) TestAcceptLogin_WhenHydraAcceptLoginSuccess() {
	loginChallenge := "some-login-challenge"

	profile := db.UserProfile{
		UserID:   "some-user-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptRequest := oauth.LoginAcceptRequest{
		Subject:     "some-user-id",
		Remember:    false,
		RememberFor: 0,
		Acr:         "1",
		Userprofile: profile,
	}

	queryParams := make(map[string]string)

	queryParams["login_challenge"] = loginChallenge

	acceptResponse := response.AcceptResponse{}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		tempResponsePointer := responseBuilder.(*response.AcceptResponse)
		*tempResponsePointer = response.AcceptResponse{
			RedirectTo: "http://redirect-to-url",
		}
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/login/accept").Return(nil)

	loginAcceptResponse, genericErr := suite.oauthHandler.AcceptLogin(suite.ginContext, loginChallenge, profile)

	suite.Nil(genericErr)
	suite.Equal(response.AcceptResponse{RedirectTo: "http://redirect-to-url"}, loginAcceptResponse)
}

func (suite *OauthLoginHandlerTest) TestAcceptLogin_WhenHydraReturnsGenericNotFoundError() {
	loginChallenge := "some-login-challenge"

	profile := db.UserProfile{
		UserID:   "some-user-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptRequest := oauth.LoginAcceptRequest{
		Subject:     "some-user-id",
		Remember:    false,
		RememberFor: 0,
		Acr:         "1",
		Userprofile: profile,
	}

	authError := response.GenericAuthError{
		Error:      "Not Found",
		StatusCode: 404,
	}

	bytes, _ := json.Marshal(authError)

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusNotFound,
		ResponseBody: bytes,
	}

	queryParams := make(map[string]string)

	queryParams["login_challenge"] = loginChallenge

	acceptResponse := response.AcceptResponse{}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/login/accept").Return(httpError)

	loginAcceptResponse, genericErr := suite.oauthHandler.AcceptLogin(suite.ginContext, loginChallenge, profile)

	suite.NotNil(genericErr)
	suite.Empty(loginAcceptResponse)
	suite.Equal(&constants.InvalidLoginChallengeError, genericErr)
}

func (suite *OauthLoginHandlerTest) TestAcceptLogin_WhenHydraReturnsGenericUnAuthorizedError() {
	loginChallenge := "some-login-challenge"

	profile := db.UserProfile{
		UserID:   "some-user-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptRequest := oauth.LoginAcceptRequest{
		Subject:     "some-user-id",
		Remember:    false,
		RememberFor: 0,
		Acr:         "1",
		Userprofile: profile,
	}

	authError := response.GenericAuthError{
		Error:      "Unauthorised error",
		StatusCode: 401,
	}

	bytes, _ := json.Marshal(authError)

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusUnauthorized,
		ResponseBody: bytes,
	}

	queryParams := make(map[string]string)

	queryParams["login_challenge"] = loginChallenge

	acceptResponse := response.AcceptResponse{}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/login/accept").Return(httpError)

	loginAcceptResponse, genericErr := suite.oauthHandler.AcceptLogin(suite.ginContext, loginChallenge, profile)

	suite.NotNil(genericErr)
	suite.Empty(loginAcceptResponse)
	suite.Equal(&constants.InvalidLoginChallengeError, genericErr)
}

func (suite *OauthLoginHandlerTest) TestAcceptLogin_WhenHydraReturnsGenericInternalServerError() {
	loginChallenge := "some-login-challenge"

	profile := db.UserProfile{
		UserID:   "some-user-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptRequest := oauth.LoginAcceptRequest{
		Subject:     "some-user-id",
		Remember:    false,
		RememberFor: 0,
		Acr:         "1",
		Userprofile: profile,
	}

	authError := response.GenericAuthError{
		Error:      "Internal server error",
		StatusCode: 500,
	}

	bytes, _ := json.Marshal(authError)

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusInternalServerError,
		ResponseBody: bytes,
	}

	queryParams := make(map[string]string)

	queryParams["login_challenge"] = loginChallenge

	acceptResponse := response.AcceptResponse{}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/login/accept").Return(httpError)

	loginAcceptResponse, genericErr := suite.oauthHandler.AcceptLogin(suite.ginContext, loginChallenge, profile)

	suite.NotNil(genericErr)
	suite.Empty(loginAcceptResponse)
	suite.Equal(&constants.InvalidLoginChallengeError, genericErr)
}

func (suite *OauthLoginHandlerTest) TestAcceptLogin_WhenUnableToMarshalGenericError() {
	loginChallenge := "some-login-challenge"

	profile := db.UserProfile{
		UserID:   "some-user-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptRequest := oauth.LoginAcceptRequest{
		Subject:     "some-user-id",
		Remember:    false,
		RememberFor: 0,
		Acr:         "1",
		Userprofile: profile,
	}

	invalidBytes, err := json.Marshal("hello")

	suite.Nil(err)

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusInternalServerError,
		ResponseBody: invalidBytes,
	}

	queryParams := make(map[string]string)

	queryParams["login_challenge"] = loginChallenge

	acceptResponse := response.AcceptResponse{}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/login/accept").Return(httpError)

	loginAcceptResponse, genericErr := suite.oauthHandler.AcceptLogin(suite.ginContext, loginChallenge, profile)

	suite.NotNil(genericErr)
	suite.Empty(loginAcceptResponse)
	suite.Equal(&constants.InternalServerError, genericErr)
}

func (suite *OauthLoginHandlerTest) TestAcceptLogin_WhenHydraReturnsGenericOtherErrors() {
	loginChallenge := "some-login-challenge"

	profile := db.UserProfile{
		UserID:   "some-user-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptRequest := oauth.LoginAcceptRequest{
		Subject:     "some-user-id",
		Remember:    false,
		RememberFor: 0,
		Acr:         "1",
		Userprofile: profile,
	}

	authError := response.GenericAuthError{
		Error:      "Conflict",
		StatusCode: 409,
	}

	bytes, _ := json.Marshal(authError)

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusConflict,
		ResponseBody: bytes,
	}

	queryParams := make(map[string]string)

	queryParams["login_challenge"] = loginChallenge

	acceptResponse := response.AcceptResponse{}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/login/accept").Return(httpError)

	loginAcceptResponse, genericErr := suite.oauthHandler.AcceptLogin(suite.ginContext, loginChallenge, profile)

	suite.NotNil(genericErr)
	suite.Empty(loginAcceptResponse)
	suite.Equal(&constants.InternalServerError, genericErr)
}

func (suite *OauthLoginHandlerTest) TestAcceptConsentRequest_WhenSuccess() {
	consentChallenge := "consent-challenge-code"

	consentResponse := response.ConsentResponse{}

	queryParams := make(map[string]string)

	queryParams["consent_challenge"] = consentChallenge
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&consentResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		tempResponsePointer := responseBuilder.(*response.ConsentResponse)
		*tempResponsePointer = response.ConsentResponse{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
			RequestedAccessTokenAudience: []string{"some-token"},
			RequestedScope:               []string{"some-scope"},
		}
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Get("https://oauth2/oauth2/auth/requests/consent").Return(nil).Times(1)

	acceptResponse := response.AcceptResponse{}

	acceptRequest := oauth.ConsentAcceptRequest{
		GrantAccessTokenAudience: []string{"some-token"},
		GrantScope:               []string{"some-scope"},
		Remember:                 true,
		RememberFor:              1,
		Session: oauth.Session{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
		},
	}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		tempResponsePointer := responseBuilder.(*response.AcceptResponse)
		*tempResponsePointer = response.AcceptResponse{
			RedirectTo: "http://redirect-to-url",
		}
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/requests/consent/accept").Return(nil).Times(1)

	expectedAcceptConsentResponse := response.AcceptResponse{
		RedirectTo: "http://redirect-to-url",
	}
	actualResponse := response.AcceptResponse{}
	acceptConsentResponse, genericErr := suite.oauthHandler.AcceptConsentRequest(suite.ginContext, consentChallenge)
	suite.Nil(genericErr)
	marshal, err := json.Marshal(acceptConsentResponse)
	suite.Nil(err)
	err = json.Unmarshal(marshal, &actualResponse)
	suite.Nil(err)
	suite.Equal(expectedAcceptConsentResponse, actualResponse)
}

func (suite *OauthLoginHandlerTest) TestAcceptConsentRequest_WhenGetConsentRequestFails() {
	consentChallenge := "consent-challenge-code"

	consentResponse := response.ConsentResponse{}

	queryParams := make(map[string]string)

	authError := response.GenericAuthError{
		Error:      "Not Found",
		StatusCode: 404,
	}

	bytes, _ := json.Marshal(authError)

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusNotFound,
		ResponseBody: bytes,
	}

	queryParams["consent_challenge"] = consentChallenge
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&consentResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		tempResponsePointer := responseBuilder.(*response.ConsentResponse)
		*tempResponsePointer = response.ConsentResponse{}
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Get("https://oauth2/oauth2/auth/requests/consent").Return(httpError).Times(1)

	acceptResponse := response.AcceptResponse{}

	acceptRequest := oauth.ConsentAcceptRequest{
		GrantAccessTokenAudience: []string{"some-token"},
		GrantScope:               []string{"some-scope"},
		Remember:                 true,
		RememberFor:              1,
		Session: oauth.Session{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
		},
	}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(0)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(0)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(0)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest).Times(0)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	}).Times(0)
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/requests/consent/accept").Return(nil).Times(0)

	acceptConsentResponse, genericErr := suite.oauthHandler.AcceptConsentRequest(suite.ginContext, consentChallenge)

	suite.Equal(nil, acceptConsentResponse)
	suite.Equal(&constants.InternalServerError, genericErr)
}

func (suite *OauthLoginHandlerTest) TestAcceptConsentRequest_WhenAcceptConsentRequestFails() {
	consentChallenge := "consent-challenge-code"

	consentResponse := response.ConsentResponse{}

	queryParams := make(map[string]string)

	authError := response.GenericAuthError{
		Error:      "Not Found",
		StatusCode: 404,
	}

	bytes, _ := json.Marshal(authError)

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusNotFound,
		ResponseBody: bytes,
	}

	queryParams["consent_challenge"] = consentChallenge
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&consentResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		tempResponsePointer := responseBuilder.(*response.ConsentResponse)
		*tempResponsePointer = response.ConsentResponse{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
			RequestedAccessTokenAudience: []string{"some-token"},
			RequestedScope:               []string{"some-scope"},
		}
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Get("https://oauth2/oauth2/auth/requests/consent").Return(nil).Times(1)

	acceptResponse := response.AcceptResponse{}

	acceptRequest := oauth.ConsentAcceptRequest{
		GrantAccessTokenAudience: []string{"some-token"},
		GrantScope:               []string{"some-scope"},
		Remember:                 true,
		RememberFor:              1,
		Session: oauth.Session{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
		},
	}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/requests/consent/accept").Return(httpError).Times(1)

	acceptConsentResponse, genericErr := suite.oauthHandler.AcceptConsentRequest(suite.ginContext, consentChallenge)
	suite.NotNil(genericErr)
	suite.Equal(&constants.InvalidConsentChallengeError, genericErr)
	suite.Empty(acceptConsentResponse)
}

func (suite *OauthLoginHandlerTest) TestAcceptConsentRequest_WhenAcceptConsentRequestFailsAndUnmarshalGenericError() {
	consentChallenge := "consent-challenge-code"

	consentResponse := response.ConsentResponse{}

	queryParams := make(map[string]string)

	bytes, _ := json.Marshal("some error")

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusInternalServerError,
		ResponseBody: bytes,
	}

	queryParams["consent_challenge"] = consentChallenge
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&consentResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		tempResponsePointer := responseBuilder.(*response.ConsentResponse)
		*tempResponsePointer = response.ConsentResponse{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
			RequestedAccessTokenAudience: []string{"some-token"},
			RequestedScope:               []string{"some-scope"},
		}
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Get("https://oauth2/oauth2/auth/requests/consent").Return(nil).Times(1)

	acceptResponse := response.AcceptResponse{}

	acceptRequest := oauth.ConsentAcceptRequest{
		GrantAccessTokenAudience: []string{"some-token"},
		GrantScope:               []string{"some-scope"},
		Remember:                 true,
		RememberFor:              1,
		Session: oauth.Session{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
		},
	}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/requests/consent/accept").Return(httpError).Times(1)

	acceptConsentResponse, genericErr := suite.oauthHandler.AcceptConsentRequest(suite.ginContext, consentChallenge)
	suite.NotNil(genericErr)
	suite.Equal(&constants.InternalServerError, genericErr)
	suite.Empty(acceptConsentResponse)
}

func (suite *OauthLoginHandlerTest) TestAcceptConsentRequest_WhenAcceptConsentRequestFailsAndDifferentStatusCodeForGenericError() {
	consentChallenge := "consent-challenge-code"

	consentResponse := response.ConsentResponse{}

	queryParams := make(map[string]string)

	authError := response.GenericAuthError{
		Error:      "Un authorized",
		StatusCode: http.StatusUnauthorized,
	}

	bytes, _ := json.Marshal(authError)

	httpError := golaerror.HttpError{
		StatusCode:   http.StatusUnauthorized,
		ResponseBody: bytes,
	}

	queryParams["consent_challenge"] = consentChallenge
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&consentResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		tempResponsePointer := responseBuilder.(*response.ConsentResponse)
		*tempResponsePointer = response.ConsentResponse{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
			RequestedAccessTokenAudience: []string{"some-token"},
			RequestedScope:               []string{"some-scope"},
		}
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Get("https://oauth2/oauth2/auth/requests/consent").Return(nil).Times(1)

	acceptResponse := response.AcceptResponse{}

	acceptRequest := oauth.ConsentAcceptRequest{
		GrantAccessTokenAudience: []string{"some-token"},
		GrantScope:               []string{"some-scope"},
		Remember:                 true,
		RememberFor:              1,
		Session: oauth.Session{
			Userprofile: db.UserProfile{
				UserID:   "some-user",
				Username: "some-user-name",
				Email:    "dummy@gmail.com",
				IsActive: true,
			},
		},
	}

	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().AddQueryParameters(queryParams).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithJSONBody(acceptRequest).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&acceptResponse).DoAndReturn(func(responseBuilder interface{}) request.HttpRequest {
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Put("https://oauth2/oauth2/auth/requests/consent/accept").Return(httpError).Times(1)

	acceptConsentResponse, genericErr := suite.oauthHandler.AcceptConsentRequest(suite.ginContext, consentChallenge)
	suite.NotNil(genericErr)
	suite.Equal(&constants.InternalServerError, genericErr)
	suite.Empty(acceptConsentResponse)
}

func (suite *OauthLoginHandlerTest) TestShouldReturnAccessTokenAndEncryptedIdTokenWhenHttpClientRequestReturnsIdAndAccessToken() {
	encryptedIdToken, idToken, userId := "encrypted-id-token", "random-id-token", "userId"
	istLocation, _ := time.LoadLocation("Asia/Kolkata")
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2020-08-27 14:04:01", istLocation)
	exchangeRequest := oauth.TokenExchangeRequest{
		CodeVerifier: "random-code-verifier",
		ClientId:     "random-client-id",
		RedirectUri:  "random-redirect-uri",
		Code:         "random-code",
	}
	expectedOauthResponse := response.TokenExchangeResponse{
		AccessToken: "random-access-token",
		ExpiresIn:   3600,
		IdToken:     idToken,
	}
	expectedServiceResponse := response.TokenExchangeResponse{
		AccessToken:      "random-access-token",
		ExpiresIn:        3600,
		IdToken:          idToken,
		EncryptedIdToken: encryptedIdToken,
		ExpiresAt:        "Thu, 27 Aug 2020 09:34:01 GMT",
	}

	token := model.IdToken{UserId: userId}

	formData := make(map[string]interface{})
	formData["grant_type"] = "authorization_code"
	formData["code_verifier"] = "random-code-verifier"
	formData["client_id"] = "random-client-id"
	formData["redirect_uri"] = "random-redirect-uri"
	formData["code"] = "random-code"

	tokenExchangeResponse := response.TokenExchangeResponse{}
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(formData).Return(suite.mockHttpRequest).Times(1)
	suite.mockHttpRequest.EXPECT().ResponseAs(&tokenExchangeResponse).DoAndReturn(func(responseHolder interface{}) request.HttpRequest {
		tempResponsePointer := responseHolder.(*response.TokenExchangeResponse)
		*tempResponsePointer = expectedOauthResponse
		return suite.mockHttpRequest
	}).Times(1)
	suite.mockHttpRequest.EXPECT().Post("https://oauth2/token-url").Return(nil).Times(1)
	suite.mockClock.EXPECT().Now().Return(currentTime).Times(1)
	suite.mockOauthUtils.EXPECT().DecodeIdTokenFromJWT(idToken).Return(token, nil).Times(1)
	suite.mockOauthUtils.EXPECT().EncryptIdToken(suite.ginContext, idToken).Return(encryptedIdToken, nil).Times(1)

	tokenExchangeResponse, idTokenResponse, err := suite.oauthHandler.ExchangeToken(suite.ginContext, exchangeRequest)

	suite.Equal(expectedServiceResponse, tokenExchangeResponse)
	suite.Equal(token, idTokenResponse)
	suite.Nil(err)
}

func (suite *OauthLoginHandlerTest) TestShouldReturnExpiresAtFormattedInHttpTimeFormat() {
	encryptedIdToken, idToken, userId := "encrypted-id-token", "random-id-token", "userId"
	currentTime, _ := time.Parse("2006-01-02 15:04:05 MST", "2020-08-27 14:04:01 UTC")
	exchangeRequest := oauth.TokenExchangeRequest{
		CodeVerifier: "random-code-verifier",
		ClientId:     "random-client-id",
		RedirectUri:  "random-redirect-uri",
		Code:         "random-code",
	}
	expectedOauthResponse := response.TokenExchangeResponse{
		AccessToken: "random-access-token",
		ExpiresIn:   3600,
		IdToken:     idToken,
	}
	expectedServiceResponse := response.TokenExchangeResponse{
		AccessToken:      "random-access-token",
		ExpiresIn:        3600,
		IdToken:          idToken,
		EncryptedIdToken: encryptedIdToken,
		ExpiresAt:        "Thu, 27 Aug 2020 15:04:01 GMT",
	}

	token := model.IdToken{UserId: userId}

	formData := make(map[string]interface{})
	formData["grant_type"] = "authorization_code"
	formData["code_verifier"] = "random-code-verifier"
	formData["client_id"] = "random-client-id"
	formData["redirect_uri"] = "random-redirect-uri"
	formData["code"] = "random-code"

	tokenExchangeResponse := response.TokenExchangeResponse{}
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(formData).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&tokenExchangeResponse).DoAndReturn(func(responseHolder interface{}) request.HttpRequest {
		tempResponsePointer := responseHolder.(*response.TokenExchangeResponse)
		*tempResponsePointer = expectedOauthResponse
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Post("https://oauth2/token-url").Return(nil)
	suite.mockClock.EXPECT().Now().Return(currentTime)
	suite.mockOauthUtils.EXPECT().DecodeIdTokenFromJWT(idToken).Return(token, nil)
	suite.mockOauthUtils.EXPECT().EncryptIdToken(suite.ginContext, idToken).Return(encryptedIdToken, nil)

	tokenExchangeResponse, idTokenResponse, err := suite.oauthHandler.ExchangeToken(suite.ginContext, exchangeRequest)

	suite.Equal(expectedServiceResponse, tokenExchangeResponse)
	suite.Equal(token, idTokenResponse)
	suite.Nil(err)
}

func (suite *OauthLoginHandlerTest) TestShouldReturnInternalServerErrorWhenHttpClientRequestReturnsAnError() {
	exchangeRequest := oauth.TokenExchangeRequest{
		CodeVerifier: "random-code-verifier",
		ClientId:     "random-client-id",
		RedirectUri:  "random-redirect-uri",
		Code:         "random-code",
	}

	formData := make(map[string]interface{})
	formData["grant_type"] = "authorization_code"
	formData["code_verifier"] = "random-code-verifier"
	formData["client_id"] = "random-client-id"
	formData["redirect_uri"] = "random-redirect-uri"
	formData["code"] = "random-code"

	bytes, _ := json.Marshal("Http Error!!")
	httpError := golaerror.HttpError{ResponseBody: bytes, StatusCode: http.StatusInternalServerError}
	tokenExchangeResponse := response.TokenExchangeResponse{}
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(formData).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&tokenExchangeResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("https://oauth2/token-url").Return(httpError)

	tokenExchangeResponse, token, err := suite.oauthHandler.ExchangeToken(suite.ginContext, exchangeRequest)
	suite.Equal(&constants.InternalServerError, err)
	suite.Equal(response.TokenExchangeResponse{}, tokenExchangeResponse)
	suite.Equal(model.IdToken{}, token)
}

func (suite *OauthLoginHandlerTest) TestShouldReturnBadRequestErrorWhenHttpResponseDoesNotContainAccessToken() {
	exchangeRequest := oauth.TokenExchangeRequest{
		CodeVerifier: "random-code-verifier",
		ClientId:     "random-client-id",
		RedirectUri:  "random-redirect-uri",
		Code:         "random-code",
	}

	formData := make(map[string]interface{})
	formData["grant_type"] = "authorization_code"
	formData["code_verifier"] = "random-code-verifier"
	formData["client_id"] = "random-client-id"
	formData["redirect_uri"] = "random-redirect-uri"
	formData["code"] = "random-code"

	genericAuthError := response.GenericAuthError{Error: "Error", StatusCode: 401}
	bytes, _ := json.Marshal(genericAuthError)
	httpError := golaerror.HttpError{ResponseBody: bytes, StatusCode: http.StatusUnauthorized}
	tokenExchangeResponse := response.TokenExchangeResponse{}
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(formData).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&tokenExchangeResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("https://oauth2/token-url").Return(httpError)

	tokenExchangeResponse, token, actualError := suite.oauthHandler.ExchangeToken(suite.ginContext, exchangeRequest)

	suite.Equal(response.TokenExchangeResponse{}, tokenExchangeResponse)
	suite.Equal(model.IdToken{}, token)
	suite.Equal(&constants.PayloadValidationError, actualError)
}

func (suite *OauthLoginHandlerTest) TestShouldReturnServerErrorWhenDecodingOfIdTokenReturnsError() {
	idToken := "random-id-token"
	errorMsg := "encryption error"
	exchangeRequest := oauth.TokenExchangeRequest{
		CodeVerifier: "random-code-verifier",
		ClientId:     "random-client-id",
		RedirectUri:  "random-redirect-uri",
		Code:         "random-code",
	}
	expectedOauthResponse := response.TokenExchangeResponse{
		AccessToken: "random-access-token",
		ExpiresIn:   0,
		IdToken:     idToken,
	}

	formData := make(map[string]interface{})
	formData["grant_type"] = "authorization_code"
	formData["code_verifier"] = "random-code-verifier"
	formData["client_id"] = "random-client-id"
	formData["redirect_uri"] = "random-redirect-uri"
	formData["code"] = "random-code"

	tokenExchangeResponse := response.TokenExchangeResponse{}
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(formData).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&tokenExchangeResponse).DoAndReturn(func(responseHolder interface{}) request.HttpRequest {
		tempResponsePointer := responseHolder.(*response.TokenExchangeResponse)
		*tempResponsePointer = expectedOauthResponse
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Post("https://oauth2/token-url").Return(nil)
	suite.mockClock.EXPECT().Now().Return(time.Now())
	suite.mockOauthUtils.EXPECT().DecodeIdTokenFromJWT(idToken).Return(model.IdToken{}, errors.New(errorMsg))

	tokenExchangeResponse, token, err := suite.oauthHandler.ExchangeToken(suite.ginContext, exchangeRequest)

	suite.Equal(&constants.InternalServerError, err)
	suite.Equal(model.IdToken{}, token)
	suite.Empty(tokenExchangeResponse)
}

func (suite *OauthLoginHandlerTest) TestShouldReturnServerErrorWhenEncryptionOfIdTokenReturnsError() {
	idToken, userId := "random-id-token", "userId"
	errorMsg := "encryption error"
	exchangeRequest := oauth.TokenExchangeRequest{
		CodeVerifier: "random-code-verifier",
		ClientId:     "random-client-id",
		RedirectUri:  "random-redirect-uri",
		Code:         "random-code",
	}
	expectedOauthResponse := response.TokenExchangeResponse{
		AccessToken: "random-access-token",
		ExpiresIn:   0,
		IdToken:     idToken,
	}

	formData := make(map[string]interface{})
	formData["grant_type"] = "authorization_code"
	formData["code_verifier"] = "random-code-verifier"
	formData["client_id"] = "random-client-id"
	formData["redirect_uri"] = "random-redirect-uri"
	formData["code"] = "random-code"

	tokenExchangeResponse := response.TokenExchangeResponse{}
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(formData).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&tokenExchangeResponse).DoAndReturn(func(responseHolder interface{}) request.HttpRequest {
		tempResponsePointer := responseHolder.(*response.TokenExchangeResponse)
		*tempResponsePointer = expectedOauthResponse
		return suite.mockHttpRequest
	})
	suite.mockHttpRequest.EXPECT().Post("https://oauth2/token-url").Return(nil)
	suite.mockClock.EXPECT().Now().Return(time.Now())
	suite.mockOauthUtils.EXPECT().DecodeIdTokenFromJWT(idToken).Return(model.IdToken{UserId: userId}, nil)
	suite.mockOauthUtils.EXPECT().EncryptIdToken(suite.ginContext, idToken).Return("", errors.New(errorMsg))

	tokenExchangeResponse, token, err := suite.oauthHandler.ExchangeToken(suite.ginContext, exchangeRequest)

	suite.Equal(&constants.InternalServerError, err)
	suite.Equal(model.IdToken{}, token)
	suite.Empty(tokenExchangeResponse)
}

func (suite *OauthLoginHandlerTest) TestShouldReturnInternalServerErrorWhenGetTokenReturnsGenericErrorOfInternalServerError() {
	exchangeRequest := oauth.TokenExchangeRequest{
		CodeVerifier: "random-code-verifier",
		ClientId:     "random-client-id",
		RedirectUri:  "random-redirect-uri",
		Code:         "random-code",
	}

	formData := make(map[string]interface{})
	formData["grant_type"] = "authorization_code"
	formData["code_verifier"] = "random-code-verifier"
	formData["client_id"] = "random-client-id"
	formData["redirect_uri"] = "random-redirect-uri"
	formData["code"] = "random-code"

	genericAuthError := response.GenericAuthError{Error: "Error", StatusCode: 500}
	bytes, _ := json.Marshal(genericAuthError)
	httpError := golaerror.HttpError{ResponseBody: bytes, StatusCode: http.StatusInternalServerError}
	tokenExchangeResponse := response.TokenExchangeResponse{}
	suite.mockHttpRequestBuilder.EXPECT().NewRequest().Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithContext(suite.ginContext).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().WithFormURLEncoded(formData).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().ResponseAs(&tokenExchangeResponse).Return(suite.mockHttpRequest)
	suite.mockHttpRequest.EXPECT().Post("https://oauth2/token-url").Return(httpError)

	tokenExchangeResponse, token, actualError := suite.oauthHandler.ExchangeToken(suite.ginContext, exchangeRequest)

	suite.Equal(response.TokenExchangeResponse{}, tokenExchangeResponse)
	suite.Equal(model.IdToken{}, token)
	suite.Equal(&constants.InternalServerError, actualError)
}
