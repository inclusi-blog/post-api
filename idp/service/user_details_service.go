package service

import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"post-api/idp/constants"
	"post-api/idp/models/request"
	"post-api/idp/repository"
)

type UserDetailsService interface {
	UpdateUserDetails(ctx context.Context, userID uuid.UUID, update request.UserDetailsUpdate) *golaerror.Error
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
		if !usernameAvailabilityResponse.IsAvailable {
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
		err := service.repository.UpdateName(ctx, update.Username, userID)
		if err != nil {
			logger.Error("unable to update username for user %v", userID)
			return &constants.NameUpdateError
		}
	}
	if update.About != "" {
		err := service.repository.UpdateAbout(ctx, update.Username, userID)
		if err != nil {
			logger.Error("unable to update username for user %v", userID)
			return &constants.AboutUpdateError
		}
	}

	return nil
}

func NewUserDetailsService(repository repository.UserDetailsRepository, service UserRegistrationService) UserDetailsService {
	return userDetailsService{repository: repository, userRegistrationService: service}
}
