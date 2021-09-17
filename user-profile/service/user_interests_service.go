package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"post-api/story/models"
	"post-api/user-profile/constants"
	"post-api/user-profile/repository"
)

type UserInterestsService interface {
	GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*models.JSONString, error)
}

type userInterestsService struct {
	repository repository.UserInterestsRepository
}

func (service userInterestsService) GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*models.JSONString, error) {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsService").WithField("method", "GetFollowedInterest")
	log.Info("fetching followed interests")

	followedInterests, err := service.repository.GetFollowedInterest(ctx, userId)
	if err != nil {
		log.Errorf("unable to get followed interests %v", err)
		return nil, constants.InternalServerError
	}
	return followedInterests, nil
}

func NewUserInterestsService(repository repository.UserInterestsRepository) UserInterestsService {
	return userInterestsService{repository: repository}
}
