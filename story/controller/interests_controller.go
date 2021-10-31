package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"net/http"
	"post-api/story/constants"
	"post-api/story/models/request"
	"post-api/story/service"
	"post-api/story/utils"
)

type InterestsController struct {
	service service.InterestsService
}

func (controller InterestsController) GetInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "interestsController").WithField("method", "GetInterests")
	logger.Info("Entered interests controller to get interests")
	logger.Info("binding request body for search interests request")

	interests, err := controller.service.GetInterests(ctx)
	if err != nil {
		logger.Errorf("Error occurred while fetching over all interests from interests service %v", err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Info("Successfully got interests")
	ctx.JSON(http.StatusOK, interests)
}

func (controller InterestsController) GetInterestDetails(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "GetReadPosts")
	logger.Info("Started get post to fetch post for the given post id")

	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entering post controller to publish post")

	var interestURIRequest request.InterestNameURIRequest
	if err := ctx.ShouldBindUri(&interestURIRequest); err != nil {
		logger.Errorf("unable to bind request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	details, fetchErr := controller.service.GetFollowCount(ctx, interestURIRequest.Name, userUUID)
	if fetchErr != nil {
		logger.Errorf("unable to fetch interest details %v", err)
		constants.RespondWithGolaError(ctx, fetchErr)
		return
	}

	ctx.JSON(http.StatusOK, details)
}

func NewInterestsController(interestsService service.InterestsService) InterestsController {
	return InterestsController{service: interestsService}
}
