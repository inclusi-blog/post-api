package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"net/http"
	"post-api/idp/constants"
	"post-api/idp/models/request"
	"post-api/idp/service"
	commonService "post-api/service"
	"post-api/story/utils"
)

type UserDetailsController struct {
	service    service.UserDetailsService
	awsService commonService.AwsServices
}

func (controller UserDetailsController) UpdateUserDetails(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsController").WithField("method", "UpdateUserDetails")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to get draft request for user %v", userUUID)

	var detailsUpdate request.UserDetailsUpdate

	if err = ctx.ShouldBindBodyWith(&detailsUpdate, binding.JSON); err != nil {
		logger.Errorf("Unable to bind request body for new details update request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	logger.Infof("Successfully bind request body %v", detailsUpdate)
	updateErr := controller.service.UpdateUserDetails(ctx, userUUID, detailsUpdate)

	if updateErr != nil {
		logger.Errorf("Unable to update user details. Error %v", updateErr)
		constants.RespondWithGolaError(ctx, updateErr)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller UserDetailsController) GetPreSignURLForProfilePic(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsController").WithField("method", "UpdateUserDetails")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	logger.Infof("Entered controller to get draft request for user %v", token.UserId)

	var p request.CoverPreSign
	p.Extension = "jpg"
	if err := ctx.ShouldBindQuery(&p); err != nil {
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	key := fmt.Sprintf("profile/%s.%s", token.UserId, p.Extension)

	url, s3Err := controller.awsService.PutObjectInS3(key)
	if s3Err != nil {
		logger.Errorf("unable to put object in s3 %v", s3Err)
		ctx.JSON(http.StatusBadRequest, constants.UnableToAssignPreSignURLError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

func NewUserDetailsController(service service.UserDetailsService, services commonService.AwsServices) UserDetailsController {
	return UserDetailsController{
		service:    service,
		awsService: services,
	}
}
