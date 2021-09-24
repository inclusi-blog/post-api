package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"net/http"
	"post-api/story/models/request"
	storyApi "post-api/story/service"
	"post-api/story/utils"
	"post-api/user-profile/constants"
	"post-api/user-profile/service"
)

type UserInterestsController struct {
	service     service.UserInterestsService
	postService storyApi.PostService
}

func (controller UserInterestsController) GetFollowedInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "GetFollowedInterests")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to upsert draft request for user %v", userUUID)

	followedInterests, interestsErr := controller.service.GetFollowedInterest(ctx, userUUID)
	if interestsErr != nil {
		logger.Errorf("unable to get followed interests %v", err)
		constants.RespondWithGolaError(ctx, interestsErr)
		return
	}

	ctx.JSON(http.StatusOK, followedInterests)
}

func (controller UserInterestsController) FollowInterest(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "GetFollowedInterests")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to upsert draft request for user %v", userUUID)
	interest := ctx.Query("interest")
	interestID, err := uuid.Parse(interest)
	if err != nil {
		logger.Errorf("unable to bind request body %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	interestsErr := controller.service.FollowInterest(ctx, interestID, userUUID)
	if interestsErr != nil {
		logger.Errorf("unable to get followed interests %v", err)
		constants.RespondWithGolaError(ctx, interestsErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (controller UserInterestsController) GetExploreInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "GetExploreInterests")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to upsert draft request for user %v", userUUID)

	followedInterests, interestsErr := controller.service.GetExploreInterests(ctx, userUUID)
	if interestsErr != nil {
		logger.Errorf("unable to get followed interests %v", err)
		constants.RespondWithGolaError(ctx, interestsErr)
		return
	}

	ctx.JSON(http.StatusOK, followedInterests)
}

func (controller UserInterestsController) UnFollowInterest(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "UnFollowInterest")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to upsert draft request for user %v", userUUID)
	interest := ctx.Query("interest")
	interestID, err := uuid.Parse(interest)
	if err != nil {
		logger.Errorf("unable to bind request body %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	interestsErr := controller.service.UnFollowInterest(ctx, interestID, userUUID)
	if interestsErr != nil {
		logger.Errorf("unable to get followed interests %v", err)
		constants.RespondWithGolaError(ctx, interestsErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (controller UserInterestsController) GetPublishedPosts(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "UnFollowInterest")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to get drafts request for user %v", userUUID)

	var postRequest request.GetPublishedPostRequest
	err = ctx.ShouldBindBodyWith(&postRequest, binding.JSON)
	if err != nil {
		logger.Errorf("Unable to bind all draft request for user %v. Error %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	postRequest.UserID = userUUID
	posts, fetchErr := controller.postService.GetPublishedPostByUser(ctx, postRequest)
	if fetchErr != nil {
		logger.Errorf("unable to get published post %v", fetchErr)
		constants.RespondWithGolaError(ctx, fetchErr)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func NewUserInterestsController(interestsService service.UserInterestsService, postService storyApi.PostService) UserInterestsController {
	return UserInterestsController{
		service:     interestsService,
		postService: postService,
	}
}
