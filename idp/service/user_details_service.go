package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
	"post-api/idp/constants"
	"post-api/idp/models/request"
	"post-api/idp/repository"
)

type UserDetailsService interface {
	UpdateUserDetails(ctx context.Context, userID uuid.UUID, update request.UserDetailsUpdate) *golaerror.Error
	UpdateProfileImage(ctx context.Context, avatarKey string, userID uuid.UUID) *golaerror.Error
}

type userDetailsService struct {
	repository              repository.UserDetailsRepository
	userRegistrationService UserRegistrationService
}

func (service userDetailsService) UpdateUserDetails(ctx context.Context, userID uuid.UUID, update request.UserDetailsUpdate) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsService").WithField("method", "UpdateUserDetails")
	logger.Info("updating user details")

	if update.Username != "" {
		usernameAvailabilityResponse, err := service.userRegistrationService.IsUsernameRegistered(ctx, request.UsernameAvailabilityRequest{Username: update.Username})
		if err != nil {
			logger.Error("unable to get username availability")
			return &constants.InternalServerError
		}
		if usernameAvailabilityResponse.IsAvailable {
			logger.Error("same username already found")
			return &constants.UsernameAlreadyPresentError
		}

		updateErr := service.repository.UpdateUsername(ctx, update.Username, userID)
		if updateErr != nil {
			logger.Errorf("unable to update username for user %v. Error %v", userID, updateErr)
			return &constants.UsernameUpdateError
		}
	}
	if update.Name != "" {
		err := service.repository.UpdateName(ctx, update.Name, userID)
		if err != nil {
			logger.Error("unable to update username for user %v", userID)
			return &constants.NameUpdateError
		}
	}
	if update.About != "" {
		err := service.repository.UpdateAbout(ctx, update.About, userID)
		if err != nil {
			logger.Error("unable to update username for user %v", userID)
			return &constants.AboutUpdateError
		}
	}
	if update.FacebookURL != "" {
		err := service.repository.UpdateFacebookURL(ctx, update.FacebookURL, userID)
		if err != nil {
			logger.Error("unable to update facebook url for user %v", userID)
			return &constants.SocialUpdateError
		}
	}
	if update.LinkedInURL != "" {
		err := service.repository.UpdateLinkedInURL(ctx, update.LinkedInURL, userID)
		if err != nil {
			logger.Error("unable to update linkedin url for user %v", userID)
			return &constants.SocialUpdateError
		}
	}
	if update.TwitterURL != "" {
		err := service.repository.UpdateTwitterURL(ctx, update.TwitterURL, userID)
		if err != nil {
			logger.Error("unable to update twitter url for user %v", userID)
			return &constants.SocialUpdateError
		}
	}

	return nil
}

func (service userDetailsService) UpdateProfileImage(ctx context.Context, avatarKey string, userID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "UserDetailsService").WithField("method", "UpdateUserDetails")
	logger.Info("updating user avatar")

	err := service.repository.UpdateProfileImage(ctx, avatarKey, userID)
	if err != nil {
		logger.Errorf("unable to update user profile avatar %v", err)
		return &constants.UnableToUpdateAvatarError
	}

	return nil
}

func NewUserDetailsService(repository repository.UserDetailsRepository, service UserRegistrationService) UserDetailsService {
	return userDetailsService{repository: repository, userRegistrationService: service}
}
