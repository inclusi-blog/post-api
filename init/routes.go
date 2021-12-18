package init

import (
	"context"
	cacheMiddleware "github.com/gola-glitch/gola-utils/middleware/cache_control"
	tokenMiddleware "github.com/gola-glitch/gola-utils/middleware/introspection"
	"github.com/gola-glitch/gola-utils/oauth"
	"net/http"
	"post-api/configuration"
	"post-api/idp/middlewares"
	commonService "post-api/service"
	"post-api/story/constants"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	cors "github.com/gola-glitch/gola-utils/middleware/cors"
	"github.com/gola-glitch/gola-utils/middleware/request_response_trace"
	middleware "github.com/gola-glitch/gola-utils/middleware/session_trace"
	corsModel "github.com/gola-glitch/gola-utils/model"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func globalPanicMiddleware(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logging.GetLogger(ctx).Errorf("Error occurred: %+v, \n %s", err, string(debug.Stack()))
			serverError := constants.InternalServerError
			constants.RespondWithGolaError(ctx, serverError)
		}
	}()
	ctx.Next()
}

func tokenIntrospectionMiddleware(oauthBaseUrl string, oauthUtil oauth.Utils, data *configuration.ConfigData) gin.HandlerFunc {
	protectedUrlService := commonService.NewProtectedUrlService(data)
	introspectionMiddleware := tokenMiddleware.NewIntrospectionAndDecryptionMiddleware(protectedUrlService, oauthBaseUrl, oauthUtil)
	return introspectionMiddleware.TokenValidationMiddleware()
}

func RegisterRouter(router *gin.Engine, configData *configuration.ConfigData) {
	router.GET("api/post/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"Status": "Up",
		})
	})

	golaLoggerRegistry := logging.NewLoggerEntry()

	router.Use(logging.LoggingMiddleware(golaLoggerRegistry))

	logLevel := configData.LogLevel
	logger := logging.GetLogger(context.TODO())

	if logLevel != "" {
		logLevelInitErr := golaLoggerRegistry.SetLevel(logLevel)
		if logLevelInitErr != nil {
			logger.Warning("gola_logger.SetLevel failed. Default log level being used", logLevelInitErr.Error())
		}
	}

	router.Use(middleware.SessionTracingMiddleware)
	router.Use(request_response_trace.HttpRequestResponseTracingAllMiddlewareWithCustomHealthEndpoint("api/post/healthz"))

	corsConfig := corsModel.CorsConfig{
		AllowedOrigins: configData.AllowedOrigins,
	}

	router.Use(cors.CORSMiddleware(corsConfig))
	router.Use(globalPanicMiddleware)
	cacheControl := cacheMiddleware.NewCacheControlMiddleware([]string{})
	router.Use(cacheControl.StopCaching())

	oauthUtil := oauth.NewOauthUtils(configData.CryptoServiceURL)

	router.GET("api/post/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	idpRoute := router.Group("api/idp/v1")
	registrationGroup := idpRoute.Group("user")
	{
		registrationGroup.POST("/register", registrationController.NewRegistration)
		registrationGroup.GET("/activate/:activation_hash", middlewares.UserActionMiddleware(registrationCacheService), registrationController.ActivateUser)
		registrationGroup.POST("/emailAvailable", registrationController.IsEmailAvailable)
		registrationGroup.POST("/usernameAvailability", registrationController.IsUsernameAvailable)
	}

	loginGroup := idpRoute.Group("login")
	{
		loginGroup.POST("password", loginController.LoginByEmailAndPassword)
		loginGroup.GET("accept", loginController.AcceptChallenge)
	}

	consentGroup := idpRoute.Group("consent")
	{
		consentGroup.GET("", loginController.GrantConsent)
	}

	tokenGroup := idpRoute.Group("token")
	{
		tokenGroup.POST("/exchange", tokenController.ExchangeToken)
	}

	userDetailsGroup := idpRoute.Group("user-details", tokenIntrospectionMiddleware(configData.OauthUrl, oauthUtil, configData))
	{
		userDetailsGroup.PUT("", userDetailsController.UpdateUserDetails)
	}

	defaultRouterGroup := router.Group("api/post/v1")
	defaultRouterGroup.GET("/interests", interestsController.GetInterests)
	defaultRouterGroup.Use(tokenIntrospectionMiddleware(configData.OauthUrl, oauthUtil, configData))
	{
		draftGroup := defaultRouterGroup.Group("/draft")
		{
			draftGroup.POST("", draftController.CreateDraft)
			draftGroup.PUT("", draftController.SaveDraft)
			draftGroup.PUT("/tagline", draftController.SaveTagline)
			draftGroup.PUT("/interests", draftController.SaveInterests)
			draftGroup.GET("", draftController.GetDraft)
			draftGroup.POST("/get-all-draft", draftController.GetAllDraft)
			draftGroup.GET("/preview-draft/:draft_id", draftController.GetPreviewDraft)
			draftGroup.DELETE("id/:draft_id", draftController.DeleteDraft)
			draftGroup.GET("pre-sign/:draft_id", draftController.GetPreSignURLForDraftPreview)
			draftGroup.POST("preview/:draft_id/upload", draftController.UploadImageKey)
		}

		postGroup := defaultRouterGroup.Group("/post")
		{
			postGroup.POST("/publish", postController.PublishPost)
			postGroup.GET("/like", postController.Like)
			postGroup.GET("/unlike", postController.UnLike)
			postGroup.GET("/id/:post_id", postController.GetPost)
			postGroup.DELETE("/id/:post_id", postController.Delete)
			postGroup.POST("comment/:post_id", postController.Comment)
			postGroup.GET("/id/:post_id/comments", postController.GetComments)
			postGroup.GET("/id/:post_id/save", postController.SavePost)
			postGroup.GET("/id/:post_id/remove", postController.RemoveBookmark)
			postGroup.GET("/id/:post_id/viewed", postController.MarkAsViewed)
			postGroup.GET("/saved", postController.GetReadLaterPosts)
			postGroup.GET("/viewed", postController.GetReadPosts)
		}
	}

	interestGroup := defaultRouterGroup.Group("interests")
	{
		interestGroup.POST("/details", interestsController.GetInterestDetails)
		interestGroup.Use(tokenIntrospectionMiddleware(configData.OauthUrl, oauthUtil, configData))
		{
			interestGroup.GET("/posts/:interest_id", postController.GetPostsByInterests)
		}
	}

	userGroup := router.Group("api/user-profile/v1")
	noAuthUserprofile := router.Group("api/user-profile/v1")
	userGroup.Use(tokenIntrospectionMiddleware(configData.OauthUrl, oauthUtil, configData))
	{
		interests := userGroup.Group("/interests")
		{
			interests.GET("/followed", profileController.GetFollowedInterests)
			interests.POST("", profileController.FollowInterest)
			interests.DELETE("", profileController.UnFollowInterest)
			interests.GET("/explore", profileController.GetExploreInterests)
		}
		userBehaviourGroup := userGroup.Group("user")
		{
			userBehaviourGroup.GET(":user_id/follow", profileController.FollowUser)
			userBehaviourGroup.GET(":user_id/unfollow", profileController.UnFollowUser)
			userBehaviourGroup.GET(":user_id/block", profileController.BlockUser)
		}
		posts := userGroup.Group("posts")
		{
			posts.POST("", profileController.GetPublishedPosts)
		}
		profileGroup := userGroup.Group("profile")
		{
			noAuthUserprofile.GET("user/:user_id/avatar", profileController.ViewProfileAvatar)
			profileGroup.GET("/presign", userDetailsController.GetPreSignURLForProfilePic)
			profileGroup.GET("", profileController.GetDetails)
			profileGroup.POST("avatar/upload", userDetailsController.UploadImageKey)
		}
	}
}
