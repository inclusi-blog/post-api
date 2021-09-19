package repository

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	storyModels "post-api/story/models"
	"post-api/user-profile/models"
)

const (
	GetFollowedInterests = "select json_agg(json_build_object('id', interest_id, 'name', i.name)) as interests from user_interests inner join interests i on user_interests.interest_id = i.id where user_id = $1"
	MapInterests         = "insert into user_interests(user_id, interest_id) values ($1, (select id from interests where id = $2))"
	GetExploreInterests  = "select i.name as category, (select json_agg(jsonb_build_object('id', ii.id, 'name', ii.name, 'cover_pic', ii.cover_pic, 'is_followed', ii.id = ui.interest_id)) from interests ii left join (select * from user_interests where user_id = $1) as ui on ii.id = ui.interest_id where ii.id in (select interest_id from category_x_interests where category_id = cxi.category_id) and ii.cover_pic is not null) as interests from interests i inner join category_x_interests cxi on i.id = cxi.category_id group by category_id, i.name"
)

type UserInterestsRepository interface {
	GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*storyModels.JSONString, error)
	FollowInterest(ctx *gin.Context, interestID, userID uuid.UUID) error
	GetExploreInterests(ctx *gin.Context, userID uuid.UUID) ([]models.ExploreInterest, error)
}

type userInterestsRepository struct {
	db *sqlx.DB
}

func (repository userInterestsRepository) GetFollowedInterest(ctx *gin.Context, userId uuid.UUID) (*storyModels.JSONString, error) {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsRepository").WithField("method", "GetFollowedInterest")

	type followedInterests struct {
		Interests storyModels.JSONString `json:"interests"`
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

func (repository userInterestsRepository) GetExploreInterests(ctx *gin.Context, userID uuid.UUID) ([]models.ExploreInterest, error) {
	log := logging.GetLogger(ctx).WithField("class", "UserInterestsRepository").WithField("method", "GetExploreInterests")

	var interests []models.ExploreInterest
	err := repository.db.SelectContext(ctx, &interests, GetExploreInterests, userID)
	if err != nil {
		log.Errorf("unable to get explore interests %v", err)
		return nil, err
	}

	return interests, nil
}

func NewUserInterestsRepository(db *sqlx.DB) UserInterestsRepository {
	return userInterestsRepository{db: db}
}
