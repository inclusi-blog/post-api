package repository

//go:generate mockgen -source=interests_repository.go -destination=./../mocks/mock_interests_repositiry.go -package=mocks

import (
	"context"
	"errors"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"post-api/constants"
	"post-api/models/db"
)

type InterestsRepository interface {
	GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, error)
}

type interestsRepository struct {
	db neo4j.Session
}

const (
	GetInterestsWithoutSelectedTags = "match (interest:Interest) where NOT interest.name IN $selectedInterests and interest.name =~ '(?i).*'+ $searchKeyword +'.*' return interest.name as name"
)

func (repository interestsRepository) GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsRepository").WithField("method", "GetInterests")

	logger.Info("fetching over all interests")

	var interests []db.Interest

	arg := map[string]interface{}{
		"selectedInterests": selectedTags,
		"searchKeyword":     searchKeyword,
	}

	result, err := repository.db.Run(GetInterestsWithoutSelectedTags, arg)

	if err != nil {
		logger.Errorf("Error occurred while naming in query params for get interests %v", err)
		return []db.Interest{}, err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error occurred while getting summary  get interests %v", err)
		return []db.Interest{}, err
	}

	for result.Next() {
		value, isPresent := result.Record().Get("name")
		if !isPresent {
			logger.Errorf("Error occurred while binding In query params for get interests %v", err)
			return []db.Interest{}, err
		}
		if value != nil {
			interest := value.(string)
			interests = append(interests, db.Interest{Name: interest})
		}
	}

	if len(interests) == 0 {
		logger.Errorf("Error no interests found")
		return []db.Interest{}, errors.New(constants.NoInterestsFoundCode)
	}
	logger.Info("Successfully fetched interests from db")

	return interests, nil
}

func NewInterestRepository(db neo4j.Session) InterestsRepository {
	return interestsRepository{db: db}
}
