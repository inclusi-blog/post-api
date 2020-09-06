package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"post-api/models"
	"post-api/service"
)

type DraftController struct {
	service service.DraftService
}

func (draftController DraftController) SaveDraft(ctx *gin.Context) {
	var upsertPost models.UpsertDraft

	err := ctx.ShouldBindBodyWith(&upsertPost, binding.JSON)

	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = draftController.service.SaveDraft(upsertPost, ctx)

	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.Status(http.StatusOK)
}

func NewDraftController(service service.DraftService) DraftController {
	return DraftController{
		service: service,
	}
}
