package init

import (
	"context"
	"github.com/gola-glitch/gola-utils/alert/email"
	"github.com/gola-glitch/gola-utils/crypto"
	"github.com/gola-glitch/gola-utils/logging"
	oauth2 "github.com/gola-glitch/gola-utils/oauth"
	"github.com/gola-glitch/gola-utils/redis_util"
	"github.com/jmoiron/sqlx"
	"post-api/configuration"
	"post-api/helper"
	idpController "post-api/idp/controller"
	"post-api/idp/handlers/login"
	idpRepository "post-api/idp/repository"
	idpService "post-api/idp/service"
	idpUtil "post-api/idp/utils"
	storyController "post-api/story/controller"
	"post-api/story/repository"
	"post-api/story/service"
	"post-api/story/utils"
	userProfileController "post-api/user-profile/controller"
	userProfileRepository "post-api/user-profile/repository"
	userProfileService "post-api/user-profile/service"
)

var (
	draftController          storyController.DraftController
	interestsController      storyController.InterestsController
	postController           storyController.PostController
	registrationController   idpController.RegistrationController
	loginController          idpController.LoginController
	tokenController          idpController.TokenController
	userInterestsController  userProfileController.UserInterestsController
	registrationCacheService idpService.RegistrationCacheService
)

func Objects(db *sqlx.DB, configData *configuration.ConfigData) {
	logger := logging.GetLogger(context.TODO())
	redisClient, redisError := redis_util.NewRedisClientWith(configData.RedisStoreConfig)
	if redisError != nil {
		logger.Errorf("Error occurred while initializing redis cache %v", redisError)
	}
	postValidator := utils.NewPostValidator(configData)
	interestsRepository := repository.NewInterestRepository(db)
	interestsService := service.NewInterestsService(interestsRepository)
	interestsController = storyController.NewInterestsController(interestsService)
	manager := helper.NewTransactionManager(db)
	draftRepository := repository.NewDraftRepository(db)
	draftService := service.NewDraftService(draftRepository, interestsRepository, postValidator)
	draftController = storyController.NewDraftController(draftService)
	postRepository := repository.NewPostsRepository(db)
	previewPostRepository := repository.NewAbstractPostRepository(db)
	postService := service.NewPostService(postRepository, draftRepository, postValidator, previewPostRepository, interestsRepository, manager)
	postController = storyController.NewPostController(postService)

	detailsRepository := idpRepository.NewUserDetailsRepository(db)
	util := crypto.NewCryptoUtil(configData.CryptoServiceURL)
	hashUtil := idpUtil.NewHashUtil()
	uuidGenerator := idpUtil.NewUUIDGenerator()
	emailUtil := email.NewEmailUtil(configData.Email.GatewayURL)
	userRegistrationService := idpService.NewUserRegistrationService(detailsRepository, util, redisClient, hashUtil)
	registrationCacheService = idpService.NewRegistrationCacheService(redisClient, uuidGenerator, configData, emailUtil)
	registrationController = idpController.NewRegistrationController(registrationCacheService, userRegistrationService)
	oauthUtils := oauth2.NewOauthUtils(configData.CryptoServiceURL)
	clockUtil := idpUtil.NewClock()

	authenticatorService := idpService.NewAuthenticatorService(hashUtil, detailsRepository)
	oauthHandler := login.NewOauthLoginHandler(requestBuilder, configData, oauthUtils, clockUtil)
	loginService := idpService.NewLoginService(detailsRepository, util, authenticatorService, oauthHandler)
	loginController = idpController.NewLoginController(loginService, oauthHandler)
	tokenController = idpController.NewTokenController(oauthHandler, configData.AllowInsecureCookies)

	userInterestsRepository := userProfileRepository.NewUserInterestsRepository(db)
	userInterestsService := userProfileService.NewUserInterestsService(userInterestsRepository)
	userInterestsController = userProfileController.NewUserInterestsController(userInterestsService, postService)
}
