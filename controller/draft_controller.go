package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"net/http"
	"post-api/constants"
	"post-api/models"
	"post-api/models/request"
	"post-api/service"
)

type DraftController struct {
	service service.DraftService
}

// SaveDraft godoc
// @Tags draft
// @Summary SaveDraft
// @Description Save new draft or update existing draft
// @Accept json
// @Param request body models.UpsertDraft true "Request Body"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/post/v1/draft/upsertDraft [post]
func (controller DraftController) SaveDraft(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "CreateNewPostWithData")

	log.Infof("Entered controller to upsert draft request for user")
	var upsertPost models.UpsertDraft

	err := ctx.ShouldBindBodyWith(&upsertPost, binding.JSON)

	if err != nil {
		log.Errorf("Unable to bind upsert draft request for user %v. Error %v", "12", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with upsert draft request for user %v", upsertPost.UserID)

	draftSaveErr := controller.service.SaveDraft(upsertPost, ctx)

	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving draft for user %v. Error %v", "12", draftSaveErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("writing response to draft request for user %v", upsertPost.UserID)

	ctx.Status(http.StatusOK)
}

// SaveTagline godoc
// @Tags draft
// @Summary SaveTagline
// @Description save or update tagline for draft
// @Accept json
// @Param request body request.TaglineSaveRequest true "Request Body"
// @Success 200
// @Failure 400 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/draft/tagline [post]
func (controller DraftController) SaveTagline(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "SaveTagline")

	log.Infof("Entered controller to save tagline request for user %v", "12")
	var upsertTagline request.TaglineSaveRequest

	err := ctx.ShouldBindBodyWith(&upsertTagline, binding.JSON)

	if err != nil {
		log.Errorf("Unable to bind upsert draft request for user %v. Error %v", "12", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with save tagline request for user %v", "12")

	draftSaveErr := controller.service.UpsertTagline(upsertTagline, ctx)

	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving tagline for user %v. Error %v", "12", draftSaveErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("writing response to tagline request for user %v", "12")

	ctx.Status(http.StatusOK)
}

// SaveInterests godoc
// @Tags draft
// @Summary SaveInterests
// @Description save or update interest for draft
// @Accept json
// @Param request body request.InterestsSaveRequest true "Request Body"
// @Success 200
// @Failure 400 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/draft/upsert-interests [post]
func (controller DraftController) SaveInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "SaveInterests")

	log.Infof("Entered controller to save Interests request for user %v", "12")
	var upsertInterests request.InterestsSaveRequest

	err := ctx.ShouldBindBodyWith(&upsertInterests, binding.JSON)

	if err != nil {
		log.Errorf("Unable to bind upsert interests request for user %v. Error %v", "12", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with save interests request for user %v", upsertInterests.UserID)

	draftSaveErr := controller.service.UpsertInterests(upsertInterests, ctx)

	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving interests for user %v. Error %v", "12", draftSaveErr)
		constants.RespondWithGolaError(ctx, draftSaveErr)
		return
	}

	log.Infof("writing response to interests request for user %v", upsertInterests.UserID)

	ctx.Status(http.StatusOK)
}

// DeleteInterest godoc
// @Tags draft
// @Summary DeleteInterest
// @Description delete interest for draft
// @Accept json
// @Param request body request.InterestsSaveRequest true "Request Body"
// @Success 200
// @Failure 400 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/draft/delete-interest [post]
func (controller DraftController) DeleteInterest(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "SaveInterests")

	log.Infof("Entered controller to save Interests request for user %v", "12")
	var upsertInterests request.InterestsSaveRequest

	err := ctx.ShouldBindBodyWith(&upsertInterests, binding.JSON)

	if err != nil {
		log.Errorf("Unable to bind upsert interests request for user %v. Error %v", "12", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with save interests request for user %v", upsertInterests.UserID)

	draftSaveErr := controller.service.DeleteInterest(ctx, upsertInterests)

	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving interests for user %v. Error %v", "12", draftSaveErr)
		constants.RespondWithGolaError(ctx, draftSaveErr)
		return
	}

	log.Infof("writing response to interests request for user %v", upsertInterests.UserID)

	ctx.Status(http.StatusOK)
}

// GetDraft godoc
// @Tags draft
// @Summary GetDraft
// @Description get draft for given draft id
// @Accept json
// @Param request body request.DraftURIRequest true "Request Body"
// @Success 200 {object} db.DraftDB
// @Failure 400 {object} golaerror.Error
// @Failure 404 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/draft/get-draft/:draft_id [get]
func (controller DraftController) GetDraft(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftController").WithField("method", "GetDraft")
	logger.Infof("Entered controller to get draft request for user %v", "12")

	var draftURIRequest request.DraftURIRequest
	if err := ctx.ShouldBindUri(&draftURIRequest); err != nil {
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	logger.Infof("Request body bind successful with get draft request for user %v", "12")

	draftData, draftGetErr := controller.service.GetDraft(draftURIRequest.DraftID, "some-user", ctx)
	if draftGetErr != nil {
		logger.Errorf("Error occurred in draft service while saving tagline for user %v. Error %v", "12", draftGetErr)
		constants.RespondWithGolaError(ctx, draftGetErr)
		return
	}

	logger.Infof("writing response to draft data request for user %v %s", "12", draftURIRequest.DraftID)

	ctx.JSON(http.StatusOK, draftData)
}

// SavePreviewImage godoc
// @Tags draft
// @Summary SavePreviewImage
// @Description saves preview image for draft
// @Accept json
// @Param request body request.PreviewImageSaveRequest true "Request Body"
// @Success 200
// @Failure 400 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/draft/upsert-preview-image [post]
func (controller DraftController) SavePreviewImage(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "SavePreviewImage")

	log.Infof("Entered controller to save preview image for user %v", "12")

	var imageSaveRequest request.PreviewImageSaveRequest

	if err := ctx.ShouldBindBodyWith(&imageSaveRequest, binding.JSON); err != nil {
		logger.Errorf("Unable to bind request body with image save request for draft %v", err)
		ctx.JSON(http.StatusBadRequest, &constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with image save request for user %v", "12")

	imageSaveErr := controller.service.SavePreviewImage(imageSaveRequest, ctx)
	if imageSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving image for user %v. Error %v", "12", imageSaveErr)
		constants.RespondWithGolaError(ctx, imageSaveErr)
		return
	}

	log.Infof("writing response to draft image save request for user %v %s", "12", imageSaveRequest.DraftID)

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// GetAllDraft godoc
// @Tags draft
// @Summary GetAllDraft
// @Description get all draft for given draft id
// @Accept json
// @Param request body models.GetAllDraftRequest true "Request Body"
// @Success 200 {object} []db.AllDraft
// @Failure 400 {object} golaerror.Error
// @Failure 404 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/draft/get-all-draft [post]
func (controller DraftController) GetAllDraft(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "GetAllDraft")

	log.Infof("Entered controller to get all draft request for user %v", "12")

	var allDraftReq models.GetAllDraftRequest

	err := ctx.ShouldBindBodyWith(&allDraftReq, binding.JSON)
	if err != nil {
		log.Errorf("Unable to bind all draft request for user %v. Error %v", "12", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with get all draft request for user %v", "12")

	allDraftData, draftSaveErr := controller.service.GetAllDraft(allDraftReq, ctx)
	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving tagline for user %v. Error %v", "12", draftSaveErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("writing response to draft all data request for user %v", "12")

	ctx.JSON(http.StatusOK, allDraftData)
}

// DeleteDraft godoc
// @Tags draft
// @Summary DeleteDraft
// @Description delete draft for a user
// @Accept json
// @Param request body request.DraftURIRequest true "Request Body"
// @Success 200
// @Failure 400 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/draft/delete/:draft_id [post]
func (controller DraftController) DeleteDraft(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftController").WithField("method", "DeleteDraft")

	logger.Info("Entering the controller layer to delete draft")

	var draftDeleteRequest request.DraftURIRequest

	if err := ctx.ShouldBindUri(&draftDeleteRequest); err != nil {
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	logger.Infof("Successfully bind request uri with draft delete request for draft id %v", draftDeleteRequest.DraftID)

	err := controller.service.DeleteDraft(draftDeleteRequest.DraftID, "some-user", ctx)

	if err != nil {
		logger.Errorf("error occurred while deleting draft for draft id %v", draftDeleteRequest.DraftID)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
	return
}

// GetPreviewDraft godoc
// @Tags draft
// @Summary GetPreviewDraft
// @Description get preview draft for a given draft id
// @Accept json
// @Param request body request.DraftURIRequest true "Request Body"
// @Success 200
// @Failure 400 {object} golaerror.Error
// @Failure 406 {object} golaerror.Error
// @Failure 500 {object} golaerror.Error
// @Router /api/post/v1/draft/preview-draft/:draft_id [get]
func (controller DraftController) GetPreviewDraft(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftController").WithField("method", "GetPreviewDraft")
	logger.Info("Entered preview draft controller")

	var draftURIRequest request.DraftURIRequest
	if err := ctx.ShouldBindUri(&draftURIRequest); err != nil {
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	logger.Infof("Request body bind successful with get draft request for user %v", "12")

	draftData, draftGetErr := controller.service.ValidateAndGetDraft(ctx, draftURIRequest.DraftID, "some-user")
	if draftGetErr != nil {
		logger.Errorf("Error occurred in draft service while saving tagline for user %v. Error %v", "12", draftGetErr)
		constants.RespondWithGolaError(ctx, draftGetErr)
		return
	}

	logger.Infof("writing response to draft data request for user %v %s", "12", draftURIRequest.DraftID)

	ctx.JSON(http.StatusOK, draftData)
}

func NewDraftController(service service.DraftService) DraftController {
	return DraftController{
		service: service,
	}
}
