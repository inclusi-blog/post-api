package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/idp/constants"
	"post-api/idp/mocks"
	"post-api/idp/models/request"
	"post-api/idp/models/response"
	"testing"
)

type RegistrationControllerTest struct {
	suite.Suite
	mockCtrl                     *gomock.Controller
	recorder                     *httptest.ResponseRecorder
	context                      *gin.Context
	mockRegistrationCacheService *mocks.MockRegistrationCacheService
	mockUserRegistrationService  *mocks.MockUserRegistrationService
	registrationController       RegistrationController
}

func (suite *RegistrationControllerTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.mockRegistrationCacheService = mocks.NewMockRegistrationCacheService(suite.mockCtrl)
	suite.mockUserRegistrationService = mocks.NewMockUserRegistrationService(suite.mockCtrl)
	suite.registrationController = NewRegistrationController(suite.mockRegistrationCacheService, suite.mockUserRegistrationService)
}

func (suite *RegistrationControllerTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestRegistrationControllerTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrationControllerTest))
}

func (suite *RegistrationControllerTest) TestNewRegistration_WhenSuccess() {
	requestBody := request.NewRegistrationRequest{
		Email:    "dummy@gmail.com",
		Password: "encrypted-password",
		Username: "dummy-user",
	}

	suite.mockRegistrationCacheService.EXPECT().SaveUserDetailsInCache(requestBody, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(requestBody)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/idp/user/register", bytes.NewBufferString(string(jsonBytes)))
	suite.registrationController.NewRegistration(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *RegistrationControllerTest) TestNewRegistration_WhenRequestIsInvalid() {
	requestBody := request.NewRegistrationRequest{
		Email:    "dummy@gmail.com",
		Password: "encrypted-password",
		Username: "dummy-user",
	}

	requestString := `{"email": "dummy@gmail.com", "password": "encrypted-password"}`
	suite.mockRegistrationCacheService.EXPECT().SaveUserDetailsInCache(requestBody, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/idp/user/register", bytes.NewBufferString(requestString))
	suite.registrationController.NewRegistration(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	bytesResponse, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(bytesResponse, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestNewRegistration_WhenCacheServiceThrowsError() {
	requestBody := request.NewRegistrationRequest{
		Email:    "dummy@gmail.com",
		Password: "encrypted-password",
		Username: "dummy-user",
	}

	suite.mockRegistrationCacheService.EXPECT().SaveUserDetailsInCache(requestBody, suite.context).Return(&constants.InternalServerError).Times(1)

	jsonBytes, err := json.Marshal(requestBody)
	suite.Nil(err)

	bytesResponse, err := json.Marshal(constants.InternalServerError)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/idp/user/register", bytes.NewBufferString(string(jsonBytes)))
	suite.registrationController.NewRegistration(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(bytesResponse, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestActivateUser_WhenValidRequest() {
	requestBody := request.InitiateRegistrationRequest{
		Email:    "dummy@gmail.com",
		Password: "encrypted-password",
		Username: "dummy-user",
		UUID:     "some-uuid",
	}

	suite.mockUserRegistrationService.EXPECT().InitiateRegistration(requestBody, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(requestBody)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/idp/user/activate/some-uuid", bytes.NewBufferString(string(jsonBytes)))
	suite.registrationController.ActivateUser(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *RegistrationControllerTest) TestActivateUser_WhenInvalidRequest() {
	requestBody := request.InitiateRegistrationRequest{}
	requestString := `{"email": "dummy@gmail.com", "password": "encrypted-password"}`
	suite.mockUserRegistrationService.EXPECT().InitiateRegistration(requestBody, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/idp/user/activate/some-uuid", bytes.NewBufferString(requestString))
	suite.registrationController.ActivateUser(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	marshalBytes, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(marshalBytes, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestActivateUser_WhenInitiateRegistrationThrowsError() {
	requestBody := request.InitiateRegistrationRequest{
		Email:    "dummy@gmail.com",
		Password: "encrypted-password",
		Username: "dummy-user",
		UUID:     "some-uuid",
	}

	suite.mockUserRegistrationService.EXPECT().InitiateRegistration(requestBody, suite.context).Return(&constants.InternalServerError).Times(1)

	jsonBytes, err := json.Marshal(requestBody)
	suite.Nil(err)

	bytesResponse, err := json.Marshal(constants.InternalServerError)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/idp/user/activate/some-uuid", bytes.NewBufferString(string(jsonBytes)))
	suite.registrationController.ActivateUser(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(bytesResponse, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestIsEmailAvailable_WhenValidRequest() {
	requestBody := request.EmailAvailabilityRequest{Email: "dummy@gmail.com"}

	expectedResponse := response.EmailAvailabilityResponse{IsAvailable: true}

	suite.mockUserRegistrationService.EXPECT().IsEmailRegistered(requestBody, suite.context).Return(response.EmailAvailabilityResponse{IsAvailable: true}, nil).Times(1)

	jsonBytes, err := json.Marshal(requestBody)

	suite.Nil(err)

	jsonResponseBytes, err := json.Marshal(expectedResponse)

	suite.Nil(err)
	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/idp/v1/user/emailAvailability", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)
	suite.registrationController.IsEmailAvailable(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(jsonResponseBytes, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestIsEmailAvailable_WhenInvalidRequest() {
	requestBody := request.EmailAvailabilityRequest{Email: "dummy@gmail.com"}

	suite.mockUserRegistrationService.EXPECT().IsEmailRegistered(requestBody, suite.context).Return(response.EmailAvailabilityResponse{IsAvailable: true}, nil).Times(0)

	jsonResponseBytes, err := json.Marshal(&constants.PayloadValidationError)

	suite.Nil(err)
	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/idp/v1/user/emailAvailability", bytes.NewBufferString(`{"email": ""}`))
	suite.Nil(err)
	suite.registrationController.IsEmailAvailable(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(jsonResponseBytes, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestIsEmailAvailable_WhenServiceFailure() {
	requestBody := request.EmailAvailabilityRequest{Email: "dummy@gmail.com"}

	suite.mockUserRegistrationService.EXPECT().IsEmailRegistered(requestBody, suite.context).Return(response.EmailAvailabilityResponse{}, &constants.IDPServiceFailureError).Times(1)

	jsonBytes, err := json.Marshal(requestBody)

	suite.Nil(err)

	jsonResponseBytes, err := json.Marshal(&constants.IDPServiceFailureError)

	suite.Nil(err)
	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/idp/v1/user/emailAvailability", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)
	suite.registrationController.IsEmailAvailable(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(jsonResponseBytes, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestIsUsernameAvailable_WhenValidRequest() {
	requestBody := request.UsernameAvailabilityRequest{Username: "someuser"}

	expectedResponse := response.UsernameAvailabilityResponse{IsAvailable: true}

	suite.mockUserRegistrationService.EXPECT().IsUsernameRegistered(requestBody, suite.context).Return(response.UsernameAvailabilityResponse{IsAvailable: true}, nil).Times(1)

	jsonBytes, err := json.Marshal(requestBody)

	suite.Nil(err)

	jsonResponseBytes, err := json.Marshal(expectedResponse)

	suite.Nil(err)
	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/idp/v1/user/usernameAvailability", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)
	suite.registrationController.IsUsernameAvailable(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(jsonResponseBytes, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestIsUsernameAvailable_WhenInvalidRequest() {
	requestBody := request.UsernameAvailabilityRequest{Username: "someuser"}

	suite.mockUserRegistrationService.EXPECT().IsUsernameRegistered(requestBody, suite.context).Return(response.UsernameAvailabilityResponse{IsAvailable: true}, nil).Times(0)

	jsonResponseBytes, err := json.Marshal(&constants.PayloadValidationError)

	suite.Nil(err)
	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/idp/v1/user/usernameAvailability", bytes.NewBufferString(`{"username": ""}`))
	suite.Nil(err)
	suite.registrationController.IsUsernameAvailable(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(jsonResponseBytes, suite.recorder.Body.Bytes())
}

func (suite *RegistrationControllerTest) TestIsUsernameAvailable_WhenServiceFailure() {
	requestBody := request.UsernameAvailabilityRequest{Username: "someuser"}

	suite.mockUserRegistrationService.EXPECT().IsUsernameRegistered(requestBody, suite.context).Return(response.UsernameAvailabilityResponse{}, &constants.IDPServiceFailureError).Times(1)

	jsonBytes, err := json.Marshal(requestBody)

	suite.Nil(err)

	jsonResponseBytes, err := json.Marshal(&constants.IDPServiceFailureError)

	suite.Nil(err)
	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/idp/v1/user/usernameAvailability", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)
	suite.registrationController.IsUsernameAvailable(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(jsonResponseBytes, suite.recorder.Body.Bytes())
}
