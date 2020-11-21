package repository

//go:generate mockgen -source=posts_repository.go -destination=./../mocks/mock_posts_repository.go -package=mocks

import (
	"context"
	"database/sql"
	"post-api/models/db"

	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, post db.PublishPost) (int64, error)
	// GetLikesIdByPost(ctx context.Context, postID string, userID string) (string, error)
	// SaveUserToLikedBy(postID string, userID string, ctx context.Context) error
	// RemoveUserFromLikedBy(postID string, userID string, ctx context.Context) error
	GetLikeCountByPost(ctx context.Context, postID string) (string, error)
	AppendOrRemoveUserFromLikedBy(postID string, userID string, ctx context.Context) error
}

type postRepository struct {
	db *sqlx.DB
}

const (
	PublishPost           = "INSERT INTO POSTS (puid, user_id, post_data, read_time, view_count) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	GetLikeIdByUser       = "SELECT id  FROM likes WHERE post_id=$1 AND $2=ANY(liked_by)"
	UpdateLikedBy         = "UPDATE likes SET liked_by = array_append(liked_by,$1) WHERE post_id=$2"
	RemoveLikedBy         = "UPDATE likes SET liked_by = array_remove(liked_by,$1) WHERE post_id=$2"
	GetLikedByCount       = "SELECT array_length(liked_by, 1) FROM likes WHERE post_id=$1 "
	UpdateOrRemoveLikedBy = "UPDATE likes SET liked_by = case when (SELECT count(id) as id  FROM likes WHERE post_id=$1 AND $2=ANY(liked_by)) = '1' then array_remove(liked_by, $2) else array_append(liked_by, $2) end where post_id = $1"
)

func (repository postRepository) CreatePost(ctx context.Context, post db.PublishPost) (int64, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "CreatePost")

	logger.Infof("Publishing the post in posts table for post postID %v", post.PUID)

	var postID int64
	err := repository.db.QueryRowContext(ctx, PublishPost, post.PUID, post.UserID, post.PostData, post.ReadTime, post.ViewCount).Scan(&postID)

	if err != nil {
		logger.Errorf("Error occurred while publishing user post in posts table %v", err)
		return 0, err
	}

	logger.Infof("Successfully posted the post in posts table for post postID %v", post.PUID)

	return postID, nil
}

// func (repository postRepository) GetLikesIdByPost(ctx context.Context, postID string, userID string) (string, error) {
// 	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "GetLikesIdByPost")

// 	logger.Infof("Fetching like id from likes table for the given post id %v %v", postID, userID)

// 	var likeID db.LikedByRes

// 	err := repository.db.GetContext(ctx, &likeID, GetLikeIdByUser, postID, userID)

// 	if err != nil && err.Error() != sql.ErrNoRows.Error() {
// 		logger.Errorf("Error occurred while fetching likeID from likes table %v", err.Error())
// 		return likeID.LikedByID, err
// 	}

// 	logger.Infof("Successfully fetching like ID from likes table for given post id %v", postID)

// 	return likeID.LikedByID, nil
// }

// func (repository postRepository) SaveUserToLikedBy(postID string, userID string, ctx context.Context) error {
// 	logger := logging.GetLogger(ctx)
// 	log := logger.WithField("class", "PostsRepository").WithField("method", "SaveUserToLikedBy")

// 	log.Infof("updating the existing liked by in like table for post %v", postID)

// 	_, err := repository.db.ExecContext(ctx, UpdateLikedBy, userID, postID)

// 	if err != nil {
// 		log.Errorf("Error occurred while updating liked by in likes table %v", err)
// 		return err
// 	}

// 	return nil
// }

// func (repository postRepository) RemoveUserFromLikedBy(postID string, userID string, ctx context.Context) error {
// 	logger := logging.GetLogger(ctx)
// 	log := logger.WithField("class", "PostsRepository").WithField("method", "RemoveUserFromLikedBy")

// 	log.Infof("updating the existing liked by in like table for post %v", postID)

// 	_, err := repository.db.ExecContext(ctx, RemoveLikedBy, userID, postID)

// 	if err != nil {
// 		log.Errorf("Error occurred while updating liked by in likes table %v", err)
// 		return err
// 	}

// 	return nil
// }

func (repository postRepository) GetLikeCountByPost(ctx context.Context, postID string) (string, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "GetLikeCountByPost")

	logger.Infof("Fetching likedby count from likes table for the given post id %v %v", postID)

	var likeCount sql.NullString

	err := repository.db.GetContext(ctx, &likeCount, GetLikedByCount, postID)

	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		logger.Errorf("Error occurred while fetching likedby count from likes table %v", err.Error())
		return likeCount.String, err
	}

	logger.Infof("Successfully fetching likedby count from likes table for given post id %v", postID)

	return likeCount.String, nil
}

func (repository postRepository) AppendOrRemoveUserFromLikedBy(postID string, userID string, ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("class", "PostsRepository").WithField("method", "AppendOrRemoveUserFromLikedBy")

	log.Infof("updating the likedby in like table for post %v", postID)

	_, err := repository.db.ExecContext(ctx, UpdateOrRemoveLikedBy, postID, userID)

	if err != nil {
		log.Errorf("Error occurred while updating liked by in likes table %v", err)
		return err
	}

	return nil
}
func NewPostsRepository(db *sqlx.DB) PostsRepository {
	return postRepository{db: db}
}
