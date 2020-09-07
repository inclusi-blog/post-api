package init

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/middleware/request_response_trace"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"post-api/configuration"
)

func RegisterRouter(router *gin.Engine, configData *configuration.ConfigData) {
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

	draftGroup := router.Group("api/post/v1/draft")
	{
		draftGroup.POST("/upsertDraft", draftController.SaveDraft)
	}
}
