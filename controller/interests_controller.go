package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"net/http"
	"post-api/constants"
	"post-api/models/request"
	"post-api/service"
)

type InterestsController struct {
	service service.InterestsService
}

// GetInterests godoc
// @Tags interest
// @Summary GetInterests
// @Description get all interests or search a interest with keyword
// @Accept json
// @Param request body request.SearchInterests true "Request Body"
// @Success 200 {object} []db.Interest
// @Failure 400 {object} golaerror.Error
// @Failure 404 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/post/interest/get-interests [post]
func (controller InterestsController) GetInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "interestsController").WithField("method", "GetInterests")

	logger.Info("Entered interests controller to get interests")

	logger.Info("binding request body for search interests request")

	var searchInterestRequest request.SearchInterests

	bindingErr := ctx.ShouldBindBodyWith(&searchInterestRequest, binding.JSON)

	if bindingErr != nil {
		logger.Errorf("Unable to bind request body for search key interests %v", bindingErr)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	interests, err := controller.service.GetInterests(ctx, searchInterestRequest.SearchKeyword, searchInterestRequest.SelectedTags)

	if err != nil {
		logger.Errorf("Error occurred while fetching over all interests from interests service %v", err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Info("Successfully got interests")

	ctx.JSON(http.StatusOK, interests)
}

// GetExploreInterests godoc
// @Tags interest
// @Summary GetExploreInterests
// @Description get all categories and interests
// @Accept json
// @Success 200 {object} []response.CategoryAndInterest
// @Failure 400 {object} golaerror.Error
// @Failure 404 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/interest/topics-and-interests [get]
func (controller InterestsController) GetExploreInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsController").WithField("method", "GetExploreInterests")

	logger.Info("Entered Interests controller to get categories and interests")

	categoriesAndInterests, err := controller.service.GetExploreCategoriesAndInterests(ctx)

	if err != nil {
		logger.Errorf("Error occurred while fetching over all categories and interests from interests service %v", err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Info("Successfully got categories and interests")

	ctx.JSON(http.StatusOK, categoriesAndInterests)
}

func NewInterestsController(interestsService service.InterestsService) InterestsController {
	return InterestsController{service: interestsService}
}
