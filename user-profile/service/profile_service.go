package service

import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"post-api/user-profile/constants"
	"post-api/user-profile/models"
	"post-api/user-profile/repository"
)

type profileService struct {
	repository repository.ProfileRepository
}

type ProfileService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (models.Profile, *golaerror.Error)
}

func (service profileService) GetProfile(ctx context.Context, userID uuid.UUID) (models.Profile, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "ProfileService").WithField("method", "GetProfile")
	logger.Infof("calling db to fetch profile details for user %v", userID)

	details, err := service.repository.GetDetails(ctx, userID)
	if err != nil {
		return models.Profile{}, &constants.InternalServerError
	}
	logger.Info("successfully fetched user details")

	return details, nil
}

func NewProfileService(repository repository.ProfileRepository) ProfileService {
	return profileService{repository: repository}
}
