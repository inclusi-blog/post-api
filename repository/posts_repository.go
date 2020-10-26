package repository

//go:generate mockgen -source=posts_repository.go -destination=./../mocks/mock_posts_repository.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
	"post-api/models/db"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, post db.PublishPost) error
}

type postRepository struct {
	db *sqlx.DB
}

const (
	PublishPost = "INSERT INTO POSTS (puid, user_id, post_data, title_data, read_time, view_count) VALUES ($1, $2, $3, $4, $5, $6)"
)

func (repository postRepository) CreatePost(ctx context.Context, post db.PublishPost) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "CreatePost")

	logger.Infof("Publishing the post in posts table for post id %v", post.PUID)

	result, err := repository.db.ExecContext(ctx, PublishPost, post.PUID, post.UserID, post.PostData, post.TitleData, post.ReadTime, post.ViewCount)

	if err != nil {
		logger.Errorf("Error occurred while publishing user post in posts table %v", err)
		return err
	}

	if affectedRows, err := result.RowsAffected(); affectedRows != 1 || err != nil {
		logger.Errorf("Error occurred while inserting new post in posts table %v", err)
		return err
	}

	logger.Infof("Successfully posted the post in posts table for post id %v", post.PUID)

	return nil
}

func NewPostsRepository(db *sqlx.DB) PostsRepository {
	return postRepository{db: db}
}
