package repository

//go:generate mockgen -source=posts_repository.go -destination=./../mocks/mock_posts_repository.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
	"post-api/models/db"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, post db.PublishPost) (int64, error)
}

type postRepository struct {
	db *sqlx.DB
}

const (
	PublishPost = "INSERT INTO POSTS (puid, user_id, post_data, read_time, view_count) VALUES ($1, $2, $3, $4, $5) RETURNING id"
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

func NewPostsRepository(db *sqlx.DB) PostsRepository {
	return postRepository{db: db}
}
