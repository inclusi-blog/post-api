package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/logging"
	"net/http"
	"post-api/idp/constants"
	"post-api/idp/models/request"
	"post-api/idp/service"
	commonModels "post-api/models"
	commonService "post-api/service"
	"post-api/story/utils"
	"strings"
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

	var p commonModels.CoverPreSign
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

func (controller UserDetailsController) UploadImageKey(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsController").WithField("method", "UpdateUserDetails")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)

	logger.Infof("Entered controller to update avatar key for user %v", token.UserId)
	var upload commonModels.UploadImage
	if err := ctx.ShouldBindJSON(&upload); err != nil {
		logger.Errorf("unable to bind request body %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	hasPrefix := strings.HasPrefix(upload.UploadID, "profile/")
	if !hasPrefix {
		logger.Error("invalid image to upload for profile")
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	ps, err := controller.awsService.CheckS3Object(upload.UploadID)
	if err != nil {
		logger.Errorf("unable to check object existence for upload id %v, Error: %v", upload.UploadID, err)
		constants.RespondWithGolaError(ctx, &constants.UnableToFetchObjectError)
		return
	}

	if !ps {
		logger.Error("object not found")
		constants.RespondWithGolaError(ctx, &constants.ObjectNotFoundError)
		return
	}

	uploadErr := controller.service.UpdateProfileImage(ctx, upload.UploadID, userUUID)
	if uploadErr != nil {
		logger.Errorf("unable to upload image %v", uploadErr)
		constants.RespondWithGolaError(ctx, uploadErr)
		return
	}

	ctx.Status(http.StatusOK)
}

func NewUserDetailsController(service service.UserDetailsService, services commonService.AwsServices) UserDetailsController {
	return UserDetailsController{
		service:    service,
		awsService: services,
	}
}
