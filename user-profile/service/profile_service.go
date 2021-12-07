package service

import (
	"context"
	"database/sql"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"post-api/service"
	"post-api/user-profile/constants"
	"post-api/user-profile/models"
	"post-api/user-profile/repository"
	"time"
)

type profileService struct {
	repository  repository.ProfileRepository
	awsServices service.AwsServices
}

type ProfileService interface {
	GetProfile(ctx context.Context, userID uuid.UUID) (models.Profile, *golaerror.Error)
	FetchProfileAvatar(ctx context.Context, id uuid.UUID) (string, *golaerror.Error)
	FollowUser(ctx context.Context, userID, followingID uuid.UUID) *golaerror.Error
	UnFollowUser(ctx context.Context, userID, followingID uuid.UUID) *golaerror.Error
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

func (service profileService) FetchProfileAvatar(ctx context.Context, id uuid.UUID) (string, *golaerror.Error) {
	avatar, err := service.repository.FetchProfileAvatar(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &constants.NoUserFoundError
		}
		return "", &constants.InternalServerError
	}

	var url string
	if avatar != "" {
		url, err = service.awsServices.GetObjectInS3(avatar, time.Hour*time.Duration(6))
		if err != nil {
			return "", &constants.InternalServerError
		}
	}

	return url, nil
}

func (service profileService) FollowUser(ctx context.Context, userID, followingID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "ProfileService").WithField("method", "FollowUser")
	err := service.repository.FollowUser(ctx, userID, followingID)
	if err != nil {
		logger.Errorf("unable to follow user %v", err)
		return &constants.InternalServerError
	}
	logger.Info("successfully followed")
	return nil
}

func (service profileService) UnFollowUser(ctx context.Context, userID, followingID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "ProfileService").WithField("method", "UnFollowUser")
	err := service.repository.UnFollowUser(ctx, userID, followingID)
	if err != nil {
		logger.Errorf("unable to follow user %v", err)
		return &constants.InternalServerError
	}
	logger.Info("successfully followed")
	return nil
}

func NewProfileService(repository repository.ProfileRepository, services service.AwsServices) ProfileService {
	return profileService{
		repository:  repository,
		awsServices: services,
	}
}
