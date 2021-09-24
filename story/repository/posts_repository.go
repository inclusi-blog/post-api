package repository

//go:generate mockgen -source=posts_repository.go -destination=./../mocks/mock_posts_repository.go -package=mocks

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"post-api/helper"
	"post-api/story/models/db"
	"post-api/story/models/request"
	"post-api/story/models/response"
	"strings"

	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, tx helper.Transaction, post db.PublishPost) (uuid.UUID, error)
	Like(ctx context.Context, postID, userID uuid.UUID) error
	UnLike(ctx context.Context, postID, userID uuid.UUID) error
	AddInterests(ctx context.Context, transaction helper.Transaction, postID uuid.UUID, interests []uuid.UUID) error
	FetchPost(ctx context.Context, postId, userId uuid.UUID) (response.Post, error)
	GetPublishedPostByUser(ctx context.Context, request request.GetPublishedPostRequest) ([]response.PublishedPost, error)
}

type postRepository struct {
	db *sqlx.DB
}

const (
	PublishPost       = "insert into posts (id, data, author_id, draft_id) values (uuid_generate_v4(), $1, $2, $3) returning id"
	LikePost          = "insert into likes(post_id, liked_by)values($1, $2)"
	UnLike            = "delete from likes where post_id = $1 and liked_by = $2"
	AddInterests      = "insert into post_x_interests (post_id, interest_id)values %s"
	GetPost           = "select posts.id, posts.data, count(l) as likes_count, count(c) as comments_count, json_agg(jsonb_build_object('id', interests.id, 'name', interests.name)) as interests, u.id as author_id, u.username as author_name, ap.preview_image as preview_image, posts.created_at as published_at, case when $1 in (l.post_id) then true else false end as is_viewer_liked, case when $2 = u.id then true else false end as is_viewer_is_author from posts inner join users u on u.id = posts.author_id inner join post_x_interests on posts.id = post_x_interests.post_id inner join interests on post_x_interests.interest_id = interests.id inner join abstract_post ap on posts.id = ap.post_id left join comments c on posts.id = c.post_id left join likes l on posts.id = l.post_id where posts.id = $3 group by posts.id, posts.data, u.id, u.username, preview_image, posts.created_at, is_viewer_liked, is_viewer_is_author"
	GetPublishedPosts = "select posts.id, ap.title, ap.tagline, posts.created_at, json_agg(jsonb_build_object('id', i.id, 'name', i.name)) as interests, count(l) as likes_count, preview_image from posts inner join users on posts.author_id = users.id inner join abstract_post ap on posts.id = ap.post_id inner join post_x_interests pxi on posts.id = pxi.post_id inner join interests i on pxi.interest_id = i.id left join likes l on posts.id = l.post_id where users.id = $1 group by posts.id, posts.created_at, ap.title, ap.tagline, posts.id, preview_image order by posts.created_at limit $2 offset $3"
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

func (repository postRepository) FetchPost(ctx context.Context, postId, userId uuid.UUID) (response.Post, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "FetchPost")

	logger.Infof("fetching post to view for user %v of post id %v", userId, postId)
	var post response.Post
	err := repository.db.GetContext(ctx, &post, GetPost, userId, userId, postId)

	if err != nil {
		logger.Errorf("Error occurred while fetching post to view for user %v of post id %v, Error %v", userId, postId, err)
		return response.Post{}, err
	}

	return post, nil
}

func (repository postRepository) GetPublishedPostByUser(ctx context.Context, request request.GetPublishedPostRequest) ([]response.PublishedPost, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "GetPublishedPostByUser")
	logger.Infof("fetching user published post for user %v", request.UserID)

	var posts []response.PublishedPost
	err := repository.db.SelectContext(ctx, &posts, GetPublishedPosts, request.UserID, request.Limit, request.StartValue)
	if err != nil {
		logger.Errorf("unable to get published posts %v", err)
		return nil, err
	}

	return posts, nil
}

func NewPostsRepository(db *sqlx.DB) PostsRepository {
	return postRepository{db: db}
}
