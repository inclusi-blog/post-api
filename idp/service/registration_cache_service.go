package service

import (
	"context"
	"github.com/gola-glitch/gola-utils/alert/email"
	"github.com/gola-glitch/gola-utils/alert/email/models"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/redis_util"
	"post-api/configuration"
	"post-api/idp/constants"
	"post-api/idp/models/request"
	idputil "post-api/idp/utils"
)

type RegistrationCacheService interface {
	SaveUserDetailsInCache(request request.NewRegistrationRequest, ctx context.Context) *golaerror.Error
	GetUserDetailsFromCache(activationHash string, ctx context.Context) (request.InitiateRegistrationRequest, *golaerror.Error)
}

type registrationCacheService struct {
	store         redis_util.RedisStore
	uuidGenerator idputil.UUIDGenerator
	configData    *configuration.ConfigData
	emailUtil     email.Util
}

func (service registrationCacheService) GetUserDetailsFromCache(activationHash string, ctx context.Context) (request.InitiateRegistrationRequest, *golaerror.Error) {
	log := logging.GetLogger(ctx).WithField("class", "RegistrationCacheService").WithField("method", "GetUserDetailsFromCache")

	log.Infof("Fetching user registration request from cache %v", activationHash)

	var initRequest request.InitiateRegistrationRequest
	if redisErr := service.store.Get(ctx, activationHash, &initRequest); redisErr != nil {
		log.Errorf("User details expired or not available for user uuid %v", activationHash)
		return request.InitiateRegistrationRequest{}, &constants.ActivationLinkExpiredError
	}

	log.Info("Successfully fetched user registration")

	return initRequest, nil
}

func (service registrationCacheService) SaveUserDetailsInCache(newRequest request.NewRegistrationRequest, ctx context.Context) *golaerror.Error {
	log := logging.GetLogger(ctx).WithField("class", "RegistrationCacheService").WithField("method", "SaveUserDetailsInCache")
	generatedUUID := service.uuidGenerator.Generate()

	registrationRequest := request.InitiateRegistrationRequest{
		Email:    newRequest.Email,
		Password: newRequest.Password,
		Username: newRequest.Username,
		Id:       generatedUUID,
	}
	log.Infof("Saving new user registration request in cache %v for userEmail and user uuid %v", newRequest.Email, generatedUUID)

	err := service.store.Set(ctx, generatedUUID.String(), registrationRequest, 120)

	if err != nil {
		log.Errorf("Error occurred while saving user registration request in cache for userEmail %v . %v", newRequest.Email, err)
		return &constants.IDPServiceFailureError
	}

	type NewUserActivation struct {
		ActivationUrl string
	}

	// TODO Change the url when frontend route is played
	activation := NewUserActivation{ActivationUrl: service.configData.ActivationCallback + `?token=` + generatedUUID.String() + ``}
	emailContent, _ := idputil.ParseTemplate(ctx, service.configData.Email.TemplatePaths.NewUserActivation, activation)

	userEmail := registrationRequest.Email
	emailDetails := models.EmailDetails{
		From:    service.configData.Email.DefaultSender,
		To:      []string{userEmail},
		Subject: constants.VERIFY_EMAIL,
		Content: emailContent,
	}

	log.Infof("Sending activation link to user userEmail %v", userEmail)
	emailErr := service.emailUtil.SendWithContext(ctx, emailDetails, true)

	if emailErr != nil {
		log.Errorf("Unable to send userEmail to user %v .%v", userEmail, err)
		return emailErr
	}

	log.Infof("Successfully saved user registration request in cache for userEmail %v", userEmail)

	return nil
}

func NewRegistrationCacheService(store redis_util.RedisStore, generator idputil.UUIDGenerator, data *configuration.ConfigData, util email.Util) RegistrationCacheService {
	return registrationCacheService{
		store:         store,
		uuidGenerator: generator,
		configData:    data,
		emailUtil:     util,
	}
}
