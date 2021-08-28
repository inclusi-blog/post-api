package service

// mockgen -source=service/login_service.go -destination=mocks/mock_login_service.go -package=mocks
import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/crypto"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/mask_util"
	"post-api/idp/constants"
	"post-api/idp/handlers/login"
	"post-api/idp/models/request"
	"post-api/idp/repository"
)

type LoginService interface {
	LoginWithEmailAndPassword(request request.UserLoginRequest, ctx *gin.Context) (interface{}, *golaerror.Error)
}

type loginService struct {
	userDetailsRepository repository.UserDetailsRepository
	cryptoUtils           crypto.CryptoUtil
	authenticator         AuthenticatorService
	loginHandler          login.OauthLoginHandler
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

func NewLoginService(detailsRepository repository.UserDetailsRepository, util crypto.CryptoUtil, authService AuthenticatorService, handler login.OauthLoginHandler) LoginService {
	return loginService{
		userDetailsRepository: detailsRepository,
		cryptoUtils:           util,
		authenticator:         authService,
		loginHandler:          handler,
	}
}
