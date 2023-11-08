package init

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/inclusi-blog/gola-utils/alert/email"
	"github.com/inclusi-blog/gola-utils/crypto"
	"github.com/inclusi-blog/gola-utils/logging"
	oauth2 "github.com/inclusi-blog/gola-utils/oauth"
	"github.com/inclusi-blog/gola-utils/redis_util"
	"github.com/jmoiron/sqlx"
	"log"
	"post-api/configuration"
	"post-api/helper"
	idpController "post-api/idp/controller"
	"post-api/idp/handlers/login"
	idpRepository "post-api/idp/repository"
	idpService "post-api/idp/service"
	idpUtil "post-api/idp/utils"
	commonService "post-api/service"
	storyController "post-api/story/controller"
	"post-api/story/repository"
	"post-api/story/service"
	"post-api/story/utils"
	userProfileController "post-api/user-profile/controller"
	userProfileRepository "post-api/user-profile/repository"
	userProfileService "post-api/user-profile/service"
	"strings"
)

var (
	draftController          storyController.DraftController
	interestsController      storyController.InterestsController
	postController           storyController.PostController
	registrationController   idpController.RegistrationController
	loginController          idpController.LoginController
	tokenController          idpController.TokenController
	profileController        userProfileController.UserProfileController
	registrationCacheService idpService.RegistrationCacheService
	userDetailsController    idpController.UserDetailsController
	reportController         storyController.ReportController
)

func Objects(db *sqlx.DB, configData *configuration.ConfigData, aws *session.Session) {
	logger := logging.GetLogger(context.TODO())
	redisPassword, err := getRedisPassword(aws, configData)
	if err != nil {
		log.Fatal(err)
	}
	configData.RedisStoreConfig.Password = redisPassword
	redisClient, redisError := redis_util.NewRedisClientWith(configData.RedisStoreConfig)
	if redisError != nil {
		logger.Errorf("Error occurred while initializing redis cache %v", redisError)
	}
	awsServices := commonService.NewAwsServices(aws, configData)

	postValidator := utils.NewPostValidator(configData)
	interestsRepository := repository.NewInterestRepository(db)
	interestsService := service.NewInterestsService(interestsRepository)
	interestsController = storyController.NewInterestsController(interestsService)
	manager := helper.NewTransactionManager(db)
	draftRepository := repository.NewDraftRepository(db)
	draftService := service.NewDraftService(draftRepository, interestsRepository, postValidator, awsServices)
	draftController = storyController.NewDraftController(draftService, awsServices)
	postRepository := repository.NewPostsRepository(db)
	previewPostRepository := repository.NewAbstractPostRepository(db)
	postService := service.NewPostService(postRepository, draftRepository, postValidator, previewPostRepository, interestsRepository, manager, awsServices)
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
	loginService := idpService.NewLoginService(detailsRepository, util, authenticatorService, oauthHandler, emailUtil, configData, redisClient, uuidGenerator, hashUtil)
	loginController = idpController.NewLoginController(loginService, oauthHandler)
	tokenController = idpController.NewTokenController(oauthHandler, configData.AllowInsecureCookies)

	profileRepository := userProfileRepository.NewProfileRepository(db)
	profileService := userProfileService.NewProfileService(profileRepository, awsServices)

	userInterestsRepository := userProfileRepository.NewUserInterestsRepository(db)
	userInterestsService := userProfileService.NewUserInterestsService(userInterestsRepository, awsServices)
	profileController = userProfileController.NewUserProfileController(userInterestsService, postService, profileService, awsServices)

	userDetailsService := idpService.NewUserDetailsService(detailsRepository, userRegistrationService)
	userDetailsController = idpController.NewUserDetailsController(userDetailsService, awsServices)

	reportRepository := repository.NewReportRepository(db)
	reportService := service.NewReportService(reportRepository)
	reportController = storyController.NewReportController(reportService)
}

func getRedisPassword(awsSession *session.Session, config *configuration.ConfigData) (string, error) {
	manager := secretsmanager.New(awsSession, nil)
	value, err := manager.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(strings.ToLower(config.Environment + "/redis/password")),
	})
	if err != nil {
		log.Println("unable to get db password from aws secrets ", err)
		return "", err
	}
	secretString := []byte(*value.SecretString)
	var data map[string]string
	err = json.Unmarshal(secretString, &data)
	if err != nil {
		log.Println("unable to unmarshal db password from aws secrets ", err)
		return "", err
	}

	return data[config.RedisPasswordKey], nil
}
