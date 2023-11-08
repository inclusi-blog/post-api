package service

// mockgen -source=service/authenticator_service.go -destination=mocks/mock_authenticator_service.go -package=mocks
import (
	"github.com/gin-gonic/gin"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/inclusi-blog/gola-utils/mask_util"
	"post-api/idp/constants"
	"post-api/idp/repository"
	util "post-api/idp/utils"
)

type AuthenticatorService interface {
	Authenticate(ctx *gin.Context, planTextPassword, email string) *golaerror.Error
}

type authenticatorService struct {
	hashUtils             util.HashUtil
	userDetailsRepository repository.UserDetailsRepository
}

func (service authenticatorService) Authenticate(ctx *gin.Context, planTextPassword, email string) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "AuthenticatorService").WithField("method", "Authenticate")
	hashedPassword, err := service.userDetailsRepository.GetPassword(email, ctx)

	maskEmail := mask_util.MaskEmail(ctx, email)
	if err != nil {
		logger.Errorf("Error occurred while fetching user credentials for email %v. %v", maskEmail, err)
		return &constants.InternalServerError
	}

	err = service.hashUtils.MatchBcryptHash(hashedPassword, planTextPassword)

	if err != nil {
		logger.Errorf("Invalid credentials for email %v .%v", maskEmail, err)
		return &constants.InvalidCredentialsError
	}

	logger.Infof("Successfully validated login request credentials for email %v", maskEmail)

	return nil
}

func NewAuthenticatorService(hashUtil util.HashUtil, detailsRepository repository.UserDetailsRepository) AuthenticatorService {
	return authenticatorService{
		hashUtils:             hashUtil,
		userDetailsRepository: detailsRepository,
	}
}
