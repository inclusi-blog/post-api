package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/jmoiron/sqlx"
	"post-api/user-profile/models"
)

type profileRepository struct {
	db *sqlx.DB
}

type ProfileRepository interface {
	GetDetails(ctx context.Context, id uuid.UUID) (models.Profile, error)
	FetchProfileAvatar(ctx context.Context, id uuid.UUID) (string, error)
	FollowUser(ctx context.Context, userID, followingID uuid.UUID) error
	UnFollowUser(ctx context.Context, userID, followingID uuid.UUID) error
	BlockUser(ctx context.Context, userID, toBlockID uuid.UUID) error
}

const (
	GetProfile   = "select users.id as id, name, username, email, about, avatar as profile_pic, facebook as facebook_url, linkedin as linked_in_url, twitter as twitter_url from users left join social_links sl on users.id = sl.user_id where users.id = $1"
	GetAvatar    = "select avatar from users where id = $1"
	FollowUser   = "insert into followings(follower_id, following_id)values($1, $2)"
	UnfollowUser = "delete from followings where follower_id = $1 and following_id = $2"
	BlockUser    = "insert into user_blocks(blocked_id, blocked_by) values ($1, $2)"
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

func (repository profileRepository) FetchProfileAvatar(ctx context.Context, id uuid.UUID) (string, error) {
	logger := logging.GetLogger(ctx).WithField("class", "ProfileRepository").WithField("method", "FetchProfileAvatar")
	var avatar *string
	err := repository.db.GetContext(ctx, &avatar, GetAvatar, id)

	if err != nil {
		logger.Errorf("unable to fetch user profile avatar %v", err)
		return "", err
	}

	if avatar == nil {
		return "", sql.ErrNoRows
	}

	return *avatar, nil
}

func (repository profileRepository) FollowUser(ctx context.Context, userID, followingID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "ProfileRepository").WithField("method", "FollowUser")

	_, err := repository.db.ExecContext(ctx, FollowUser, userID, followingID)
	if err != nil {
		logger.Errorf("unable to follow user %v", err)
		return err
	}

	return nil
}

func (repository profileRepository) UnFollowUser(ctx context.Context, userID, followingID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "ProfileRepository").WithField("method", "FollowUser")

	_, err := repository.db.ExecContext(ctx, UnfollowUser, userID, followingID)
	if err != nil {
		logger.Errorf("unable to unfollow user %v", err)
		return err
	}

	return nil
}

func (repository profileRepository) BlockUser(ctx context.Context, userID, toBlockID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "ProfileRepository").WithField("method", "BlockUser")
	_, err := repository.db.ExecContext(ctx, BlockUser, toBlockID, userID)
	if err != nil {
		logger.Errorf("unable to block %v author by author %v. Error %v", toBlockID, userID, err)
		return err
	}

	return nil
}

func NewProfileRepository(db *sqlx.DB) ProfileRepository {
	return profileRepository{db: db}
}
