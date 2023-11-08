package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"post-api/idp/constants"
	"post-api/idp/mocks"
	"post-api/idp/models/request"
	"post-api/idp/models/response"
	"testing"
)

type LoginControllerTest struct {
	suite.Suite
	mockCtrl         *gomock.Controller
	recorder         *httptest.ResponseRecorder
	context          *gin.Context
	mockLoginService *mocks.MockLoginService
	mockAuthHandler  *mocks.MockOauthLoginHandler
	loginController  LoginController
}

func (suite *LoginControllerTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockLoginService = mocks.NewMockLoginService(suite.mockCtrl)
	suite.mockAuthHandler = mocks.NewMockOauthLoginHandler(suite.mockCtrl)
	suite.loginController = NewLoginController(suite.mockLoginService, suite.mockAuthHandler)
}

func (suite *LoginControllerTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestLoginControllerTestSuite(t *testing.T) {
	suite.Run(t, new(LoginControllerTest))
}

func (suite *LoginControllerTest) TestLoginByEmailAndPassword_WhenSuccess() {
	requestBody := `{"email":"dummy@gmail.com","password":"!QASW!@scv##@$TV7rgb","login_challenge":"some-challenge-code"}`

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "api/idp/v1/login/password", bytes.NewBufferString(requestBody))

	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "!QASW!@scv##@$TV7rgb",
		LoginChallenge: "some-challenge-code",
	}

	acceptResponse := response.AcceptResponse{RedirectTo: "http://some-redirect-url"}
	suite.mockLoginService.EXPECT().LoginWithEmailAndPassword(loginRequest, suite.context).Return(acceptResponse, nil).Times(1)

	suite.loginController.LoginByEmailAndPassword(suite.context)
	actualAcceptResponse := response.AcceptResponse{}
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &actualAcceptResponse)
	suite.Nil(err)
	suite.Equal(acceptResponse, actualAcceptResponse)
}

func (suite *LoginControllerTest) TestLoginByEmailAndPassword_WhenBadRequest() {
	requestBody := `{"email":"dummy@gmail.com","password":"","login_challenge":"some-challenge-code"}`

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "api/idp/v1/login/password", bytes.NewBufferString(requestBody))

	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "!QASW!@scv##@$TV7rgb",
		LoginChallenge: "some-challenge-code",
	}

	acceptResponse := response.AcceptResponse{RedirectTo: "http://some-redirect-url"}
	suite.mockLoginService.EXPECT().LoginWithEmailAndPassword(loginRequest, suite.context).Return(acceptResponse, nil).Times(0)

	suite.loginController.LoginByEmailAndPassword(suite.context)
	actualError := golaerror.Error{}
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &actualError)
	suite.Nil(err)
	suite.Equal(constants.PayloadValidationError, actualError)
}

func (suite *LoginControllerTest) TestLoginByEmailAndPassword_WhenLoginServiceFails() {
	requestBody := `{"email":"dummy@gmail.com","password":"!QASW!@scv##@$TV7rgb","login_challenge":"some-challenge-code"}`

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "api/idp/v1/login/password", bytes.NewBufferString(requestBody))

	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "!QASW!@scv##@$TV7rgb",
		LoginChallenge: "some-challenge-code",
	}

	acceptResponse := response.AcceptResponse{}
	suite.mockLoginService.EXPECT().LoginWithEmailAndPassword(loginRequest, suite.context).Return(acceptResponse, &constants.InvalidLoginChallengeError).Times(1)

	suite.loginController.LoginByEmailAndPassword(suite.context)
	actualErr := golaerror.Error{}
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &actualErr)
	suite.Nil(err)
	suite.Equal(constants.InvalidLoginChallengeError, actualErr)
}

func (suite *LoginControllerTest) TestShouldReturnStatusOkWithRedirectUrlForConsentIsAccepted() {
	suite.context.Request, _ = http.NewRequest("GET", "/?consent_challenge=challenge", nil)

	response := response.AcceptResponse{RedirectTo: "http://redirectTo"}
	suite.mockAuthHandler.EXPECT().AcceptConsentRequest(suite.context, "challenge").Return(response, nil)

	suite.loginController.GrantConsent(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	bodyBytes, _ := ioutil.ReadAll(suite.recorder.Body)
	body := string(bodyBytes)

	expectedBody, _ := json.Marshal(response)
	suite.Equal(string(expectedBody), body)
}

func (suite *LoginControllerTest) TestShouldReturnErrorWhenConsentRequestFails() {
	suite.context.Request, _ = http.NewRequest("GET", "/?consent_challenge=challenge", nil)
	suite.mockAuthHandler.EXPECT().AcceptConsentRequest(suite.context, "challenge").Return(nil, &constants.InternalServerError)

	suite.loginController.GrantConsent(suite.context)
	actualErr := golaerror.Error{}
	err := json.Unmarshal(suite.recorder.Body.Bytes(), &actualErr)
	suite.Nil(err)
	suite.Equal(constants.InternalServerError, actualErr)
}
