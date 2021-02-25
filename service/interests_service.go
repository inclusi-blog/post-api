package service

//go:generate mockgen -source=interests_service.go -destination=./../mocks/mock_interests_service.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/clients/user_profile"
	"post-api/constants"
	"post-api/mapper"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/response"
	"post-api/repository"
)

type InterestsService interface {
	GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, *golaerror.Error)
	GetExploreCategoriesAndInterests(ctx context.Context) ([]response.CategoryAndInterest, *golaerror.Error)
}

type interestsService struct {
	repository        repository.InterestsRepository
	userProfileClient user_profile.Client
	mapper            mapper.InterestsMapper
}

func (service interestsService) GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsService").WithField("method", "GetInterests")

	logger.Info("Calling repository to get all interests")
	interests, err := service.repository.GetInterests(ctx, searchKeyword, selectedTags)
	if err != nil {
		if err.Error() == constants.NoInterestsFoundCode {
			logger.Errorf("No interests found for keyword, Error %v", err)
			return nil, &constants.NoInterestsFoundError
		}
		logger.Errorf("error occurred while fetching over all interests from interest repository %v", err)
		return nil, &constants.PostServiceFailureError
	}

	logger.Info("successfully fetched interests")
	return interests, nil
}

func (service interestsService) GetExploreCategoriesAndInterests(ctx context.Context) ([]response.CategoryAndInterest, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsService").WithField("method", "GetExploreCategoriesAndInterests")
	logger.Info("Calling repository to fetch explore categories and interest")

	var apiChannel = make(chan models.APIChannel)
	go service.getUserFollowingInterests(ctx, apiChannel)

	userFollowingInterestsChannel := <-apiChannel
	userProfileErr := userFollowingInterestsChannel.Error.(*golaerror.Error)
	if userProfileErr != nil {
		return nil, &constants.InternalServerError
	}
	userFollowingInterests := userFollowingInterestsChannel.Response.([]string)

	categoriesAndInterests, err := service.repository.FetchCategoriesAndInterests(ctx)
	if err != nil {
		if err.Error() == constants.NoInterestsAndCategoriesCode {
			logger.Errorf("No interests and categories found, Error %v", err)
			return nil, &constants.NoInterestsAndCategoriesErr
		}
		logger.Errorf("Error occurred while fetching explore categories and interests, Error %v", err)
		return nil, &constants.PostServiceFailureError
	}

	if categoriesAndInterests == nil || len(categoriesAndInterests) == 0 {
		logger.Error("No interests and categories found")
		return nil, &constants.NoInterestsAndCategoriesErr
	}

	logger.Info("successfully fetched categories and interests")

	categoriesAndInterestsWithUserFollowStatus := service.mapper.MapUserFollowedInterest(ctx, categoriesAndInterests, userFollowingInterests)

	return categoriesAndInterestsWithUserFollowStatus, nil
}

func (service interestsService) getUserFollowingInterests(ctx context.Context, ch chan<- models.APIChannel) {
	logger := logging.GetLogger(ctx).WithField("class", "").WithField("method", "GetExploreCategoriesAndInterests").WithField("operation", "getUserFollowingInterests")
	logger.Info("Fetching user following interests from user profile")

	userprofileInterets, err := service.userProfileClient.FetchUserFollowingInterests(ctx)

	var data models.APIChannel
	data.Error = err
	data.Response = userprofileInterets

	ch <- data
}

func NewInterestsService(interestsRepository repository.InterestsRepository, client user_profile.Client, interestsMapper mapper.InterestsMapper) InterestsService {
	return interestsService{
		repository:        interestsRepository,
		userProfileClient: client,
		mapper:            interestsMapper,
	}
}
