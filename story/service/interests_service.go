package service

//go:generate mockgen -source=interests_service.go -destination=./../mocks/mock_interests_service.go -package=mocks
import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
	"post-api/story/constants"
	"post-api/story/models/db"
	"post-api/story/models/response"
	"post-api/story/repository"
)

type InterestsService interface {
	GetInterests(ctx context.Context) ([]db.Interests, *golaerror.Error)
	GetFollowCount(ctx context.Context, interestName string, userID uuid.UUID) (response.InterestCountDetails, *golaerror.Error)
}

type interestsService struct {
	repository repository.InterestsRepository
}

func (service interestsService) GetInterests(ctx context.Context) ([]db.Interests, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsService").WithField("method", "GetInterests")

	logger.Info("Calling repository to get all interests")
	interests, err := service.repository.GetInterests(ctx)
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

func (service interestsService) GetFollowCount(ctx context.Context, interestName string, userID uuid.UUID) (response.InterestCountDetails, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsService").WithField("method", "GetFollowCount")
	logger.Info("Calling repository to get all interests")

	details, err := service.repository.GetFollowCount(ctx, interestName, userID)
	if err != nil {
		logger.Errorf("unable to fetch interest details %v", err)
		if err == sql.ErrNoRows {
			logger.Error("no interest found")
			return response.InterestCountDetails{}, &constants.NoInterestsFoundError
		}
		return response.InterestCountDetails{}, &constants.InternalServerError
	}

	return details, nil
}

func NewInterestsService(interestsRepository repository.InterestsRepository) InterestsService {
	return interestsService{repository: interestsRepository}
}
