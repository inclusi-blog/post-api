package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"net/http"
	"post-api/models"
	"post-api/service"
)

type DraftController struct {
	service service.DraftService
}

func (draftController DraftController) SaveDraft(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)

	log := logger.WithField("class", "DraftController").WithField("method", "SaveDraft")

	log.Infof("Entered controller to upsert draft request for user %v", "12")
	var upsertPost models.UpsertDraft

	err := ctx.ShouldBindBodyWith(&upsertPost, binding.JSON)

	if err != nil {
		log.Errorf("Unable to bind upsert draft request for user %v. Error %v", "12", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	log.Infof("Request body bind successful with upsert draft request for user %v", "12")

	err = draftController.service.SaveDraft(upsertPost, ctx)

	if err != nil {
		log.Errorf("Error occurred in draft service while saving draft for user %v. Error %v", "12", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("writing response to draft request for user %v", "12")

	ctx.Status(http.StatusOK)
}

func NewDraftController(service service.DraftService) DraftController {
	return DraftController{
		service: service,
	}
}
