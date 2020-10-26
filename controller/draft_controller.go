package controller

import (
	"net/http"
	"post-api/constants"
	"post-api/models"
	"post-api/models/request"
	"post-api/service"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
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
func (draftController DraftController) SaveDraft(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "SavePostDraft")

	log.Infof("Entered controller to upsert draft request for user %v", "12")
	var upsertPost models.UpsertDraft

	err := ctx.ShouldBindBodyWith(&upsertPost, binding.JSON)

	if err != nil {
		log.Errorf("Unable to bind upsert draft request for user %v. Error %v", "12", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Request body bind successful with upsert draft request for user %v", "12")

	draftSaveErr := draftController.service.SaveDraft(upsertPost, ctx)

	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving draft for user %v. Error %v", "12", draftSaveErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("writing response to draft request for user %v", "12")

	ctx.Status(http.StatusOK)
}

func (draftController DraftController) SaveTagline(ctx *gin.Context) {
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

	draftSaveErr := draftController.service.UpsertTagline(upsertTagline, ctx)

	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving tagline for user %v. Error %v", "12", draftSaveErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("writing response to tagline request for user %v", "12")

	ctx.Status(http.StatusOK)
}

func (draftController DraftController) SaveInterests(ctx *gin.Context) {
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

	log.Infof("Request body bind successful with save interests request for user %v", "12")

	draftSaveErr := draftController.service.UpsertInterests(upsertInterests, ctx)

	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving interests for user %v. Error %v", "12", draftSaveErr)
		constants.RespondWithGolaError(ctx, draftSaveErr)
		return
	}

	log.Infof("writing response to interests request for user %v", "12")

	ctx.Status(http.StatusOK)
}

func (draftController DraftController) GetDraft(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "GetDraft")

	log.Infof("Entered controller to get draft request for user %v", "12")

	queryParams := ctx.Request.URL.Query()
	draftUID := queryParams.Get("draft_id")
	if draftUID == "" {
		log.Errorf("Draft ID is not recieved %v", "12")
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}
	log.Infof("Request body bind successful with get draft request for user %v", "12")

	draftData, draftSaveErr := draftController.service.GetDraft(draftUID, ctx)
	if draftSaveErr != nil {
		log.Errorf("Error occurred in draft service while saving tagline for user %v. Error %v", "12", draftSaveErr)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("writing response to draft data request for user %v %s", "12", draftUID)

	ctx.JSON(http.StatusOK, draftData)

}

func NewDraftController(service service.DraftService) DraftController {
	return DraftController{
		service: service,
	}
}
