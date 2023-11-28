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

	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/jmoiron/sqlx"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, tx helper.Transaction, post db.PublishPost) (uuid.UUID, error)
	Like(ctx context.Context, postID, userID uuid.UUID) error
	UnLike(ctx context.Context, postID, userID uuid.UUID) error
	AddInterests(ctx context.Context, transaction helper.Transaction, postID uuid.UUID, interests []uuid.UUID) error
	FetchPost(ctx context.Context, postId, userId uuid.UUID) (response.Post, error)
	GetPublishedPostByUser(ctx context.Context, request request.GetPublishedPostRequest) ([]response.PublishedPost, error)
	Comment(ctx context.Context, comment request.Comment) error
	FetchComments(ctx context.Context, commentsRequest request.FetchComments) ([]response.Comment, error)
	BookmarkPost(ctx context.Context, postID, userID uuid.UUID) error
	MarkAsViewed(ctx context.Context, postID, userID uuid.UUID) error
	FetchReadLater(ctx context.Context, postRequest request.PostRequest) ([]response.PostView, error)
	FetchViewedPosts(ctx context.Context, postRequest request.PostRequest) ([]response.PostView, error)
	FetchPostsByInterests(ctx context.Context, interestRequest request.InterestRequest, userID uuid.UUID) ([]response.PostView, error)
	RemovePostBookmark(ctx context.Context, postID, userID uuid.UUID) error
	Delete(ctx context.Context, postID, userID uuid.UUID) error
	GetHomeFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]db.HomeFeedPost, error)
}

type postRepository struct {
	db *sqlx.DB
}

const (
	PublishPost        = "insert into posts (id, data, author_id, draft_id) values (uuid_generate_v4(), $1, $2, $3) returning id"
	LikePost           = "insert into likes(post_id, liked_by)values($1, $2)"
	UnLike             = "delete from likes where post_id = $1 and liked_by = $2"
	CommentPost        = "insert into comments (id, data, post_id, commented_by) values (uuid_generate_v4(), $1, $2, $3)"
	AddInterests       = "insert into post_x_interests (post_id, interest_id)values %s"
	GetPost            = "with post_interests as (select jsonb_agg(jsonb_build_object('id', interests.id, 'name', interests.name)) as interests, post_id from posts inner join post_x_interests on posts.id = post_x_interests.post_id inner join interests on post_x_interests.interest_id = interests.id where posts.id = $1 group by post_x_interests.post_id) select posts.id, posts.data, count(distinct l.liked_by) as likes_count, count(distinct c.id) as comments_count, post_interests.interests, u.id as author_id, u.username as author_name, ap.preview_image as preview_image, posts.created_at as published_at, ap.url, case when $2 in (l.post_id) then true else false end as is_viewer_liked, case when $3 = u.id then true else false end as is_viewer_is_author from posts inner join post_interests on posts.id = post_interests.post_id inner join post_x_interests on posts.id = post_x_interests.post_id inner join interests on post_x_interests.interest_id = interests.id inner join users u on u.id = posts.author_id inner join abstract_post ap on posts.id = ap.post_id left join likes l on l.post_id = posts.id left join comments c on c.post_id = posts.id where posts.id = $4 group by posts.id, u.id, ap.preview_image, l.post_id, ap.url, post_interests.interests"
	GetPublishedPosts  = "select posts.id, ap.title, ap.tagline, posts.created_at, (select json_agg(json_build_object('id', interest_id, 'name', i.name)) from post_x_interests inner join interests i on post_x_interests.interest_id = i.id where post_x_interests.post_id = posts.id) as interests, count(l) as likes_count, username, preview_image, ap.url from posts inner join users on posts.author_id = users.id inner join abstract_post ap on posts.id = ap.post_id left join likes l on posts.id = l.post_id where users.id = $1 group by posts.id, posts.created_at, ap.title, ap.tagline, posts.id, ap.url, preview_image, username order by posts.created_at limit $2 offset $3"
	GetComments        = "select comments.id, comments.data, comments.post_id, u.username, comments.created_at from comments inner join users u on u.id = comments.commented_by where post_id = $1 order by comments.created_at desc limit $2 offset $3"
	BookmarkPost       = "insert into saved_posts (post_id, user_id) values ($1, $2)"
	RemovePostBookmark = "delete from saved_posts where post_id = $1 and user_id = $2"
	MarkAsViewed       = "insert into post_views (post_id, user_id) values ($1, $2)"
	Delete             = "update posts set deleted_at = current_timestamp where id = $1 and author_id = $2"
	GetHomeFeed        = `WITH post_interests AS (
    SELECT
        px.post_id,
        ARRAY_AGG(i.name) AS interest_names,
        ARRAY_AGG(i.id) AS interest_ids
    FROM
        post_x_interests px
            JOIN interests i ON px.interest_id = i.id
    GROUP BY
        px.post_id
),
     post_likes AS (
         SELECT
             post_id,
             COUNT(*) AS like_count
         FROM
             likes
         GROUP BY
             post_id
     ),
     post_comments AS (
         SELECT
             p.id AS post_id,
             COUNT(*) AS comment_count
         FROM
             posts p
                 LEFT JOIN comments c ON p.id = c.post_id
         GROUP BY
             p.id
     ),
     post_author AS (
         SELECT
             p.id AS post_id,
             COALESCE(a.name, u.username) AS author_name
         FROM
             posts p
                 LEFT JOIN users u ON p.author_id = u.id
                 LEFT JOIN admin a ON p.author_id = a.id
     ),
     user_likes AS (
         SELECT
             post_id,
             TRUE AS user_liked
         FROM
             likes
         WHERE
                         liked_by = $1
     ),
     post_statistics AS (
         SELECT
             p.post_id AS post_id,
             p.view_time,
             COALESCE(pl.like_count, 0) AS like_count,
             COALESCE(pc.comment_count, 0) AS comment_count,
             p.created_at
         FROM
             abstract_post p
                 LEFT JOIN post_likes pl ON p.post_id = pl.post_id
                 LEFT JOIN post_comments pc ON p.post_id = pc.post_id
     )
SELECT
    ps.post_id,
    ap.title,
    ap.tagline,
    ap.view_time,
    ap.created_at AS published_date,
    pi.interest_names,
    pa.author_name,
    ps.like_count,
    ul.user_liked,
    ap.preview_image
FROM
    post_statistics ps
        JOIN abstract_post ap ON ps.post_id = ap.post_id
        JOIN post_interests pi ON ps.post_id = pi.post_id
        LEFT JOIN post_author pa ON ps.post_id = pa.post_id
        LEFT JOIN user_likes ul ON ps.post_id = ul.post_id
WHERE
    ap.deleted_at IS NULL
  OR (pi.interest_ids @> ARRAY(SELECT interest_id FROM user_interests WHERE user_id = $2))
ORDER BY
    ap.created_at DESC limit $3 offset $4;`
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
	err := repository.db.GetContext(ctx, &post, GetPost, postId, postId, userId, postId)

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

func (repository postRepository) Comment(ctx context.Context, comment request.Comment) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "Comment")
	logger.Infof("inserting comment for post %v by user id %v ", comment.PostID, comment.CommentedBy)

	result, err := repository.db.ExecContext(ctx, CommentPost, comment.Data, comment.PostID, comment.CommentedBy)

	if err != nil {
		logger.Errorf("unable to comment %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("unable to fetch affected rows %v", err)
		return err
	}

	if rowsAffected == 0 {
		logger.Error("user never liked the post")
		return errors.New("unable to comment")
	}

	return nil
}

func (repository postRepository) FetchComments(ctx context.Context, commentsRequest request.FetchComments) ([]response.Comment, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "Comment")
	logger.Infof("fetching comment for post id %v ", commentsRequest.PostID)

	var comments []response.Comment
	err := repository.db.SelectContext(ctx, &comments, GetComments, commentsRequest.PostID, commentsRequest.Limit, commentsRequest.Start)
	if err != nil {
		logger.Errorf("unable to fetch comments %v", err)
		return nil, err
	}

	return comments, nil
}

func (repository postRepository) BookmarkPost(ctx context.Context, postID, userID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "BookmarkPost")

	_, err := repository.db.ExecContext(ctx, BookmarkPost, postID, userID)
	if err != nil {
		logger.Errorf("unable to mark post as read later %v", err)
		return err
	}

	return nil
}

func (repository postRepository) RemovePostBookmark(ctx context.Context, postID, userID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "BookmarkPost")

	_, err := repository.db.ExecContext(ctx, RemovePostBookmark, postID, userID)
	if err != nil {
		logger.Errorf("unable to mark post as read later %v", err)
		return err
	}

	return nil
}

func (repository postRepository) MarkAsViewed(ctx context.Context, postID, userID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "Comment")

	_, err := repository.db.ExecContext(ctx, MarkAsViewed, postID, userID)
	if err != nil {
		logger.Errorf("unable to mark post as viewed %v", err)
		return err
	}

	return nil
}

func (repository postRepository) FetchReadLater(ctx context.Context, postRequest request.PostRequest) ([]response.PostView, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "FetchSavedPosts")
	userID := postRequest.UserID
	var posts []response.PostView
	err := repository.db.SelectContext(ctx, &posts, FetchSavedPosts, userID, userID, userID, userID, postRequest.Limit, postRequest.Start)
	if err != nil {
		logger.Errorf("unable to fetch read later post %v", err)
		return nil, err
	}

	return posts, nil
}

func (repository postRepository) FetchViewedPosts(ctx context.Context, postRequest request.PostRequest) ([]response.PostView, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "FetchViewedPosts")
	userID := postRequest.UserID
	var posts []response.PostView
	err := repository.db.SelectContext(ctx, &posts, FetchViewedPosts, userID, userID, userID, userID, userID, userID,
		postRequest.Limit, postRequest.Start)
	if err != nil {
		logger.Errorf("unable to fetch viewed posts %v", err)
		return nil, err
	}

	return posts, nil
}

func (repository postRepository) FetchPostsByInterests(ctx context.Context, interestRequest request.InterestRequest, userID uuid.UUID) ([]response.PostView, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "FetchPostsByInterests")
	var posts []response.PostView
	err := repository.db.SelectContext(ctx, &posts, FetchPostByInterests, userID, userID, userID, interestRequest.InterestUID, interestRequest.Limit, interestRequest.Start)
	if err != nil {
		logger.Errorf("unable to fetch posts for interest %v", err)
		return nil, err
	}

	return posts, nil
}

func (repository postRepository) Delete(ctx context.Context, postID, userID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostRepository").WithField("method", "Delete")
	logger.Infof("deleting post for post id %v and author id %v", postID, userID)

	_, err := repository.db.ExecContext(ctx, Delete, postID, userID)
	if err != nil {
		logger.Errorf("unable to delete post for post id %v. Error %v", postID, err)
		return err
	}

	logger.Infof("successfully deleted post for post id %v", postID)
	return nil
}

func (repository postRepository) GetHomeFeed(ctx context.Context, userID uuid.UUID, limit, offset int) ([]db.HomeFeedPost, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "GetHomeFeed")
	var posts []db.HomeFeedPost
	err := repository.db.SelectContext(ctx, &posts, GetHomeFeed, userID, userID, limit, offset)
	if err != nil {
		logger.Errorf("unable to fetch posts for interest %v", err)
		return nil, err
	}

	return posts, nil
}

func NewPostsRepository(db *sqlx.DB) PostsRepository {
	return postRepository{db: db}
}
