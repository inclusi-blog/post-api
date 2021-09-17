package repository

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"post-api/story/models"
)

const (
	GetFollowedInterests = "select json_agg(json_build_object('id', interest_id, 'name', i.name)) as interests from user_interests inner join interests i on user_interests.interest_id = i.id where user_id = $1"
	MapInterests         = "insert into user_interests(user_id, interest_id) values ($1, (select id from interests where id = $2))"
)

type UserInterestsRepository interface {
	GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*models.JSONString, error)
	FollowInterest(ctx *gin.Context, interestID, userID uuid.UUID) error
}

type userInterestsRepository struct {
	db *sqlx.DB
}

func (repository userInterestsRepository) GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*models.JSONString, error) {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsRepository").WithField("method", "GetFollowedInterest")

	type followedInterests struct {
		Interests models.JSONString `json:"interests"`
	}
	var interests followedInterests
	err := repository.db.GetContext(ctx, &interests, GetFollowedInterests, userId)
	if err != nil {
		log.Errorf("unable to get followed interests %v", err)
		return nil, err
	}

	return &interests.Interests, nil
}

func (repository userInterestsRepository) FollowInterest(ctx *gin.Context, interestID, userID uuid.UUID) error {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsRepository").WithField("method", "FollowInterest")

	result, err := repository.db.ExecContext(ctx, MapInterests, userID, interestID)
	if err != nil {
		log.Errorf("unable to follow interest %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Errorf("unable to get affected rows %v", err)
		return err
	}

	if rowsAffected == 1 {
		return nil
	}

	return errors.New("internal error occurred")
}

func NewUserInterestsRepository(db *sqlx.DB) UserInterestsRepository {
	return userInterestsRepository{db: db}
}
