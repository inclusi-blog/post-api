package service

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/mocks"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/idp/constants"
	mocks2 "post-api/idp/mocks"
	"post-api/idp/models/db"
	"post-api/idp/models/request"
	"post-api/idp/models/response"
	"testing"
)

type LoginServiceTest struct {
	suite.Suite
	mockController           *gomock.Controller
	ginContext               *gin.Context
	detailsRepository        *mocks2.MockUserDetailsRepository
	mockAuthenticatorService *mocks2.MockAuthenticatorService
	mockCryptoUtils          *mocks.MockCryptoUtil
	mockOauthLoginHandler    *mocks2.MockOauthLoginHandler
	loginService             LoginService
}

func TestLoginServiceTestSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceTest))
}

func (suite *LoginServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.ginContext, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.ginContext.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	suite.detailsRepository = mocks2.NewMockUserDetailsRepository(suite.mockController)
	suite.mockAuthenticatorService = mocks2.NewMockAuthenticatorService(suite.mockController)
	suite.mockCryptoUtils = mocks.NewMockCryptoUtil(suite.mockController)
	suite.mockOauthLoginHandler = mocks2.NewMockOauthLoginHandler(suite.mockController)
	suite.loginService = NewLoginService(suite.detailsRepository, suite.mockCryptoUtils, suite.mockAuthenticatorService, suite.mockOauthLoginHandler)
}

func (suite *LoginServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *LoginServiceTest) TestLoginWithEmailAndPassword_WhenSuccess() {
	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "some-hashed-password",
		LoginChallenge: "login-challenge",
	}

	profile := db.UserProfile{
		Id:       uuid.New(),
		Username: "some-user-name",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptResponse := response.AcceptResponse{RedirectTo: "https://some-redirect-url"}

	suite.detailsRepository.EXPECT().IsEmailAvailable(loginRequest.Email, suite.ginContext).Return(true, nil).Times(1)
	suite.mockCryptoUtils.EXPECT().Decipher(suite.ginContext, loginRequest.Password).Return("plan-text-password", nil).Times(1)
	suite.detailsRepository.EXPECT().GetUserProfile(loginRequest.Email, suite.ginContext).Return(profile, nil).Times(1)
	suite.mockAuthenticatorService.EXPECT().Authenticate(suite.ginContext, "plan-text-password", "dummy@gmail.com").Return(nil).Times(1)
	suite.mockOauthLoginHandler.EXPECT().AcceptLogin(suite.ginContext, "login-challenge", profile).Return(acceptResponse, nil).Times(1)

	actualAcceptResponse, err := suite.loginService.LoginWithEmailAndPassword(loginRequest, suite.ginContext)
	suite.Nil(err)
	suite.Equal(acceptResponse, actualAcceptResponse)
}

func (suite *LoginServiceTest) TestLoginWithEmailAndPassword_UserEmailNotAvailable() {
	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "some-hashed-password",
		LoginChallenge: "login-challenge",
	}

	profile := db.UserProfile{
		Id:       uuid.New(),
		Username: "some-user-name",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptResponse := response.AcceptResponse{RedirectTo: "https://some-redirect-url"}

	suite.detailsRepository.EXPECT().IsEmailAvailable(loginRequest.Email, suite.ginContext).Return(false, nil).Times(1)
	suite.mockCryptoUtils.EXPECT().Decipher(suite.ginContext, loginRequest.Password).Return("plan-text-password", nil).Times(0)
	suite.detailsRepository.EXPECT().GetUserProfile(loginRequest.Email, suite.ginContext).Return(profile, nil).Times(0)
	suite.mockAuthenticatorService.EXPECT().Authenticate(suite.ginContext, "plan-text-password", "dummy@gmail.com").Return(nil).Times(0)
	suite.mockOauthLoginHandler.EXPECT().AcceptLogin(suite.ginContext, "login-challenge", profile).Return(acceptResponse, nil).Times(0)

	actualAcceptResponse, err := suite.loginService.LoginWithEmailAndPassword(loginRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.UserNotFoundError, err)
	suite.Equal(nil, actualAcceptResponse)
}

func (suite *LoginServiceTest) TestLoginWithEmailAndPassword_WhenUserRepositoryReturnEmailAvailabilityError() {
	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "some-hashed-password",
		LoginChallenge: "login-challenge",
	}

	profile := db.UserProfile{
		Id:       uuid.New(),
		Username: "some-user-name",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptResponse := response.AcceptResponse{RedirectTo: "https://some-redirect-url"}

	suite.detailsRepository.EXPECT().IsEmailAvailable(loginRequest.Email, suite.ginContext).Return(false, errors.New("something went wrong")).Times(1)
	suite.mockCryptoUtils.EXPECT().Decipher(suite.ginContext, loginRequest.Password).Return("plan-text-password", nil).Times(0)
	suite.detailsRepository.EXPECT().GetUserProfile(loginRequest.Email, suite.ginContext).Return(profile, nil).Times(0)
	suite.mockAuthenticatorService.EXPECT().Authenticate(suite.ginContext, "plan-text-password", "dummy@gmail.com").Return(nil).Times(0)
	suite.mockOauthLoginHandler.EXPECT().AcceptLogin(suite.ginContext, "login-challenge", profile).Return(acceptResponse, nil).Times(0)

	actualAcceptResponse, err := suite.loginService.LoginWithEmailAndPassword(loginRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InternalServerError, err)
	suite.Equal(nil, actualAcceptResponse)
}

func (suite *LoginServiceTest) TestLoginWithEmailAndPassword_WhenUserDeciphersReturnsError() {
	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "some-hashed-password",
		LoginChallenge: "login-challenge",
	}

	profile := db.UserProfile{
		Id:       uuid.New(),
		Username: "some-user-name",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptResponse := response.AcceptResponse{RedirectTo: "https://some-redirect-url"}

	suite.detailsRepository.EXPECT().IsEmailAvailable(loginRequest.Email, suite.ginContext).Return(true, nil).Times(1)
	suite.mockCryptoUtils.EXPECT().Decipher(suite.ginContext, loginRequest.Password).Return("", errors.New("something went wrong")).Times(1)
	suite.detailsRepository.EXPECT().GetUserProfile(loginRequest.Email, suite.ginContext).Return(profile, nil).Times(0)
	suite.mockAuthenticatorService.EXPECT().Authenticate(suite.ginContext, "plan-text-password", "dummy@gmail.com").Return(nil).Times(0)
	suite.mockOauthLoginHandler.EXPECT().AcceptLogin(suite.ginContext, "login-challenge", profile).Return(acceptResponse, nil).Times(0)

	actualAcceptResponse, err := suite.loginService.LoginWithEmailAndPassword(loginRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InternalServerError, err)
	suite.Equal(nil, actualAcceptResponse)
}

func (suite *LoginServiceTest) TestLoginWithEmailAndPassword_WhenGetUserProfileReturnsError() {
	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "some-hashed-password",
		LoginChallenge: "login-challenge",
	}

	profile := db.UserProfile{}

	acceptResponse := response.AcceptResponse{RedirectTo: "https://some-redirect-url"}

	suite.detailsRepository.EXPECT().IsEmailAvailable(loginRequest.Email, suite.ginContext).Return(true, nil).Times(1)
	suite.mockCryptoUtils.EXPECT().Decipher(suite.ginContext, loginRequest.Password).Return("plain-text-password", nil).Times(1)
	suite.detailsRepository.EXPECT().GetUserProfile(loginRequest.Email, suite.ginContext).Return(profile, sql.ErrNoRows).Times(1)
	suite.mockAuthenticatorService.EXPECT().Authenticate(suite.ginContext, "plan-text-password", "dummy@gmail.com").Return(nil).Times(0)
	suite.mockOauthLoginHandler.EXPECT().AcceptLogin(suite.ginContext, "login-challenge", profile).Return(acceptResponse, nil).Times(0)

	actualAcceptResponse, err := suite.loginService.LoginWithEmailAndPassword(loginRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InternalServerError, err)
	suite.Equal(nil, actualAcceptResponse)
}

func (suite *LoginServiceTest) TestLoginWithEmailAndPassword_WhenAuthenticatorReturnsError() {
	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "some-hashed-password",
		LoginChallenge: "login-challenge",
	}

	profile := db.UserProfile{
		Id:       uuid.New(),
		Username: "some-user-name",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptResponse := response.AcceptResponse{RedirectTo: "https://some-redirect-url"}

	suite.detailsRepository.EXPECT().IsEmailAvailable(loginRequest.Email, suite.ginContext).Return(true, nil).Times(1)
	suite.mockCryptoUtils.EXPECT().Decipher(suite.ginContext, loginRequest.Password).Return("plain-text-password", nil).Times(1)
	suite.detailsRepository.EXPECT().GetUserProfile(loginRequest.Email, suite.ginContext).Return(profile, nil).Times(1)
	suite.mockAuthenticatorService.EXPECT().Authenticate(suite.ginContext, "plain-text-password", "dummy@gmail.com").Return(&constants.InternalServerError).Times(1)
	suite.mockOauthLoginHandler.EXPECT().AcceptLogin(suite.ginContext, "login-challenge", profile).Return(acceptResponse, nil).Times(0)

	actualAcceptResponse, err := suite.loginService.LoginWithEmailAndPassword(loginRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InternalServerError, err)
	suite.Equal(nil, actualAcceptResponse)
}

func (suite *LoginServiceTest) TestLoginWithEmailAndPassword_WhenAuthenticatorReturnsInvalidCredentialsError() {
	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "some-hashed-password",
		LoginChallenge: "login-challenge",
	}

	profile := db.UserProfile{
		Id:       uuid.New(),
		Username: "some-user-name",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptResponse := response.AcceptResponse{RedirectTo: "https://some-redirect-url"}

	suite.detailsRepository.EXPECT().IsEmailAvailable(loginRequest.Email, suite.ginContext).Return(true, nil).Times(1)
	suite.mockCryptoUtils.EXPECT().Decipher(suite.ginContext, loginRequest.Password).Return("plain-text-password", nil).Times(1)
	suite.detailsRepository.EXPECT().GetUserProfile(loginRequest.Email, suite.ginContext).Return(profile, nil).Times(1)
	suite.mockAuthenticatorService.EXPECT().Authenticate(suite.ginContext, "plain-text-password", "dummy@gmail.com").Return(&constants.InvalidCredentialsError).Times(1)
	suite.mockOauthLoginHandler.EXPECT().AcceptLogin(suite.ginContext, "login-challenge", profile).Return(acceptResponse, nil).Times(0)

	actualAcceptResponse, err := suite.loginService.LoginWithEmailAndPassword(loginRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InvalidCredentialsError, err)
	suite.Equal(nil, actualAcceptResponse)
}

func (suite *LoginServiceTest) TestLoginWithEmailAndPassword_WhenAcceptLoginReturnsInvalidLoginChanllengeError() {
	loginRequest := request.UserLoginRequest{
		Email:          "dummy@gmail.com",
		Password:       "some-hashed-password",
		LoginChallenge: "login-challenge",
	}

	profile := db.UserProfile{
		Id:       uuid.New(),
		Username: "some-user-name",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	acceptResponse := response.AcceptResponse{}

	suite.detailsRepository.EXPECT().IsEmailAvailable(loginRequest.Email, suite.ginContext).Return(true, nil).Times(1)
	suite.mockCryptoUtils.EXPECT().Decipher(suite.ginContext, loginRequest.Password).Return("plain-text-password", nil).Times(1)
	suite.detailsRepository.EXPECT().GetUserProfile(loginRequest.Email, suite.ginContext).Return(profile, nil).Times(1)
	suite.mockAuthenticatorService.EXPECT().Authenticate(suite.ginContext, "plain-text-password", "dummy@gmail.com").Return(nil).Times(1)
	suite.mockOauthLoginHandler.EXPECT().AcceptLogin(suite.ginContext, "login-challenge", profile).Return(acceptResponse, &constants.InvalidLoginChallengeError).Times(1)

	actualAcceptResponse, err := suite.loginService.LoginWithEmailAndPassword(loginRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InvalidLoginChallengeError, err)
	suite.Equal(nil, actualAcceptResponse)
}
