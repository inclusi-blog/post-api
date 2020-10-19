package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"net/http"
	"post-api/constants"
	"post-api/service"
)

type InterestsController struct {
	service service.InterestsService
}

func (controller InterestsController) GetInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "interestsController").WithField("method", "GetInterests")

	logger.Info("Entered interests controller to get interests")

	interests, err := controller.service.GetInterests(ctx)

	if err != nil {
		logger.Errorf("Error occurred while fetching over all interests from interests service %v", err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Info("Successfully got interests")

	ctx.JSON(http.StatusOK, interests)
}

func NewInterestsController(interestsService service.InterestsService) InterestsController {
	return InterestsController{service: interestsService}
}
