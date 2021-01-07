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

// PublishPost godoc
// @Tags post
// @Summary PublishPost
// @Description publish a draft to post
// @Accept json
// @Param request body request.PublishRequest true "Request Body"
// @Success 200
// @Failure 400 {object} golaerror.Error
// @Failure 404 {object} golaerror.Error
// @Failure 406 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/post/publish [post]
func (controller PostController) PublishPost(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "PublishPost")
	logger.Infof("Entering post controller to publish post")

	var publishPostRequest request.PublishRequest

	if err := ctx.ShouldBindBodyWith(&publishPostRequest, binding.JSON); err != nil {
		logger.Errorf("Error occurred while binding publishPostRequest body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	id := publishPostRequest.DratID
	logger.Infof("Successfully bind publishPostRequest body for draft id %v", id)

	postUrl, publishErr := controller.postService.PublishPost(ctx, id, "some-user")

	if publishErr != nil {
		logger.Errorf("Error occurred while publishing draft for draft id %v .%v", id, publishErr)
		constants.RespondWithGolaError(ctx, publishErr)
		return
	}

	logger.Infof("Successfully published draft for draft id %v", id)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "published",
		"url":    postUrl,
	})
}

// Like godoc
// @Tags post
// @Summary Like
// @Description like a post
// @Accept json
// @Param request body request.PostURIRequest true "Request Body"
// @Success 200 {object} response.LikedByCount
// @Failure 400 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/post/:post_id/like [get]
func (controller PostController) Like(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "PostController").WithField("method", "Like")

	log.Infof("Entered controller to update likes request for user %v", "12")

	var postLikeRequest request.PostURIRequest

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

// Unlike godoc
// @Tags post
// @Summary Unlike
// @Description unlike a post
// @Accept json
// @Param request body request.PostURIRequest true "Request Body"
// @Success 200 {object} response.LikedByCount
// @Failure 400 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/post/:post_id/unlike [get]
func (controller PostController) Unlike(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "PostController").WithField("method", "Like")

	log.Infof("Entered controller to update likes request for user %v", "12")

	var postLikeRequest request.PostURIRequest

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

// Comment godoc
// @Tags post
// @Summary Comment
// @Description comment on a post
// @Accept json
// @Param request body request.CommentPost true "Request Body"
// @Success 200
// @Failure 400 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/post/comment [post]
func (controller PostController) Comment(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "Comment")
	logger.Infof("Entering controller to comment on post")

	var commentRequest request.CommentPost

	err := ctx.ShouldBindBodyWith(&commentRequest, binding.JSON)

	if err != nil {
		logger.Errorf("Error occurred while binding commentRequest body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	logger.Infof("Successfully bind body with post comment request by user %v", "some-user")

	commentErr := controller.postService.CommentPost(ctx, "some-user", commentRequest.PostUID, commentRequest.Comment)

	if commentErr != nil {
		logger.Errorf("Error occurred while commenting on post by user %v, Error %v", "some-user", commentErr)
		constants.RespondWithGolaError(ctx, commentErr)
		return
	}

	logger.Infof("Successfully commented on post %v", commentRequest.PostUID)
	ctx.Status(http.StatusOK)
}

// GetPost godoc
// @Tags post
// @Summary GetPost
// @Description get a post
// @Accept json
// @Param request body request.PostURIRequest true "Request Body"
// @Success 200 {object} response.Post
// @Failure 400 {object} golaerror.Error
// @Failure 404 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/post/:post_id [get]
func (controller PostController) GetPost(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "GetPost")
	logger.Info("Started get post to fetch post for the given post id")

	var postRequest request.PostURIRequest
	if err := ctx.ShouldBindUri(&postRequest); err != nil {
		logger.Errorf("Error occurred while binding get post request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	id := postRequest.PostUID
	logger.Infof("Successfully bind get post request body for post id %v", id)
	post, err := controller.postService.GetPost(ctx, id, "some-user")

	if err != nil {
		logger.Errorf("Error occurred while publishing draft for draft id %v .%v", id, err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Infof("Successfully fetching post for given post id %v", id)
	ctx.JSON(http.StatusOK, post)
}

func NewPostController(postService service.PostService) PostController {
	return PostController{postService: postService}
}
