package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	mocksUtil "github.com/gola-glitch/gola-utils/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/idp/constants"
	"post-api/idp/mocks"
	"post-api/idp/models/db"
	"post-api/idp/models/request"
	"post-api/idp/models/response"
	"testing"
)

type UserRegistrationServiceTest struct {
	suite.Suite
	mockController            *gomock.Controller
	ginContext                *gin.Context
	mockUserDetailsRepository *mocks.MockUserDetailsRepository
	mockCryptoUtil            *mocksUtil.MockCryptoUtil
	mockRedisStore            *mocks.MockRedisStore
	mockHashUtil              *mocks.MockHashUtil
	userRegistrationService   UserRegistrationService
}

func TestUserRegistrationServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserRegistrationServiceTest))
}

func (suite *UserRegistrationServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.ginContext, _ = gin.CreateTestContext(httptest.NewRecorder())
	suite.ginContext.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	suite.mockUserDetailsRepository = mocks.NewMockUserDetailsRepository(suite.mockController)
	suite.mockCryptoUtil = mocksUtil.NewMockCryptoUtil(suite.mockController)
	suite.mockRedisStore = mocks.NewMockRedisStore(suite.mockController)
	suite.mockHashUtil = mocks.NewMockHashUtil(suite.mockController)
	suite.userRegistrationService = NewUserRegistrationService(suite.mockUserDetailsRepository, suite.mockCryptoUtil, suite.mockRedisStore, suite.mockHashUtil)
}

func (suite *UserRegistrationServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *UserRegistrationServiceTest) TestInitiateRegistration_WhenSuccess() {
	registrationRequest := request.InitiateRegistrationRequest{
		Email:    "encrypted-username",
		Password: "encrypted-password",
		Username: "user123",
	}
	suite.mockCryptoUtil.EXPECT().Decipher(suite.ginContext, registrationRequest.Password).Return("decrypted-password", nil).Times(1)
	suite.mockUserDetailsRepository.EXPECT().IsUserNameAndEmailAvailable(registrationRequest.Username, registrationRequest.Email, suite.ginContext).Return(true, nil).Times(1)
	suite.mockHashUtil.EXPECT().GenerateBcryptHash("decrypted-password").Return("hashed-password", nil).Times(1)

	err := suite.userRegistrationService.InitiateRegistration(registrationRequest, suite.ginContext)
	suite.Nil(err)
}

func (suite *UserRegistrationServiceTest) TestInitiateRegistration_WhenPasswordDecipherFails() {
	registrationRequest := request.InitiateRegistrationRequest{
		Email:    "encrypted-username",
		Password: "encrypted-password",
		Username: "user123",
	}
	suite.mockCryptoUtil.EXPECT().Decipher(suite.ginContext, registrationRequest.Password).Return("", errors.New("unable to decrypt password")).Times(1)
	suite.mockUserDetailsRepository.EXPECT().IsUserNameAndEmailAvailable(registrationRequest.Username, registrationRequest.Email, suite.ginContext).Return(true, nil).Times(0)

	err := suite.userRegistrationService.InitiateRegistration(registrationRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InternalServerError, err)
}

func (suite *UserRegistrationServiceTest) TestInitiateRegistration_WhenUserExistenceFails() {
	registrationRequest := request.InitiateRegistrationRequest{
		Email:    "encrypted-username",
		Password: "encrypted-password",
		Username: "user123",
	}
	suite.mockCryptoUtil.EXPECT().Decipher(suite.ginContext, registrationRequest.Password).Return("decrypted-password", nil).Times(1)
	suite.mockUserDetailsRepository.EXPECT().IsUserNameAndEmailAvailable(registrationRequest.Username, registrationRequest.Email, suite.ginContext).Return(false, errors.New("something went wrong")).Times(1)

	err := suite.userRegistrationService.InitiateRegistration(registrationRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InternalServerError, err)
}

func (suite *UserRegistrationServiceTest) TestInitiateRegistration_WhenNoUserOrEmailExists() {
	userUUID := fmt.Sprintf("%s", uuid.New())
	registrationRequest := request.InitiateRegistrationRequest{
		Email:    "encrypted-username",
		Password: "encrypted-password",
		Username: "user123",
		UUID:     userUUID,
	}

	userDetails := db.SaveUserDetails{
		UUID:     userUUID,
		Username: registrationRequest.Username,
		Email:    registrationRequest.Email,
		Password: "hashed-password",
		IsActive: true,
	}
	suite.mockCryptoUtil.EXPECT().Decipher(suite.ginContext, registrationRequest.Password).Return("decrypted-password", nil).Times(1)
	suite.mockUserDetailsRepository.EXPECT().IsUserNameAndEmailAvailable(registrationRequest.Username, registrationRequest.Email, suite.ginContext).Return(false, nil).Times(1)
	suite.mockHashUtil.EXPECT().GenerateBcryptHash("decrypted-password").Return("hashed-password", nil).Times(1)
	suite.mockUserDetailsRepository.EXPECT().SaveUserDetails(userDetails, suite.ginContext).Return(nil).Times(1)
	suite.mockRedisStore.EXPECT().Delete(suite.ginContext.Request.Context(), userUUID).Return(nil).Times(1)

	err := suite.userRegistrationService.InitiateRegistration(registrationRequest, suite.ginContext)
	suite.Nil(err)
}

func (suite *UserRegistrationServiceTest) TestInitiateRegistration_WhenUnableToSaveUser() {
	userUUID := fmt.Sprintf("%s", uuid.New())
	registrationRequest := request.InitiateRegistrationRequest{
		Email:    "encrypted-username",
		Password: "encrypted-password",
		Username: "user123",
		UUID:     userUUID,
	}

	userDetails := db.SaveUserDetails{
		UUID:     userUUID,
		Username: registrationRequest.Username,
		Email:    registrationRequest.Email,
		Password: "hashed-password",
		IsActive: true,
	}
	suite.mockCryptoUtil.EXPECT().Decipher(suite.ginContext, registrationRequest.Password).Return("decrypted-password", nil).Times(1)
	suite.mockUserDetailsRepository.EXPECT().IsUserNameAndEmailAvailable(registrationRequest.Username, registrationRequest.Email, suite.ginContext).Return(false, nil).Times(1)
	suite.mockHashUtil.EXPECT().GenerateBcryptHash("decrypted-password").Return("hashed-password", nil).Times(1)
	suite.mockUserDetailsRepository.EXPECT().SaveUserDetails(userDetails, suite.ginContext).Return(errors.New("something went wrong")).Times(1)

	err := suite.userRegistrationService.InitiateRegistration(registrationRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(&constants.InternalServerError, err)
}

func (suite *UserRegistrationServiceTest) TestInitiateRegistration_WhenHashPasswordFails() {
	userUUID := fmt.Sprintf("%s", uuid.New())
	registrationRequest := request.InitiateRegistrationRequest{
		Email:    "encrypted-username",
		Password: "encrypted-password",
		Username: "user123",
		UUID:     userUUID,
	}

	userDetails := db.SaveUserDetails{
		UUID:     userUUID,
		Username: registrationRequest.Username,
		Email:    registrationRequest.Email,
		Password: "hashed-password",
		IsActive: true,
	}
	suite.mockCryptoUtil.EXPECT().Decipher(suite.ginContext, registrationRequest.Password).Return("decrypted-password", nil).Times(1)
	suite.mockUserDetailsRepository.EXPECT().IsUserNameAndEmailAvailable(registrationRequest.Username, registrationRequest.Email, suite.ginContext).Return(false, nil).Times(1)
	suite.mockHashUtil.EXPECT().GenerateBcryptHash("decrypted-password").Return("", errors.New("something went wrong")).Times(1)
	suite.mockUserDetailsRepository.EXPECT().SaveUserDetails(userDetails, suite.ginContext).Return(nil).Times(0)
	suite.mockRedisStore.EXPECT().Delete(suite.ginContext.Request.Context(), userUUID).Return(nil).Times(0)

	err := suite.userRegistrationService.InitiateRegistration(registrationRequest, suite.ginContext)
	suite.NotNil(err)

	suite.Equal(constants.InternalServerError.Error(), err.Error())
}

func (suite *UserRegistrationServiceTest) TestIsEmailRegistered_WhenDbReturnUserExists() {
	availabilityRequest := request.EmailAvailabilityRequest{Email: "dummy@gmail.com"}

	suite.mockUserDetailsRepository.EXPECT().IsEmailAvailable(availabilityRequest.Email, suite.ginContext).Return(true, nil).Times(1)

	actualResponse, err := suite.userRegistrationService.IsEmailRegistered(availabilityRequest, suite.ginContext)
	suite.Nil(err)
	suite.Equal(response.EmailAvailabilityResponse{IsAvailable: true}, actualResponse)
}

func (suite *UserRegistrationServiceTest) TestIsEmailRegistered_WhenDbReturnUserNotExists() {
	availabilityRequest := request.EmailAvailabilityRequest{Email: "dummy@gmail.com"}

	suite.mockUserDetailsRepository.EXPECT().IsEmailAvailable(availabilityRequest.Email, suite.ginContext).Return(false, nil).Times(1)

	actualResponse, err := suite.userRegistrationService.IsEmailRegistered(availabilityRequest, suite.ginContext)
	suite.Nil(err)
	suite.Equal(response.EmailAvailabilityResponse{IsAvailable: false}, actualResponse)
}

func (suite *UserRegistrationServiceTest) TestIsEmailRegistered_WhenDbReturnError() {
	availabilityRequest := request.EmailAvailabilityRequest{Email: "dummy@gmail.com"}

	suite.mockUserDetailsRepository.EXPECT().IsEmailAvailable(availabilityRequest.Email, suite.ginContext).Return(false, errors.New("something went wrong")).Times(1)

	actualResponse, err := suite.userRegistrationService.IsEmailRegistered(availabilityRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(response.EmailAvailabilityResponse{IsAvailable: false}, actualResponse)
	suite.Equal(&constants.IDPServiceFailureError, err)
}

func (suite *UserRegistrationServiceTest) TestIsUsernameRegistered_WhenDbReturnUserExists() {
	availabilityRequest := request.UsernameAvailabilityRequest{Username: "dummy-user"}

	suite.mockUserDetailsRepository.EXPECT().IsUserNameAvailable(availabilityRequest.Username, suite.ginContext).Return(true, nil).Times(1)

	actualResponse, err := suite.userRegistrationService.IsUsernameRegistered(availabilityRequest, suite.ginContext)
	suite.Nil(err)
	suite.Equal(response.UsernameAvailabilityResponse{IsAvailable: true}, actualResponse)
}

func (suite *UserRegistrationServiceTest) TestIsUsernameRegistered_WhenDbReturnUserNotExists() {
	availabilityRequest := request.UsernameAvailabilityRequest{Username: "dummy-user"}

	suite.mockUserDetailsRepository.EXPECT().IsUserNameAvailable(availabilityRequest.Username, suite.ginContext).Return(false, nil).Times(1)

	actualResponse, err := suite.userRegistrationService.IsUsernameRegistered(availabilityRequest, suite.ginContext)
	suite.Nil(err)
	suite.Equal(response.UsernameAvailabilityResponse{IsAvailable: false}, actualResponse)
}

func (suite *UserRegistrationServiceTest) TestIsUsernameRegistered_WhenDbReturnError() {
	availabilityRequest := request.UsernameAvailabilityRequest{Username: "dummy-user"}

	suite.mockUserDetailsRepository.EXPECT().IsUserNameAvailable(availabilityRequest.Username, suite.ginContext).Return(false, errors.New("something went wrong")).Times(1)

	actualResponse, err := suite.userRegistrationService.IsUsernameRegistered(availabilityRequest, suite.ginContext)
	suite.NotNil(err)
	suite.Equal(response.UsernameAvailabilityResponse{IsAvailable: false}, actualResponse)
	suite.Equal(&constants.IDPServiceFailureError, err)
}
