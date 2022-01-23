package service

// mockgen -source=service/login_service.go -destination=mocks/mock_login_service.go -package=mocks
import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/alert/email"
	"github.com/gola-glitch/gola-utils/alert/email/models"
	"github.com/gola-glitch/gola-utils/crypto"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/mask_util"
	"github.com/gola-glitch/gola-utils/redis_util"
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
	ForgetPassword(ctx context.Context, forgetPasswordRequest request.ForgetPassword) error
}

type loginService struct {
	userDetailsRepository repository.UserDetailsRepository
	cryptoUtils           crypto.CryptoUtil
	authenticator         AuthenticatorService
	loginHandler          login.OauthLoginHandler
	emailUtil             email.Util
	configData            *configuration.ConfigData
	store                 redis_util.RedisStore
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

func (service loginService) ForgetPassword(ctx context.Context, forgetPasswordRequest request.ForgetPassword) error {
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

func NewLoginService(detailsRepository repository.UserDetailsRepository, util crypto.CryptoUtil, authService AuthenticatorService, handler login.OauthLoginHandler, emailUtil email.Util, data *configuration.ConfigData, store redis_util.RedisStore) LoginService {
	return loginService{
		userDetailsRepository: detailsRepository,
		cryptoUtils:           util,
		authenticator:         authService,
		loginHandler:          handler,
		emailUtil:             emailUtil,
		configData:            data,
		store:                 store,
	}
}
