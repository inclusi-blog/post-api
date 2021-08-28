package service

import (
	"context"
	"errors"
	"github.com/gola-glitch/gola-utils/alert/email/models"
	mocksUtil "github.com/gola-glitch/gola-utils/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"post-api/configuration"
	"post-api/idp/constants"
	"post-api/idp/mocks"
	"post-api/idp/models/request"
	"testing"
)

type RegistrationCacheServiceTest struct {
	suite.Suite
	mockController           *gomock.Controller
	goContext                context.Context
	mockRedisStore           *mocks.MockRedisStore
	mockUUIDGenerator        *mocks.MockUUIDGenerator
	configData               *configuration.ConfigData
	mockEmailUtil            *mocksUtil.MockUtil
	registrationCacheService RegistrationCacheService
}

func TestRegistrationCacheServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RegistrationCacheServiceTest))
}

func (suite *RegistrationCacheServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	suite.mockRedisStore = mocks.NewMockRedisStore(suite.mockController)
	suite.configData = &configuration.ConfigData{
		Email: configuration.Email{
			GatewayURL:    "http://localhost:8083/api/ccg/v1/email/send",
			DefaultSender: "noreply@gola.xyz",
			TemplatePaths: configuration.TemplatesPaths{
				NewUserActivation: "../assets/email_templates/new_user_activation.html",
			},
		},
		ActivationCallback: "http://localhost:3000/m/callback/email",
	}
	suite.mockEmailUtil = mocksUtil.NewMockUtil(suite.mockController)
	suite.mockUUIDGenerator = mocks.NewMockUUIDGenerator(suite.mockController)
	suite.registrationCacheService = NewRegistrationCacheService(suite.mockRedisStore, suite.mockUUIDGenerator, suite.configData, suite.mockEmailUtil)
}

func (suite *RegistrationCacheServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *RegistrationCacheServiceTest) TestSaveUserDetailsInCache_WhenSuccess() {
	registrationRequest := request.NewRegistrationRequest{
		Email:    "dummy@email.com",
		Password: "dummy-encrypted-password",
		Username: "dummy user",
	}

	suite.mockUUIDGenerator.EXPECT().Generate().Return("some-value").Times(1)
	initiateRegistrationRequest := request.InitiateRegistrationRequest{
		Email:    registrationRequest.Email,
		Password: registrationRequest.Password,
		Username: registrationRequest.Username,
		UUID:     "some-value",
	}

	details := models.EmailDetails{
		From:    "noreply@gola.xyz",
		To:      []string{"dummy@email.com"},
		Subject: constants.VERIFY_EMAIL,
		Content: constants.NewUserActivation,
	}
	suite.mockEmailUtil.EXPECT().SendWithContext(suite.goContext, details, true).Return(nil).Times(1)
	suite.mockRedisStore.EXPECT().Set(suite.goContext, "some-value", initiateRegistrationRequest, 120).Return(nil).Times(1)
	err := suite.registrationCacheService.SaveUserDetailsInCache(registrationRequest, suite.goContext)
	suite.Nil(err)
}

func (suite *RegistrationCacheServiceTest) TestSaveUserDetailsInCache_WhenUnableToSetUserDetails() {
	registrationRequest := request.NewRegistrationRequest{
		Email:    "dummy@email.com",
		Password: "dummy-encrypted-password",
		Username: "dummy user",
	}

	suite.mockUUIDGenerator.EXPECT().Generate().Return("some-value").Times(1)
	initiateRegistrationRequest := request.InitiateRegistrationRequest{
		Email:    registrationRequest.Email,
		Password: registrationRequest.Password,
		Username: registrationRequest.Username,
		UUID:     "some-value",
	}
	suite.mockRedisStore.EXPECT().Set(suite.goContext, "some-value", initiateRegistrationRequest, 120).Return(errors.New("something went wrong")).Times(1)
	err := suite.registrationCacheService.SaveUserDetailsInCache(registrationRequest, suite.goContext)
	suite.NotNil(err)
	suite.Equal(&constants.IDPServiceFailureError, err)
}

func (suite *RegistrationCacheServiceTest) TestGetUserDetailsFromCache_WhenSuccess() {
	var initiateRegistration request.InitiateRegistrationRequest
	expectedRegistrationRequest := request.InitiateRegistrationRequest{
		Email:    "dummy@gmail.com",
		Password: "encrypted-password",
		Username: "selected-username",
		UUID:     "some-uuid",
	}
	suite.mockRedisStore.EXPECT().Get(suite.goContext, "some-hash", &initiateRegistration).Do(func(ctx context.Context, uuid string, destination *request.InitiateRegistrationRequest) {
		destination.Email = "dummy@gmail.com"
		destination.Password = "encrypted-password"
		destination.Username = "selected-username"
		destination.UUID = "some-uuid"
	}).Return(nil).Times(1)
	registrationRequest, err := suite.registrationCacheService.GetUserDetailsFromCache("some-hash", suite.goContext)
	suite.Nil(err)
	suite.Equal(expectedRegistrationRequest, registrationRequest)
}

func (suite *RegistrationCacheServiceTest) TestGetUserDetailsFromCache_WhenUnableToGetUserDetailsFromCache() {
	var initiateRegistration request.InitiateRegistrationRequest
	suite.mockRedisStore.EXPECT().Get(suite.goContext, "some-hash", &initiateRegistration).Return(errors.New("something went wrong")).Times(1)
	registrationRequest, err := suite.registrationCacheService.GetUserDetailsFromCache("some-hash", suite.goContext)
	suite.NotNil(err)
	expectedRegistrationRequest := request.InitiateRegistrationRequest{}
	suite.Equal(expectedRegistrationRequest, registrationRequest)
}

func (suite *RegistrationCacheServiceTest) TestSaveUserDetailsInCache_WhenUnableToSendEmail() {
	registrationRequest := request.NewRegistrationRequest{
		Email:    "dummy@email.com",
		Password: "dummy-encrypted-password",
		Username: "dummy user",
	}

	suite.mockUUIDGenerator.EXPECT().Generate().Return("some-value").Times(1)
	initiateRegistrationRequest := request.InitiateRegistrationRequest{
		Email:    registrationRequest.Email,
		Password: registrationRequest.Password,
		Username: registrationRequest.Username,
		UUID:     "some-value",
	}

	details := models.EmailDetails{
		From:    "noreply@gola.xyz",
		To:      []string{"dummy@email.com"},
		Subject: constants.VERIFY_EMAIL,
		Content: constants.NewUserActivation,
	}
	suite.mockEmailUtil.EXPECT().SendWithContext(suite.goContext, details, true).Return(&constants.IDPServiceFailureError).Times(1)
	suite.mockRedisStore.EXPECT().Set(suite.goContext, "some-value", initiateRegistrationRequest, 120).Return(nil).Times(1)
	err := suite.registrationCacheService.SaveUserDetailsInCache(registrationRequest, suite.goContext)
	suite.NotNil(err)
	suite.Equal(&constants.IDPServiceFailureError, err)
}
