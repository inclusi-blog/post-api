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
	GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, error)
}

type interestsRepository struct {
	db *sqlx.DB
}

const (
	GetInterestsWithoutSelectedTags = "SELECT ID, NAME FROM INTERESTS WHERE NAME LIKE '%%%s%%' AND NAME NOT IN (:tags)"
	GetInterests                    = "SELECT ID, NAME FROM INTERESTS WHERE NAME LIKE '%%%s%%'"
)

func (repository interestsRepository) GetInterests(ctx context.Context, searchKeyword string, selectedTags []string) ([]db.Interest, error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsRepository").WithField("method", "GetInterests")

	logger.Info("fetching over all interests")

	var interests []db.Interest
	selectionQuery := GetInterestsWithoutSelectedTags

	if selectedTags == nil || len(selectedTags) == 0 {
		logger.Info("Fetching all the interests related to search keyword")
		selectionQuery = GetInterests
	}

	arg := map[string]interface{}{
		"tags": selectedTags,
	}

	updatedQuery := fmt.Sprintf(selectionQuery, searchKeyword)

	query, args, err := sqlx.Named(updatedQuery, arg)
	if err != nil {
		logger.Errorf("Error occurred while naming in query params for get interests %v", err)
		return []db.Interest{}, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		logger.Errorf("Error occurred while binding In query params for get interests %v", err)
		return []db.Interest{}, err
	}

	query = repository.db.Rebind(query)
	err = repository.db.SelectContext(ctx, &interests, query, args...)

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
