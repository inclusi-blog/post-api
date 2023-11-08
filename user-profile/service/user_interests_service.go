package service

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
	"post-api/service"
	storyModels "post-api/story/models"
	"post-api/user-profile/constants"
	"post-api/user-profile/models"
	"post-api/user-profile/repository"
	"time"
)

type UserInterestsService interface {
	GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*storyModels.JSONString, *golaerror.Error)
	FollowInterest(ctx *gin.Context, interestID, userID uuid.UUID) *golaerror.Error
	GetExploreInterests(ctx *gin.Context, userID uuid.UUID) ([]models.ExploreInterests, *golaerror.Error)
	UnFollowInterest(ctx *gin.Context, interestID, userID uuid.UUID) *golaerror.Error
}

type userInterestsService struct {
	repository repository.UserInterestsRepository
	awsService service.AwsServices
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

func (service userInterestsService) GetExploreInterests(ctx *gin.Context, userID uuid.UUID) ([]models.ExploreInterests, *golaerror.Error) {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsService").WithField("method", "GetExploreInterests")
	log.Info("fetching followed interests")

	followedInterests, err := service.repository.GetExploreInterests(ctx, userID)
	if err != nil {
		log.Errorf("unable to get explore interests %v", err)
		return nil, &constants.InternalServerError
	}
	var exploreInterest []models.ExploreInterests
	for _, category := range followedInterests {
		var interests []models.Interest
		err := category.Interests.Unmarshal(&interests)
		if err != nil {
			log.Errorf("unable to marshal interest model %v", err)
			return nil, &constants.InternalServerError
		}
		for i, _ := range interests {
			interests[i].CoverPic, err = service.awsService.GetObjectInS3(interests[i].CoverPic, time.Hour*time.Duration(6))
			if err != nil {
				log.Errorf("unable to fetch interest image cover %v", err)
				return nil, &constants.InternalServerError
			}
		}
		exploreInterest = append(exploreInterest, models.ExploreInterests{
			Category:  category.Category,
			Interests: interests,
		})
	}

	return exploreInterest, nil
}

func (service userInterestsService) UnFollowInterest(ctx *gin.Context, interestID, userID uuid.UUID) *golaerror.Error {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsService").WithField("method", "GetExploreInterests")
	log.Info("unfollowing user followed interests")

	err := service.repository.UnFollowInterest(ctx, interestID, userID)
	if err != nil {
		log.Errorf("unable to unfollow interest %v", err)
		return &constants.InternalServerError
	}

	return nil
}

func NewUserInterestsService(repository repository.UserInterestsRepository, services service.AwsServices) UserInterestsService {
	return userInterestsService{repository: repository, awsService: services}
}
