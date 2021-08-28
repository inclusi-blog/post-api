package middlewares

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"io/ioutil"
	"net/http"
	"post-api/idp/constants"
	"post-api/idp/service"
)

func UserActionMiddleware(service service.RegistrationCacheService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		activationHash := ctx.Param("activation_hash")
		logger := logging.GetLogger(ctx).WithField("middleware", "UserActionMiddleware")

		if activationHash == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, constants.PayloadValidationError)
			return
		}

		logger.Info("fetching user registration request from cache")
		userDetails, err := service.GetUserDetailsFromCache(activationHash, ctx)

		if err != nil {
			logger.Errorf("unable to fetch registration request from cache %v", err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, constants.ActivationLinkExpiredError)
			return
		}

		bytesData, encodeErr := json.Marshal(userDetails)

		if encodeErr != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, constants.PayloadValidationErrorCode)
			return
		}

		logger.Infof("user registration request retrieved from cache %v", userDetails)
		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bytesData))
		ctx.Next()
	}
}
