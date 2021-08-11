package controller

import (
	"github.com/google/uuid"
	"net/http"
	"post-api/constants"
	"post-api/service"
	"post-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
)

type PostController struct {
	postService service.PostService
}

func (controller PostController) PublishPost(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "PublishPost")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entering post controller to publish post")

	draftUUID := ctx.Query("draft")
	draftID, err := uuid.Parse(draftUUID)
	if err != nil {
		logger.Errorf("invalid draft id request for user %v. Error %v", userUUID, err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	logger.Infof("Successfully bind publishPostRequest body for draft id %v", draftID)
	publishErr := controller.postService.PublishPost(ctx, draftID, userUUID)

	if publishErr != nil {
		logger.Errorf("Error occurred while publishing draft for draft id %v .%v", draftID, publishErr)
		constants.RespondWithGolaError(ctx, publishErr)
		return
	}

	logger.Infof("Successfully published draft for draft id %v", draftID)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "published",
	})
}

func (controller PostController) Like(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "UpdateLikes")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to update likes request for user %v", userUUID)

	postID := ctx.Query("post")
	draftUUID, err := uuid.Parse(postID)
	if err != nil {
		logger.Errorf("invalid draft id request for user %v. Error %v", userUUID, err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	logger.Infof("Request body bind successful with get draft request for user %v", userUUID)

	serviceErr := controller.postService.LikePost(ctx, draftUUID, userUUID)
	if serviceErr != nil {
		logger.Errorf("Error occurred in post service while updating like in likes table %v. Error %v", userUUID.String(), serviceErr.Error())
		constants.RespondWithGolaError(ctx, serviceErr)
		return
	}
	logger.Infof("writing response to draft data request for user %v %s", userUUID, postID)

	ctx.Status(http.StatusOK)
}

func NewPostController(postService service.PostService) PostController {
	return PostController{postService: postService}
}
