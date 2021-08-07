package repository

//go:generate mockgen -source=preview_posts_repository.go -destination=./../mocks/mock_preview_posts_repository.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
	"post-api/models/db"
)

type PreviewPostsRepository interface {
	SavePreview(ctx context.Context, post db.PreviewPost) (int64, error)
}

type previewPostsRepository struct {
	db *sqlx.DB
}

const (
	SavePreviewPost = "INSERT INTO preview_posts (post_id, title, tagline, preview_image, like_count, comment_count, view_time) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
)

func (repository previewPostsRepository) SavePreview(ctx context.Context, post db.PreviewPost) (int64, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PreviewPostsRepository").WithField("method", "SavePreview")

	id := post.PostID
	logger.Infof("Inserting new preview post for post %v", id)
	var previewID int64
	err := repository.db.QueryRowContext(ctx, SavePreviewPost, id, post.Title, post.Tagline, post.PreviewImage, post.LikeCount, post.CommentCount, post.ViewTime).Scan(&previewID)

	if err != nil {
		logger.Errorf("Error occurred while inserting new preview post for post id %v .%v", id, err)
		return 0, err
	}

	logger.Infof("Successfully saved preview post for post id %v", id)

	return previewID, nil
}

func NewPreviewPostsRepository(db *sqlx.DB) PreviewPostsRepository {
	return previewPostsRepository{db: db}
}
