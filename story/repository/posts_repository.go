package repository

//go:generate mockgen -source=posts_repository.go -destination=./../mocks/mock_posts_repository.go -package=mocks

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"post-api/helper"
	"post-api/story/models/db"
	"strings"

	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, tx helper.Transaction, post db.PublishPost) (uuid.UUID, error)
	Like(ctx context.Context, postID, userID uuid.UUID) error
	UnLike(ctx context.Context, postID, userID uuid.UUID) error
	AddInterests(ctx context.Context, transaction helper.Transaction, postID uuid.UUID, interests []uuid.UUID) error
}

type postRepository struct {
	db *sqlx.DB
}

const (
	PublishPost  = "insert into posts (id, data, author_id, draft_id) values (uuid_generate_v4(), $1, $2, $3) returning id"
	LikePost     = "insert into likes(post_id, liked_by)values($1, $2)"
	UnLike       = "delete from likes where post_id = $1 and liked_by = $2"
	AddInterests = "insert into post_x_interests (post_id, interest_id)values %s"
)

func (repository postRepository) CreatePost(ctx context.Context, tx helper.Transaction, post db.PublishPost) (uuid.UUID, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "CreatePost")
	logger.Infof("Publishing the post in posts table for post draft id %v", post.DraftID)

	var postID uuid.UUID
	err := tx.GetContext(ctx, &postID, PublishPost, post.PostData, post.UserID, post.DraftID)
	if err != nil {
		logger.Errorf("Error occurred while publishing user post in posts table %v", err)
		return postID, err
	}

	logger.Infof("Successfully posted the post in posts table for post draftID %v", post.DraftID)
	return postID, nil
}

func (repository postRepository) Like(ctx context.Context, postID, userID uuid.UUID) error {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("class", "PostsRepository").WithField("method", "Like")
	log.Infof("updating the likedby in like table for post %v", postID)
	_, err := repository.db.ExecContext(ctx, LikePost, postID, userID)

	if err != nil {
		log.Errorf("Error occurred while updating liked by in likes table %v", err)
		return err
	}

	return nil
}

func (repository postRepository) UnLike(ctx context.Context, postID, userID uuid.UUID) error {
	log := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "Like")
	log.Infof("updating the likedby in like table for post %v", postID)
	result, err := repository.db.ExecContext(ctx, UnLike, postID, userID)

	if err != nil {
		log.Errorf("Error occurred while updating liked by in likes table %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		log.Errorf("unable to fetch affected row %v", err)
		return err
	}

	if rowsAffected == 0 {
		log.Error("user never liked the post")
		return errors.New("never liked")
	}

	return nil
}

func (repository postRepository) AddInterests(ctx context.Context, transaction helper.Transaction, postID uuid.UUID, interests []uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "AddInterests")

	valueStrings := make([]string, 0, len(interests))
	valueArgs := make([]interface{}, 0, len(interests)*2)
	i := 0
	for _, id := range interests {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, postID)
		valueArgs = append(valueArgs, id)
		i++
	}
	stmt := fmt.Sprintf(AddInterests, strings.Join(valueStrings, ","))
	logger.Infof("query %v", stmt)
	_, err := transaction.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		logger.Errorf("unable to add interests %v", err)
		return err
	}

	return nil
}

func NewPostsRepository(db *sqlx.DB) PostsRepository {
	return postRepository{db: db}
}
