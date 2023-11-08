package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/logging"
	"net/http"
	"post-api/story/constants"
	"post-api/story/models/request"
	"post-api/story/service"
	"post-api/story/utils"
)

type ReportController struct {
	service service.ReportService
}

func (controller ReportController) ReportPost(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "ReportController").WithField("method", "PublishPost")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entering post controller to report post")

	postUID := ctx.Param("post_id")
	postID, err := uuid.Parse(postUID)
	if err != nil {
		logger.Errorf("invalid post id request for user %v. Error %v", userUUID, err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	report := new(request.Report)
	if err := ctx.ShouldBindJSON(report); err != nil {
		logger.Errorf("invalid report or report not found %v. Error %v", userUUID, err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	report.UserID = userUUID
	report.PostID = postID

	logger.Infof("Successfully bind Report post body for post id %v", postID)
	reportErr := controller.service.ReportPost(ctx, *report)

	if reportErr != nil {
		logger.Errorf("Error occurred while reporting post id %v .%v", postID, reportErr)
		constants.RespondWithGolaError(ctx, reportErr)
		return
	}

	logger.Infof("Successfully reported post for post id %v", postID)
	ctx.Status(200)
}

func NewReportController(reportService service.ReportService) ReportController {
	return ReportController{
		service: reportService,
	}
}
