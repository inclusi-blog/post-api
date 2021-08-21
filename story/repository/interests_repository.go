package repository

//go:generate mockgen -source=interests_repository.go -destination=./../mocks/mock_interests_repositiry.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type InterestsRepository interface {
	GetInterests(ctx context.Context) ([]string, error)
	GetInterestIDs(ctx context.Context, interestNames []string) ([]uuid.UUID, error)
}

type interestsRepository struct {
	db *sqlx.DB
}

const (
	GetInterests = "select name from interests"
	GetInterestIDs = "SELECT id from interests where name in (?)"
)

func (repository interestsRepository) GetInterests(ctx context.Context) ([]string, error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsRepository").WithField("method", "GetInterests")
	logger.Info("fetching over all interests")

	var interests []string
	err := repository.db.SelectContext(ctx, &interests, GetInterests)
	if err != nil {
		logger.Errorf("Error occurred while fetching over all interests from repository %v", err)
		return nil, err
	}

	logger.Info("Successfully fetched interests from db")
	return interests, nil
}

func (repository interestsRepository) GetInterestIDs(ctx context.Context, interestNames []string) ([]uuid.UUID, error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsRepository").WithField("method", "GetInterests")
	logger.Info("fetching over all interests")

	var interestsIDs []uuid.UUID
	query, args, err := sqlx.In(GetInterestIDs, interestNames)

	query = repository.db.Rebind(query)
	rows, err := repository.db.Query(query, args...)
	if err != nil {
		logger.Errorf("unable to fetch interest ids %v", err)
		return nil, err
	}
	for rows.Next() {
		var id uuid.UUID
		err = rows.Scan(&id)
		interestsIDs = append(interestsIDs, id)
	}
	if err != nil{
		logger.Errorf("error occurred while binding interest ids %v", err)
		return nil, err
	}
	return interestsIDs, nil
}

func NewInterestRepository(db *sqlx.DB) InterestsRepository {
	return interestsRepository{db: db}
}
