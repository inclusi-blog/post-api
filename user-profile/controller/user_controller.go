package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"net/http"
	commonService "post-api/service"
	"post-api/story/models/request"
	storyApi "post-api/story/service"
	"post-api/story/utils"
	"post-api/user-profile/constants"
	"post-api/user-profile/service"
	"time"
)

type UserProfileController struct {
	service            service.UserInterestsService
	userProfileService service.ProfileService
	postService        storyApi.PostService
	awsServices        commonService.AwsServices
}

func (controller UserProfileController) GetFollowedInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "GetFollowedInterests")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to upsert draft request for user %v", userUUID)

	followedInterests, interestsErr := controller.service.GetFollowedInterest(ctx, userUUID)
	if interestsErr != nil {
		logger.Errorf("unable to get followed interests %v", err)
		constants.RespondWithGolaError(ctx, interestsErr)
		return
	}

	ctx.JSON(http.StatusOK, followedInterests)
}

func (controller UserProfileController) FollowInterest(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "GetFollowedInterests")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to upsert draft request for user %v", userUUID)
	interest := ctx.Query("interest")
	interestID, err := uuid.Parse(interest)
	if err != nil {
		logger.Errorf("unable to bind request body %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	interestsErr := controller.service.FollowInterest(ctx, interestID, userUUID)
	if interestsErr != nil {
		logger.Errorf("unable to get followed interests %v", err)
		constants.RespondWithGolaError(ctx, interestsErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (controller UserProfileController) GetExploreInterests(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "GetExploreInterests")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to upsert draft request for user %v", userUUID)

	followedInterests, interestsErr := controller.service.GetExploreInterests(ctx, userUUID)
	if interestsErr != nil {
		logger.Errorf("unable to get followed interests %v", err)
		constants.RespondWithGolaError(ctx, interestsErr)
		return
	}

	ctx.JSON(http.StatusOK, followedInterests)
}

func (controller UserProfileController) UnFollowInterest(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "UnFollowInterest")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to upsert draft request for user %v", userUUID)
	interest := ctx.Query("interest")
	interestID, err := uuid.Parse(interest)
	if err != nil {
		logger.Errorf("unable to bind request body %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	interestsErr := controller.service.UnFollowInterest(ctx, interestID, userUUID)
	if interestsErr != nil {
		logger.Errorf("unable to get followed interests %v", err)
		constants.RespondWithGolaError(ctx, interestsErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (controller UserProfileController) GetPublishedPosts(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "UnFollowInterest")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to get drafts request for user %v", userUUID)

	var postRequest request.GetPublishedPostRequest
	err = ctx.ShouldBindBodyWith(&postRequest, binding.JSON)
	if err != nil {
		logger.Errorf("Unable to bind all draft request for user %v. Error %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	postRequest.UserID = userUUID
	posts, fetchErr := controller.postService.GetPublishedPostByUser(ctx, postRequest)
	if fetchErr != nil {
		logger.Errorf("unable to get published post %v", fetchErr)
		constants.RespondWithGolaError(ctx, fetchErr)
		return
	}

	ctx.JSON(http.StatusOK, posts)
}

func (controller UserProfileController) GetDetails(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "UnFollowInterest")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	userUUID, _ := uuid.Parse(token.UserId)
	logger.Infof("Entered controller to get drafts request for user %v", userUUID)

	profile, profileErr := controller.userProfileService.GetProfile(ctx, userUUID)
	if profileErr != nil {
		logger.Errorf("unable to get profile %v", err)
		constants.RespondWithGolaError(ctx, profileErr)
		return
	}

	if profile.ProfilePic != nil {
		profilePic, err := controller.awsServices.GetObjectInS3(*profile.ProfilePic, time.Hour*time.Duration(6))
		if err != nil {
			logger.Errorf("unable to get object from s3 %v", err)
		}
		profile.ProfilePic = &profilePic
	}

	ctx.JSON(http.StatusOK, profile)
}

func (controller UserProfileController) ViewProfileAvatar(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserController").WithField("method", "UnFollowInterest")
	userID := ctx.Param("user_id")
	userUID, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("unable to bind request path param %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}
	logger.Infof("Entered controller to get drafts request for user %v", userUID)

	avatar, viewErr := controller.userProfileService.FetchProfileAvatar(ctx, userUID)
	if viewErr != nil {
		logger.Errorf("unable to fetch view avatar %v", viewErr)
		constants.RespondWithGolaError(ctx, viewErr)
		return
	}

	ctx.Redirect(http.StatusMovedPermanently, avatar)
}

func (controller UserProfileController) FollowUser(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserProfileController").WithField("method", "FollowUser")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	userUID, _ := uuid.Parse(token.UserId)
	userID := ctx.Param("user_id")
	followingUID, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("unable to bind request path param %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	followErr := controller.userProfileService.FollowUser(ctx, userUID, followingUID)
	if followErr != nil {
		logger.Errorf("unable to follow user %v", followErr)
		constants.RespondWithGolaError(ctx, followErr)
		return
	}
	ctx.Status(200)
}

func (controller UserProfileController) UnFollowUser(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "UserProfileController").WithField("method", "FollowUser")
	token, err := utils.GetIDToken(ctx)
	if err != nil {
		logger.Error("id token not found", err)
		ctx.JSON(http.StatusInternalServerError, constants.InternalServerError)
		return
	}
	userUID, _ := uuid.Parse(token.UserId)
	userID := ctx.Param("user_id")
	followingUID, err := uuid.Parse(userID)
	if err != nil {
		logger.Errorf("unable to bind request path param %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	followErr := controller.userProfileService.UnFollowUser(ctx, userUID, followingUID)
	if followErr != nil {
		logger.Errorf("unable to follow user %v", followErr)
		constants.RespondWithGolaError(ctx, followErr)
		return
	}
	ctx.Status(200)
}

func NewUserProfileController(interestsService service.UserInterestsService, postService storyApi.PostService, profileService service.ProfileService, services commonService.AwsServices) UserProfileController {
	return UserProfileController{
		service:            interestsService,
		userProfileService: profileService,
		postService:        postService,
		awsServices:        services,
	}
}
