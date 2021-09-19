package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	storyModels "post-api/story/models"
	"post-api/user-profile/constants"
	"post-api/user-profile/models"
	"post-api/user-profile/repository"
)

type UserInterestsService interface {
	GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*storyModels.JSONString, *golaerror.Error)
	FollowInterest(ctx *gin.Context, interestID, userID uuid.UUID) *golaerror.Error
	GetExploreInterests(ctx *gin.Context, userID uuid.UUID) ([]models.ExploreInterest, *golaerror.Error)
}

type userInterestsService struct {
	repository repository.UserInterestsRepository
}

func (service userInterestsService) GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*storyModels.JSONString, *golaerror.Error) {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsService").WithField("method", "GetFollowedInterest")
	log.Info("fetching followed interests")

	followedInterests, err := service.repository.GetFollowedInterest(ctx, userId)
	if err != nil {
		log.Errorf("unable to get followed interests %v", err)
		return nil, &constants.InternalServerError
	}
	return followedInterests, nil
}

func (service userInterestsService) FollowInterest(ctx *gin.Context, interestID, userID uuid.UUID) *golaerror.Error {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsService").WithField("method", "GetFollowedInterest")
	log.Info("following interests")

	err := service.repository.FollowInterest(ctx, interestID, userID)
	if err != nil {
		log.Errorf("unable to follow interest %v", err)
		return &constants.InternalServerError
	}

	return nil
}

func (service userInterestsService) GetExploreInterests(ctx *gin.Context, userID uuid.UUID) ([]models.ExploreInterest, *golaerror.Error) {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsService").WithField("method", "GetExploreInterests")
	log.Info("fetching followed interests")

	followedInterests, err := service.repository.GetExploreInterests(ctx, userID)
	if err != nil {
		log.Errorf("unable to get explore interests %v", err)
		return nil, &constants.InternalServerError
	}
	return followedInterests, nil
}

func NewUserInterestsService(repository repository.UserInterestsRepository) UserInterestsService {
	return userInterestsService{repository: repository}
}
