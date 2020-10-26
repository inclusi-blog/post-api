package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"net/http"
	"post-api/constants"
	"post-api/service"
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

	var request publishRequest

	if err := ctx.ShouldBindBodyWith(&request, binding.JSON); err != nil {
		logger.Errorf("Error occurred while binding request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	id := request.DratID
	logger.Infof("Successfully bind request body for draft id %v", id)

	publishErr := controller.postService.PublishPost(ctx, id)

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

func NewPostController(postService service.PostService) PostController {
	return PostController{postService: postService}
}
