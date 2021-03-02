package init

import (
	"context"
	cacheControlMiddleware "github.com/gola-glitch/gola-utils/middleware/cache_control"
	tokenIntrospection "github.com/gola-glitch/gola-utils/middleware/introspection"
	"github.com/gola-glitch/gola-utils/oauth"
	"net/http"
	"post-api/configuration"
	"post-api/constants"
	service "post-api/service"
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
	protectedUrlService := service.NewProtectedUrlService(data)
	introspectionMiddleware := tokenIntrospection.NewIntrospectionAndDecryptionMiddleware(protectedUrlService, oauthBaseUrl, oauthUtil)
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
	router.Use(middleware.SessionTracingMiddleware)
	router.Use(request_response_trace.HttpRequestResponseTracingAllMiddleware)

	corsConfig := corsModel.CorsConfig{
		AllowedOrigins: configData.AllowedOrigins,
	}

	router.Use(cors.CORSMiddleware(corsConfig))
	router.Use(globalPanicMiddleware)
	cacheControl := cacheControlMiddleware.NewCacheControlMiddleware([]string{})
	router.Use(cacheControl.StopCaching())

	oauthUtil := oauth.NewOauthUtils(configData.CryptoServiceUrl)

	logLevel := configData.LogLevel
	logger := logging.GetLogger(context.TODO())

	if logLevel != "" {
		logLevelInitErr := golaLoggerRegistry.SetLevel(logLevel)
		if logLevelInitErr != nil {
			logger.Warning("gola_logger.SetLevel failed. Default log level being used", logLevelInitErr.Error())
		}
	}

	router.GET("api/post/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	defaultRouterGroup := router.Group("api/post/v1")
	defaultRouterGroup.Use(tokenIntrospectionMiddleware(configData.OAuthUrl, oauthUtil, configData))
	draftGroup := defaultRouterGroup.Group("/draft")
	{
		draftGroup.POST("/upsert-draft", draftController.SaveDraft)
		draftGroup.POST("/tagline", draftController.SaveTagline)
		draftGroup.POST("/upsert-interests", draftController.SaveInterests)
		draftGroup.GET("/get-draft/:draft_id", draftController.GetDraft)
		draftGroup.POST("/get-all-draft", draftController.GetAllDraft)
		draftGroup.POST("/upsert-preview-image", draftController.SavePreviewImage)
		draftGroup.POST("/delete/:draft_id", draftController.DeleteDraft)
		draftGroup.POST("/delete-interest", draftController.DeleteInterest)
		draftGroup.GET("/preview-draft/:draft_id", draftController.GetPreviewDraft)
	}

	interestsGroup := defaultRouterGroup.Group("/interest")
	{
		interestsGroup.GET("/topics-and-interests", interestsController.GetExploreInterests)
		interestsGroup.POST("/get-interests", interestsController.GetInterests)
	}

	postGroup := defaultRouterGroup.Group("/post")
	{
		postGroup.POST("/publish", postController.PublishPost)
		postGroup.GET("/:post_id/like", postController.Like)
		postGroup.GET("/:post_id/unlike", postController.Unlike)
		postGroup.POST("/comment", postController.Comment)
		postGroup.GET("/:post_id", postController.GetPost)
		postGroup.GET("/:post_id/read-later", postController.MarkReadLater)
		postGroup.GET("/:post_id/remove-read-later", postController.RemoveReadLater)
	}
}
