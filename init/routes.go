package init

import (
	"context"
	"net/http"
	"post-api/configuration"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/middleware/request_response_trace"
	middleware "github.com/gola-glitch/gola-utils/middleware/session_trace"
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
		draftGroup.POST("/upsertDraft", draftController.SaveDraft)
		draftGroup.POST("/tagline", draftController.SaveTagline)
		draftGroup.POST("/upsertInterests", draftController.SaveInterests)
	}

	defaultRouterGroup.POST("/get-interests", interestsController.GetInterests)
}
