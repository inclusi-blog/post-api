package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/idp/constants"
	mocks2 "post-api/idp/mocks"
	"testing"
)

type AuthenticatorServiceTest struct {
	suite.Suite
	mockController            *gomock.Controller
	ginContext                *gin.Context
	authenticatorService      AuthenticatorService
	mockHashUtils             *mocks2.MockHashUtil
	mockUserDetailsRepository *mocks2.MockUserDetailsRepository
}

func TestAuthenticatorServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthenticatorServiceTest))
}

func (suite *AuthenticatorServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.ginContext, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.ginContext.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	suite.mockUserDetailsRepository = mocks2.NewMockUserDetailsRepository(suite.mockController)
	suite.mockHashUtils = mocks2.NewMockHashUtil(suite.mockController)
	suite.authenticatorService = NewAuthenticatorService(suite.mockHashUtils, suite.mockUserDetailsRepository)
}

func (suite *AuthenticatorServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *AuthenticatorServiceTest) TestAuthenticate_WhenDetailsRepositoryMatchesPassword() {
	planTextPassword := "some-password"
	email := "dummy@gmail.com"

	suite.mockUserDetailsRepository.EXPECT().GetPassword(email, suite.ginContext).Return("some-hashed-password", nil).Times(1)
	suite.mockHashUtils.EXPECT().MatchBcryptHash("some-hashed-password", planTextPassword).Return(nil).Times(1)

	authErr := suite.authenticatorService.Authenticate(suite.ginContext, planTextPassword, email)
	suite.Nil(authErr)
}

func (suite *AuthenticatorServiceTest) TestAuthenticate_WhenDetailsRepositoryReturnsError() {
	planTextPassword := "some-password"
	email := "dummy@gmail.com"

	suite.mockUserDetailsRepository.EXPECT().GetPassword(email, suite.ginContext).Return("", errors.New("something went wrong")).Times(1)
	suite.mockHashUtils.EXPECT().MatchBcryptHash("some-hashed-password", planTextPassword).Return(nil).Times(0)

	authErr := suite.authenticatorService.Authenticate(suite.ginContext, planTextPassword, email)
	suite.NotNil(authErr)
	suite.Equal(&constants.InternalServerError, authErr)
}

func (suite *AuthenticatorServiceTest) TestAuthenticate_WhenMatchHashReturnsError() {
	planTextPassword := "some-password"
	email := "dummy@gmail.com"

	suite.mockUserDetailsRepository.EXPECT().GetPassword(email, suite.ginContext).Return("some-hashed-password", nil).Times(1)
	suite.mockHashUtils.EXPECT().MatchBcryptHash("some-hashed-password", planTextPassword).Return(errors.New("invalid password")).Times(1)

	authErr := suite.authenticatorService.Authenticate(suite.ginContext, planTextPassword, email)
	suite.NotNil(authErr)
	suite.Equal(&constants.InvalidCredentialsError, authErr)
}
