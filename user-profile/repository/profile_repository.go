package repository

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"post-api/user-profile/models"
)

type profileRepository struct {
	db *sqlx.DB
}

type ProfileRepository interface {
	GetDetails(ctx context.Context, id uuid.UUID) (models.Profile, error)
}

const (
	GetProfile = "select users.id as id, name, username, email, about, avatar as profile_pic, facebook as facebook_url, linkedin as linked_in_url, twitter as twitter_url from users inner join social_links sl on users.id = sl.user_id where users.id = $1"
)

func (repository profileRepository) GetDetails(ctx context.Context, id uuid.UUID) (models.Profile, error) {
	logger := logging.GetLogger(ctx).WithField("class", "ProfileRepository").WithField("method", "GetDetails")
	var details models.Profile
	err := repository.db.GetContext(ctx, &details, GetProfile, id)
	if err != nil {
		logger.Errorf("unable to get profile %v", err)
		return models.Profile{}, err
	}

	return details, nil
}

func NewProfileRepository(db *sqlx.DB) ProfileRepository {
	return profileRepository{db: db}
}
