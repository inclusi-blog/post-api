package repository

//go:generate mockgen -source=interests_repository.go -destination=./../mocks/mock_interests_repositiry.go -package=mocks

import (
	"context"
	"fmt"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
	"post-api/models/db"
)

type InterestsRepository interface {
	GetInterests(ctx context.Context, searchKeyword string) ([]db.Interest, error)
}

type interestsRepository struct {
	db *sqlx.DB
}

const (
	GetInterests = "SELECT ID, NAME FROM INTERESTS WHERE NAME LIKE '%%%s%%'"
)

func (repository interestsRepository) GetInterests(ctx context.Context, searchKeyword string) ([]db.Interest, error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsRepository").WithField("method", "GetInterests")

	logger.Info("fetching over all interests")

	var interests []db.Interest

	key := fmt.Sprintf(GetInterests, searchKeyword)
	err := repository.db.SelectContext(ctx, &interests, key)

	if err != nil {
		logger.Errorf("Error occurred while fetching over all interests from repository %v", err)
		return nil, err
	}

	logger.Info("Successfully fetched interests from db")

	return interests, nil
}

func NewInterestRepository(db *sqlx.DB) InterestsRepository {
	return interestsRepository{db: db}
}
