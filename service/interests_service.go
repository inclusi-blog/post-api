package service

//go:generate mockgen -source=interests_service.go -destination=./../mocks/mock_interests_service.go -package=mocks
import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/constants"
	"post-api/models/db"
	"post-api/repository"
)

type InterestsService interface {
	GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, *golaerror.Error)
}

type interestsService struct {
	repository repository.InterestsRepository
}

func (service interestsService) GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsService").WithField("method", "GetInterests")

	logger.Info("Calling repository to get all interests")
	interests, err := service.repository.GetInterests(ctx, searchKeyword, selectedTags)
	if err != nil {
		logger.Errorf("error occurred while fetching over all interests from interest repository %v", err)
		return nil, &constants.PostServiceFailureError
	}

	if len(interests) == 0 {
		logger.Error("no results found for interests")
		return nil, &constants.NoInterestsFoundError
	}

	logger.Info("successfully fetched interests")
	return interests, nil
}

func NewInterestsService(interestsRepository repository.InterestsRepository) InterestsService {
	return interestsService{repository: interestsRepository}
}
