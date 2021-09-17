package repository

import (
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"post-api/story/models"
)

const (
	GetFollowedInterests = "select json_agg(json_build_object('id', interest_id, 'name', i.name)) as interests from user_interests inner join interests i on user_interests.interest_id = i.id where user_id = $1"
)

type UserInterestsRepository interface {
	GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*models.JSONString, error)
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

func NewUserInterestsRepository(db *sqlx.DB) UserInterestsRepository {
	return userInterestsRepository{db: db}
}
