package service

// mockgen -source=service/user_registration_service.go -destination=mocks/mock_user_registration_service.go -package=mocks
import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/crypto"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/redis_util"
	"post-api/idp/constants"
	"post-api/idp/models/db"
	"post-api/idp/models/request"
	"post-api/idp/models/response"
	"post-api/idp/repository"
	util "post-api/idp/utils"
)

type UserRegistrationService interface {
	InitiateRegistration(request request.InitiateRegistrationRequest, ctx *gin.Context) *golaerror.Error
	IsEmailRegistered(availabilityRequest request.EmailAvailabilityRequest, ctx *gin.Context) (response.EmailAvailabilityResponse, *golaerror.Error)
	IsUsernameRegistered(ctx context.Context, availabilityRequest request.UsernameAvailabilityRequest) (response.UsernameAvailabilityResponse, *golaerror.Error)
}

type userRegistrationService struct {
	repository  repository.UserDetailsRepository
	cryptoUtils crypto.CryptoUtil
	redis       redis_util.RedisStore
	hashUtil    util.HashUtil
}

func (service userRegistrationService) IsUsernameRegistered(ctx context.Context, availabilityRequest request.UsernameAvailabilityRequest) (response.UsernameAvailabilityResponse, *golaerror.Error) {
	log := logging.GetLogger(ctx).WithField("class", "UserRegistrationService").WithField("method", "IsUsernameRegistered")

	username := availabilityRequest.Username
	log.Infof("Calling repository to check username availability for username %v", username)
	isAvailable, err := service.repository.IsUserNameAvailable(availabilityRequest.Username, ctx)

	if err != nil {
		log.Errorf("Error occurred while fetching username availability for username %v .%v", username, err)
		return response.UsernameAvailabilityResponse{IsAvailable: false}, &constants.IDPServiceFailureError
	}

	availabilityResponse := response.UsernameAvailabilityResponse{
		IsAvailable: isAvailable,
	}

	log.Infof("Username not available for username %v", username)
	return availabilityResponse, nil
}

func (service userRegistrationService) IsEmailRegistered(availabilityRequest request.EmailAvailabilityRequest, ctx *gin.Context) (response.EmailAvailabilityResponse, *golaerror.Error) {
	log := logging.GetLogger(ctx).WithField("class", "UserRegistrationService").WithField("method", "IsEmailRegistered")

	email := availabilityRequest.Email
	log.Infof("Calling db to check user email existence for email %v", email)
	available, err := service.repository.IsEmailAvailable(email, ctx)

	if err != nil {
		log.Errorf("Error occurred while fetching user availability for email %v .%v", email, err)
		return response.EmailAvailabilityResponse{}, &constants.IDPServiceFailureError
	}

	log.Infof("Successfully fetched email existence from repository for email %v", email)
	return response.EmailAvailabilityResponse{IsAvailable: available}, nil
}

func (service userRegistrationService) InitiateRegistration(request request.InitiateRegistrationRequest, ctx *gin.Context) *golaerror.Error {
	log := logging.GetLogger(ctx).WithField("class", "UserRegistrationService").WithField("method", "InitiateRegistration")
	log.Info("Deciphering user registration password")
	email := request.Email

	decryptedPassword, err := service.cryptoUtils.Decipher(ctx, request.Password)

	if err != nil {
		log.Errorf("Unable to encrypt password while registering user  %v .%v", email, err)
		return &constants.InternalServerError
	}

	log.Infof("Successfully deciphered the password for user email %v", email)
	log.Infof("User password deciphered for user email %v", email)
	log.Infof("Calling user details to check user existence for user email %v", email)

	username, err := service.repository.GenerateUsername(ctx, email)
	if err != nil {
		log.Errorf("unable to generate username %v", err)
		return &constants.InternalServerError
	}

	isAvailable, err := service.repository.IsUserNameAndEmailAvailable(username, email, ctx)

	if err != nil {
		log.Infof("Error occurred while fetching user existence for email %v .%v", email, err)
		return &constants.InternalServerError
	}

	passwordHash, err := service.hashUtil.GenerateBcryptHash(decryptedPassword)

	if err != nil {
		log.Errorf("Unable to hash decrypted password %v", err)
		return &constants.InternalServerError
	}

	// TODO : Invert the condition and move to another flow when playing existence flow
	if !isAvailable {
		userRegistrationDetails := db.SaveUserDetails{
			ID:       request.Id,
			Username: username,
			Email:    request.Email,
			Password: passwordHash,
			IsActive: true,
		}

		log.Infof("Making repository call to save user details for email %v", email)
		err := service.repository.SaveUserDetails(userRegistrationDetails, ctx)

		if err != nil {
			log.Errorf("Error occurred while inserting user details for email %v .%v", email, err)
			return &constants.InternalServerError
		}

		err = service.redis.Delete(ctx.Request.Context(), userRegistrationDetails.ID.String())

		if err != nil {
			log.Errorf("Unable to delete activation hash in cache %v", err)
			return &constants.RegistrationRetryError
		}
	}

	return nil
}

func NewUserRegistrationService(detailsRepository repository.UserDetailsRepository, util crypto.CryptoUtil, store redis_util.RedisStore, hashUtil util.HashUtil) UserRegistrationService {
	return userRegistrationService{
		repository:  detailsRepository,
		cryptoUtils: util,
		redis:       store,
		hashUtil:    hashUtil,
	}
}
