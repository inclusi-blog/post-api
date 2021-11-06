package controller

import (
	"github.com/google/uuid"
	"net/http"
	"post-api/story/constants"
	"post-api/story/models/request"
	"post-api/story/service"
	"post-api/story/utils"

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
	postURL, publishErr := controller.postService.PublishPost(ctx, draftID, userUUID)

	if publishErr != nil {
		logger.Errorf("Error occurred while publishing draft for draft id %v .%v", draftID, publishErr)
		constants.RespondWithGolaError(ctx, publishErr)
		return
	}

	logger.Infof("Successfully published draft for draft id %v", draftID)
	ctx.JSON(http.StatusOK, gin.H{
		"status": "published",
		"url":    postURL,
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

func (controller PostController) UnLike(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "UnLike")
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

	serviceErr := controller.postService.UnLikePost(ctx, draftUUID, userUUID)
	if serviceErr != nil {
		logger.Errorf("Error occurred in post service while updating like in likes table %v. Error %v", userUUID.String(), serviceErr.Error())
		constants.RespondWithGolaError(ctx, serviceErr)
		return
	}
	logger.Infof("writing response to draft data request for user %v %s", userUUID, postID)

	ctx.Status(http.StatusOK)
}

func (controller PostController) Comment(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "Comment")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to update likes request for user %v", userUUID)

	var postRequest request.PostURIRequest
	if err := ctx.ShouldBindUri(&postRequest); err != nil {
		logger.Errorf("Error occurred while binding get post request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}
	id, _ := uuid.Parse(postRequest.PostUID)

	var comment request.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		logger.Errorf("unable to bind request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	logger.Infof("Request body bind successful with get draft request for user %v", userUUID)
	comment.CommentedBy = userUUID
	comment.PostID = id

	serviceErr := controller.postService.Comment(ctx, comment)
	if serviceErr != nil {
		logger.Errorf("Error occurred in post service while updating like in likes table %v. Error %v", userUUID.String(), serviceErr.Error())
		constants.RespondWithGolaError(ctx, serviceErr)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller PostController) GetComments(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "GetComments")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to update likes request for user %v", userUUID)

	var postRequest request.PostURIRequest
	if err := ctx.ShouldBindUri(&postRequest); err != nil {
		logger.Errorf("Error occurred while binding get post request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}
	id, _ := uuid.Parse(postRequest.PostUID)

	var commentRequest request.FetchComments
	if err := ctx.ShouldBindQuery(&commentRequest); err != nil {
		logger.Errorf("unable to bind request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	commentRequest.PostID = id
	logger.Infof("Request body bind successful with get draft request for user %v", userUUID)

	comments, serviceErr := controller.postService.GetComments(ctx, commentRequest)
	if serviceErr != nil {
		logger.Errorf("Error occurred in post service while updating like in likes table %v. Error %v", userUUID.String(), serviceErr.Error())
		constants.RespondWithGolaError(ctx, serviceErr)
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func (controller PostController) SavePost(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "GetComments")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to update likes request for user %v", userUUID)

	var postRequest request.PostURIRequest
	if err := ctx.ShouldBindUri(&postRequest); err != nil {
		logger.Errorf("Error occurred while binding get post request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}
	id, _ := uuid.Parse(postRequest.PostUID)

	markErr := controller.postService.SavePost(ctx, id, userUUID)
	if markErr != nil {
		logger.Errorf("unable to mark post as read later %v", markErr)
		constants.RespondWithGolaError(ctx, markErr)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller PostController) RemoveBookmark(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "RemoveBookmark")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to update likes request for user %v", userUUID)

	var postRequest request.PostURIRequest
	if err := ctx.ShouldBindUri(&postRequest); err != nil {
		logger.Errorf("Error occurred while binding get post request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}
	id, _ := uuid.Parse(postRequest.PostUID)

	bookmarkErr := controller.postService.RemovePostBookmark(ctx, id, userUUID)
	if bookmarkErr != nil {
		logger.Errorf("unable to mark post as read later %v", bookmarkErr)
		constants.RespondWithGolaError(ctx, bookmarkErr)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller PostController) MarkAsViewed(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "MarkAsViewed")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to update likes request for user %v", userUUID)

	var postRequest request.PostURIRequest
	if err := ctx.ShouldBindUri(&postRequest); err != nil {
		logger.Errorf("Error occurred while binding get post request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}
	id, _ := uuid.Parse(postRequest.PostUID)

	markErr := controller.postService.MarkAsViewed(ctx, id, userUUID)
	if markErr != nil {
		logger.Errorf("unable to mark post as read later %v", err)
		constants.RespondWithGolaError(ctx, markErr)
		return
	}

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

	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entering post controller to publish post")

	var postRequest request.PostURIRequest
	if err := ctx.ShouldBindUri(&postRequest); err != nil {
		logger.Errorf("Error occurred while binding get post request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	id, _ := uuid.Parse(postRequest.PostUID)
	logger.Infof("Successfully bind get post request body for post id %v", id)
	post, publishErr := controller.postService.GetPost(ctx, id, userUUID)

	if publishErr != nil {
		logger.Errorf("Error occurred while publishing draft for draft id %v .%v", id, publishErr)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Infof("Successfully fetching post for given post id %v", id)
	ctx.JSON(http.StatusOK, post)
}

func (controller PostController) GetReadLaterPosts(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "GetPost")
	logger.Info("Started get post to fetch post for the given post id")

	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entering post controller to publish post")

	var postRequest request.PostRequest
	if err := ctx.ShouldBindQuery(&postRequest); err != nil {
		logger.Errorf("unable to bind request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	postRequest.UserID = userUUID
	logger.Infof("Request body bind successful with get draft request for user %v", userUUID)

	posts, fetchErr := controller.postService.FetchSavedPosts(ctx, postRequest)
	if fetchErr != nil {
		logger.Errorf("unable to get read later posts %v", fetchErr)
		constants.RespondWithGolaError(ctx, fetchErr)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (controller PostController) GetReadPosts(ctx *gin.Context) {
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

	var postRequest request.PostRequest
	if err := ctx.ShouldBindQuery(&postRequest); err != nil {
		logger.Errorf("unable to bind request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	postRequest.UserID = userUUID
	logger.Infof("Request body bind successful with get draft request for user %v", userUUID)

	posts, fetchErr := controller.postService.FetchViewedPosts(ctx, postRequest)
	if fetchErr != nil {
		logger.Errorf("unable to get read later posts %v", fetchErr)
		constants.RespondWithGolaError(ctx, fetchErr)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (controller PostController) GetPostsByInterests(ctx *gin.Context) {
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

	var interestRequest request.InterestRequest
	if err := ctx.ShouldBindQuery(&interestRequest); err != nil {
		logger.Errorf("unable to bind request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	logger.Infof("Request body bind successful with get draft request for user %v", userUUID)

	var interestURIRequest request.InterestURIRequest
	if err := ctx.ShouldBindUri(&interestURIRequest); err != nil {
		logger.Errorf("unable to bind request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	logger.Infof("Request body bind successful with get draft request for user %v", userUUID)
	interestRequest.InterestUID, _  = uuid.Parse(interestURIRequest.InterestUID)

	posts, fetchErr := controller.postService.FetchPostsByInterests(ctx, interestRequest, userUUID)
	if fetchErr != nil {
		logger.Errorf("unable to get posts by interests %v", fetchErr)
		constants.RespondWithGolaError(ctx, fetchErr)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (controller PostController) Delete(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "PostController").WithField("method", "Delete")
	logger.Info("Started get post to fetch post for the given post id")

	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entering post controller to publish post")

	var postRequest request.PostURIRequest
	if err := ctx.ShouldBindUri(&postRequest); err != nil {
		logger.Errorf("Error occurred while binding get post request body %v", err)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	id, _ := uuid.Parse(postRequest.PostUID)
	logger.Infof("Successfully bind get post request body for post id %v", id)
	deleteErr := controller.postService.Delete(ctx, id, userUUID)

	if deleteErr != nil {
		logger.Errorf("Error occurred while publishing draft for draft id %v .%v", id, deleteErr)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Infof("Successfully fetching post for given post id %v", id)
	ctx.Status(http.StatusOK)
}

func NewPostController(postService service.PostService) PostController {
	return PostController{postService: postService}
}
