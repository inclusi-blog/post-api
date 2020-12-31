package controller

import (
	"net/http"
	"post-api/constants"
	"post-api/models/request"
	"post-api/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
)

type PostController struct {
	postService service.PostService
}

func (controller PostController) PublishPost(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "PublishPost")
	logger.Infof("Entering post controller to publish post")
	type publishRequest struct {
		DratID string `json:"draft_id" binding:"required"`
	}

	var publishPostRequest publishRequest

	if err := ctx.ShouldBindBodyWith(&publishPostRequest, binding.JSON); err != nil {
		logger.Errorf("Error occurred while binding publishPostRequest body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	id := publishPostRequest.DratID
	logger.Infof("Successfully bind publishPostRequest body for draft id %v", id)

	publishErr := controller.postService.PublishPost(ctx, id, "some-user")

	if publishErr != nil {
		logger.Errorf("Error occurred while publishing draft for draft id %v .%v", id, publishErr)
		constants.RespondWithGolaError(ctx, publishErr)
		return
	}

	logger.Infof("Successfully published draft for draft id %v", id)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "published",
	})
}

func (controller PostController) Like(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "PostController").WithField("method", "Like")

	log.Infof("Entered controller to update likes request for user %v", "12")

	var postLikeRequest request.PostLikeRequest

	if err := ctx.ShouldBindUri(&postLikeRequest); err != nil {
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with get draft request for user %v", "12")

	res, err := controller.postService.LikePost("some-user", postLikeRequest.PostUID, ctx)
	if err != nil {
		log.Errorf("Error occurred in post service while updating likedby in likes table %v. Error %v", "12", err.Error())
		constants.RespondWithGolaError(ctx, err)
		return
	}

	log.Infof("writing response to draft data request for user %v %s", "12", postLikeRequest.PostUID)

	ctx.JSON(http.StatusOK, res)
}

func (controller PostController) Unlike(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "PostController").WithField("method", "Like")

	log.Infof("Entered controller to update likes request for user %v", "12")

	var postLikeRequest request.PostLikeRequest

	if err := ctx.ShouldBindUri(&postLikeRequest); err != nil {
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with get draft request for user %v", "12")

	res, err := controller.postService.UnlikePost("some-user", postLikeRequest.PostUID, ctx)
	if err != nil {
		log.Errorf("Error occurred in post service while updating likedby in likes table %v. Error %v", "12", err.Error())
		constants.RespondWithGolaError(ctx, err)
		return
	}

	log.Infof("writing response to draft data request for user %v %s", "12", postLikeRequest.PostUID)

	ctx.JSON(http.StatusOK, res)
}

func NewPostController(postService service.PostService) PostController {
	return PostController{postService: postService}
}
