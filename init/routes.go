package init

import (
	"context"
	"net/http"
	"post-api/configuration"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	cors "github.com/gola-glitch/gola-utils/middleware/cors"
	"github.com/gola-glitch/gola-utils/middleware/request_response_trace"
	middleware "github.com/gola-glitch/gola-utils/middleware/session_trace"
	corsModel "github.com/gola-glitch/gola-utils/model"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRouter(router *gin.Engine, configData *configuration.ConfigData) {
	router.GET("api/post/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"Status": "Up",
		})
	})

	router.Use(middleware.SessionTracingMiddleware)
	router.Use(request_response_trace.HttpRequestResponseTracingMiddleware([]request_response_trace.IgnoreRequestResponseLogs{
		{
			PartialApiPath:       "api/post/healthz",
			IsRequestLogAllowed:  false,
			IsResponseLogAllowed: false,
		},
	}))

	corsConfig := corsModel.CorsConfig{
		AllowedOrigins: configData.AllowedOrigins,
	}

	router.Use(cors.CORSMiddleware(corsConfig))
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

	router.GET("api/post/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	defaultRouterGroup := router.Group("api/post/v1")

	draftGroup := defaultRouterGroup.Group("/draft")
	{
		draftGroup.POST("/upsert-draft", draftController.SaveDraft)
		draftGroup.POST("/tagline", draftController.SaveTagline)
		draftGroup.POST("/upsert-interests", draftController.SaveInterests)
		draftGroup.GET("/get-draft/:draft_id", draftController.GetDraft)
		draftGroup.POST("/get-all-draft", draftController.GetAllDraft)
		draftGroup.POST("/upsert-preview-image", draftController.SavePreviewImage)
		draftGroup.POST("/delete", draftController.DeleteDraft)
		draftGroup.POST("/delete-interest", draftController.DeleteInterest)
	}

	defaultRouterGroup.POST("/get-interests", interestsController.GetInterests)

	postGroup := defaultRouterGroup.Group("/post")
	{
		postGroup.POST("/publish", postController.PublishPost)
		postGroup.GET("/:post_id/like", postController.Like)
		postGroup.GET("/:post_id/unlike", postController.Unlike)
		postGroup.POST("/comment", postController.Comment)
	}
}
