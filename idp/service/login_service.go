package service

// mockgen -source=service/login_service.go -destination=mocks/mock_login_service.go -package=mocks
import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/alert/email"
	"github.com/inclusi-blog/gola-utils/alert/email/models"
	"github.com/inclusi-blog/gola-utils/crypto"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/inclusi-blog/gola-utils/mask_util"
	"github.com/inclusi-blog/gola-utils/redis_util"
	"github.com/segmentio/ksuid"
	"post-api/configuration"
	"post-api/idp/constants"
	"post-api/idp/handlers/login"
	"post-api/idp/models/request"
	"post-api/idp/repository"
	idputil "post-api/idp/utils"
)

type LoginService interface {
	LoginWithEmailAndPassword(request request.UserLoginRequest, ctx *gin.Context) (interface{}, *golaerror.Error)
	ForgetPassword(ctx context.Context, forgetPasswordRequest request.ForgetPassword) *golaerror.Error
	CanResetPassword(ctx context.Context, uniqueID string) (string, *golaerror.Error)
	ResetPassword(ctx context.Context, encryptedPassword, uniqueUUID string) *golaerror.Error
}

type loginService struct {
	userDetailsRepository repository.UserDetailsRepository
	cryptoUtils           crypto.Util
	authenticator         AuthenticatorService
	loginHandler          login.OauthLoginHandler
	uuid                  idputil.UUIDGenerator
	emailUtil             email.Util
	configData            *configuration.ConfigData
	store                 redis_util.RedisStore
	hashUtil              idputil.HashUtil
}

func (service loginService) LoginWithEmailAndPassword(request request.UserLoginRequest, ctx *gin.Context) (interface{}, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "LoginService").WithField("method", "LoginWithEmailAndPassword")
	isEmailAvailable, err := service.userDetailsRepository.IsEmailAvailable(request.Email, ctx)

	if err != nil {
		logger.Errorf("Error occurred while fetching user maskedEmail existence %v", err)
		return nil, &constants.InternalServerError
	}

	maskedEmail := mask_util.MaskEmail(ctx, request.Email)
	if !isEmailAvailable {
		logger.Infof("User not found in gola for maskedEmail %v", maskedEmail)
		return nil, &constants.UserNotFoundError
	}

	plainTextPassword, err := service.cryptoUtils.Decipher(ctx, request.Password)

	if err != nil {
		logger.Errorf("Password decipher failed for maskedEmail %v .%v", maskedEmail, err)
		return nil, &constants.InternalServerError
	}

	profile, err := service.userDetailsRepository.GetUserProfile(request.Email, ctx)

	if err != nil {
		logger.Errorf("Error occurred while fetching user profile details for maskedEmail %v .%v", maskedEmail, err)
		return nil, &constants.InternalServerError
	}

	authenticationErr := service.authenticator.Authenticate(ctx, plainTextPassword, request.Email)

	if authenticationErr != nil {
		logger.Errorf("Error occurred while authenticating user for maskedEmail %v", maskedEmail)
		return nil, authenticationErr
	}

	loginResponse, oauthErr := service.loginHandler.AcceptLogin(ctx, request.LoginChallenge, profile)

	if oauthErr != nil {
		logger.Errorf("Error in accepting login %v", oauthErr)
		return nil, oauthErr
	}

	logger.Infof("Password login successful for user email %v", maskedEmail)

	return loginResponse, nil
}

func (service loginService) ForgetPassword(ctx context.Context, forgetPasswordRequest request.ForgetPassword) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "LoginService").WithField("method", "ForgetPassword")

	userEmail := forgetPasswordRequest.Email
	isEmailAvailable, err := service.userDetailsRepository.IsEmailAvailable(userEmail, ctx)

	if err != nil {
		logger.Errorf("Error occurred while fetching user maskedEmail existence %v", err)
		return &constants.InternalServerError
	}

	maskedEmail := mask_util.MaskEmail(ctx, userEmail)
	if !isEmailAvailable {
		logger.Infof("User not found in gola for maskedEmail %v", maskedEmail)
		return &constants.UserNotFoundError
	}

	logger.Infof("user found with userEmail %v", maskedEmail)

	type ForgetPassword struct {
		ResetURL string
	}

	generatedID := ksuid.New().String()
	forgetPassword := ForgetPassword{ResetURL: service.configData.PasswordResetCallback + generatedID}

	err = service.store.Set(ctx, generatedID, userEmail, 120)
	if err != nil {
		logger.Errorf("Error occurred while saving user password reset details in cache for userEmail %v . %v", userEmail, err)
		return &constants.IDPServiceFailureError
	}

	emailContent, parseErr := idputil.ParseTemplate(ctx, service.configData.Email.TemplatePaths.ForgetPassword, forgetPassword)
	if parseErr != nil {
		logger.Errorf("unable to parse forget password email template %v", parseErr)
		return &constants.InternalServerError
	}

	emailDetails := models.EmailDetails{
		From:    service.configData.Email.DefaultSender,
		To:      []string{userEmail},
		Subject: constants.VerifyEmail,
		Content: emailContent,
	}

	logger.Infof("Sending password reset link to user userEmail %v", userEmail)
	emailErr := service.emailUtil.SendWithContext(ctx, emailDetails, true)

	if emailErr != nil {
		logger.Errorf("Unable to send userEmail to user %v .%v", userEmail, emailErr)
		return emailErr
	}

	logger.Infof("Successfully saved user email details in cache for userEmail %v", userEmail)

	return nil
}

func (service loginService) CanResetPassword(ctx context.Context, uniqueID string) (string, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "LoginService").WithField("method", "CanResetPassword")
	var userEmail string
	err := service.store.Get(ctx, uniqueID, &userEmail)
	if err != nil {
		logger.Errorf("unable to find key %v", err)
		return "", &constants.UnauthorisedRequestError
	}

	isEmailAvailable, err := service.userDetailsRepository.IsEmailAvailable(userEmail, ctx)
	if err != nil {
		logger.Errorf("unable to find user email %v", err)
		return "", &constants.UnauthorisedRequestError
	}

	maskedEmail := mask_util.MaskEmail(ctx, userEmail)
	if !isEmailAvailable {
		logger.Infof("User not found in gola for maskedEmail %v", maskedEmail)
		return "", &constants.UnauthorisedRequestError
	}

	generatedUUID := service.uuid.Generate()
	uniqueUUID := generatedUUID.String()
	err = service.store.Set(ctx, uniqueUUID, &userEmail, 15)
	if err != nil {
		logger.Errorf("unable to set unique generated id %v", err)
		return "", &constants.UnableToResetPasswordError
	}

	go func(routineCtx context.Context) {
		err := service.store.Delete(routineCtx, uniqueID)
		if err != nil {
			logger.Errorf("unable to delete old identifier %v", err)
		}
	}(ctx)

	return uniqueUUID, nil
}

func (service loginService) ResetPassword(ctx context.Context, encryptedPassword, uniqueUUID string) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "LoginService").WithField("method", "CanResetPassword")
	var userEmail string
	err := service.store.Get(ctx, uniqueUUID, &userEmail)
	if err != nil {
		logger.Errorf("unable to find identifier %v", err)
		return &constants.UnauthorisedRequestError
	}
	go func(verifier string) {
		err := service.store.Delete(ctx, verifier)
		if err != nil {
			logger.Errorf("unable to delete verifier %v", err)
		}
	}(uniqueUUID)

	isEmailAvailable, err := service.userDetailsRepository.IsEmailAvailable(userEmail, ctx)
	if err != nil {
		logger.Errorf("unable to find user email %v", err)
		return &constants.UnauthorisedRequestError
	}

	maskedEmail := mask_util.MaskEmail(ctx, userEmail)
	if !isEmailAvailable {
		logger.Infof("User not found in gola for maskedEmail %v", maskedEmail)
		return &constants.UnauthorisedRequestError
	}

	decryptedPassword, err := service.cryptoUtils.Decipher(ctx, encryptedPassword)

	if err != nil {
		logger.Errorf("Unable to encrypt password while registering user  %v .%v", userEmail, err)
		return &constants.InternalServerError
	}

	logger.Infof("Successfully deciphered the password for user email %v", userEmail)
	logger.Infof("User password deciphered for user email %v", userEmail)

	passwordHash, err := service.hashUtil.GenerateBcryptHash(decryptedPassword)

	if err != nil {
		logger.Errorf("Unable to hash decrypted password %v", err)
		return &constants.InternalServerError
	}

	err = service.userDetailsRepository.UpdatePassword(ctx, passwordHash, userEmail)
	if err != nil {
		logger.Errorf("unable to update password %v", err)
		return &constants.UnableToResetPasswordError
	}

	return nil
}

func NewLoginService(detailsRepository repository.UserDetailsRepository, util crypto.Util, authService AuthenticatorService, handler login.OauthLoginHandler, emailUtil email.Util, data *configuration.ConfigData, store redis_util.RedisStore, generator idputil.UUIDGenerator, hashUtil idputil.HashUtil) LoginService {
	return loginService{
		userDetailsRepository: detailsRepository,
		cryptoUtils:           util,
		authenticator:         authService,
		loginHandler:          handler,
		uuid:                  generator,
		emailUtil:             emailUtil,
		configData:            data,
		store:                 store,
		hashUtil:              hashUtil,
	}
}
