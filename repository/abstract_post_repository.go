package repository

//go:generate mockgen -source=abstract_post_repository.go -destination=./../mocks/mock_abstract_post_repository.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"post-api/helper"
	"post-api/models/db"
)

type AbstractPostRepository interface {
	Save(ctx context.Context, txn helper.Transaction, post db.AbstractPost) (uuid.UUID, error)
}

type abstractPostRepository struct {
	db *sqlx.DB
}

const (
	SavePreviewPost = "INSERT INTO abstract_post (id, title, tagline, preview_image, view_time, post_id) VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5) RETURNING id"
)

func (repository abstractPostRepository) Save(ctx context.Context, txn helper.Transaction, post db.AbstractPost) (uuid.UUID, error) {
	logger := logging.GetLogger(ctx).WithField("class", "AbstractPostRepository").WithField("method", "SavePreview")

	id := post.PostID
	logger.Infof("Inserting new preview post for post %v", id)
	var abstractPostID uuid.UUID
	err := txn.QueryRowContext(ctx, SavePreviewPost, post.Title, post.Tagline, post.PreviewImage, post.ViewTime, post.PostID).Scan(&abstractPostID)

	if err != nil {
		logger.Errorf("Error occurred while inserting new preview post for post id %v .%v", id, err)
		return abstractPostID, err
	}

	logger.Infof("Successfully saved preview post for post id %v", id)

	return abstractPostID, nil
}

func NewAbstractPostRepository(db *sqlx.DB) AbstractPostRepository {
	return abstractPostRepository{db: db}
}
