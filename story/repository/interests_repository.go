package repository

//go:generate mockgen -source=interests_repository.go -destination=./../mocks/mock_interests_repositiry.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"post-api/story/models/db"
	"post-api/story/models/response"
)

type InterestsRepository interface {
	GetInterests(ctx context.Context) ([]db.Interests, error)
	GetInterestIDs(ctx context.Context, interestNames []string) ([]uuid.UUID, error)
	GetInterestsForName(ctx context.Context, interestNames []string) ([]db.Interests, error)
	GetFollowCount(ctx context.Context, interestID, userID uuid.UUID) (response.InterestCountDetails, error)
}

type interestsRepository struct {
	db *sqlx.DB
}

const (
	GetInterests            = "select id, name from interests"
	GetInterestIDs          = "SELECT id from interests where name in (?)"
	GetInterestsForNames    = "select id, name from interests where name in (?)"
	GetInterestsFollowCount = "select count(*) as followers_count, $1 in (user_id) as is_followed from user_interests where interest_id = $2 group by user_interests.user_id"
)

func (repository interestsRepository) GetInterests(ctx context.Context) ([]db.Interests, error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsRepository").WithField("method", "GetInterests")
	logger.Info("fetching over all interests")

	var interests []db.Interests
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
	if err != nil {
		logger.Errorf("error occurred while binding interest ids %v", err)
		return nil, err
	}
	return interestsIDs, nil
}

func (repository interestsRepository) GetInterestsForName(ctx context.Context, interestNames []string) ([]db.Interests, error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsRepository").WithField("method", "GetInterests")
	logger.Info("fetching over all interests")

	var interests []db.Interests
	query, args, err := sqlx.In(GetInterestsForNames, interestNames)

	query = repository.db.Rebind(query)
	rows, err := repository.db.Query(query, args...)
	if err != nil {
		logger.Errorf("unable to fetch interest ids %v", err)
		return nil, err
	}
	for rows.Next() {
		var interest db.Interests
		err = rows.Scan(&interest.ID, &interest.Name)
		interests = append(interests, interest)
	}
	if err != nil {
		logger.Errorf("error occurred while binding interest ids %v", err)
		return nil, err
	}
	return interests, nil
}

func (repository interestsRepository) GetFollowCount(ctx context.Context, interestID, userID uuid.UUID) (response.InterestCountDetails, error) {
	logger := logging.GetLogger(ctx).WithField("class", "InterestsRepository").WithField("method", "GetFollowCount")
	logger.Info("fetching over all interests")

	var interestDetails response.InterestCountDetails
	err := repository.db.GetContext(ctx, &interestDetails, GetInterestsFollowCount, userID, interestID)
	if err != nil {
		logger.Errorf("unable to fetch interest details %v", err)
		return response.InterestCountDetails{}, err
	}

	return interestDetails, nil
}

func NewInterestRepository(db *sqlx.DB) InterestsRepository {
	return interestsRepository{db: db}
}
