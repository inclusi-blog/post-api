package repository

//go:generate mockgen -source=posts_repository.go -destination=./../mocks/mock_posts_repository.go -package=mocks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx/types"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"post-api/constants"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/response"
	"post-api/utils"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, post db.PublishPost, transaction neo4j.Transaction) error
	LikePost(postID string, userID string, ctx context.Context) error
	UnlikePost(ctx context.Context, userId string, postId string) error
	IsPostLikedByPerson(ctx context.Context, userId string, postId string) (bool, error)
	CommentPost(ctx context.Context, userId string, comment string, postId string) error
	GetLikesCountByPostID(ctx context.Context, postId string) (int64, error)
	FetchPost(ctx context.Context, postId string, userId string) (response.Post, error)
}

type postRepository struct {
	db neo4j.Session
}

const (
	PublishPost           = "MATCH (author:Person) WHERE author.userId = $userId MATCH (interest:Interest) WHERE interest.name IN $interests MERGE (post:Post{title: $title, puid: $puid, postData: $postData, tagline: $tagline, previewImage: $previewImage, readTime: $readTime, url: $url})-[audit:PUBLISHED_BY{createdAt: timestamp()}]->(author) MERGE (post)-[:FALLS_UNDER]->(interest)"
	IsPersonLikedThePost  = "MATCH (user:Person{ userId: $userId}) MATCH (post:Post{ puid: $puid}) return EXISTS((user)-[:LIKED]->(post)) as isLiked"
	CommentPost           = "MATCH (user:Person{ userId: $userId}) MATCH (post:Post) WHERE post.puid = $puid MERGE (user)-[:COMMENTED{commentId: apoc.create.uuid(), comment: $comment, createdAt: timestamp()}]->(post)"
	LikePost              = "MATCH (user:Person{ userId: $userId}) MATCH (post:Post) WHERE post.puid = $puid MERGE (user)-[:LIKED]->(post)"
	UnlikePost            = "MATCH (user:Person{ userId: $userId})-[like:LIKED]->(post: Post{ puid: $puid}) delete like"
	GetLikeCountForPostID = "MATCH (readers:Person)-[likes:LIKED]->(post:Post{ puid: $puid}) RETURN count(likes) as likeCount"
	FetchPost             = "MATCH (interests:Interest)<-[tag:FALLS_UNDER]-(post:Post)-[audit:PUBLISHED_BY]->(author:Person) WHERE post.puid = $postId MATCH (user:Person{userId: $userId}) RETURN author.userId AS authorID, author.displayName AS authorName, post.puid as postId, COLLECT(interests.name) AS interests, post.postData AS data, post.previewImage AS previewImage, audit.createdAt AS publishedAt, size((:Person)-[:LIKED]->(post)) AS likeCount, size((:Person)-[:COMMENTED]->(post)) AS commentCount, EXISTS((user)-[:LIKED]->(post)) AS isViewerLiked, CASE WHEN $userId =~ author.userId THEN true ELSE false END AS isAuthorViewing"
)

func (repository postRepository) CreatePost(ctx context.Context, post db.PublishPost, transaction neo4j.Transaction) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "CreatePost")

	logger.Infof("Publishing the post for postID %v", post.PUID)

	result, err := transaction.Run(PublishPost, map[string]interface{}{
		"userId":       post.UserID,
		"postData":     post.PostData.String(),
		"puid":         post.PUID,
		"readTime":     post.ReadTime,
		"interests":    post.Interest,
		"title":        post.Title,
		"tagline":      post.Tagline,
		"previewImage": post.PreviewImage,
		"url":          post.PostUrl,
	})

	if err != nil {
		logger.Errorf("Error occurred while publishing user post in posts table %v", err)
		return err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error while getting result summary for publish post %v", err)
		return err
	}

	logger.Infof("Successfully posted the post for post postID %v", post.PUID)

	return nil
}

func (repository postRepository) LikePost(postID string, userID string, ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("class", "PostsRepository").WithField("method", "LikePost")

	log.Infof("updating the likedby in like table for post %v", postID)

	result, err := repository.db.Run(LikePost, map[string]interface{}{
		"userId": userID,
		"puid":   postID,
	})

	if err != nil {
		log.Errorf("Error occurred while updating liked by in likes table %v", err)
		return err
	}

	_, err = result.Summary()

	if err != nil {
		log.Errorf("Error occurred while getting summary for post like update for post %v ,Error %v", postID, err)
		return err
	}

	return nil
}

func (repository postRepository) IsPostLikedByPerson(ctx context.Context, userId string, postId string) (bool, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "IsPostLikedByPerson")

	logger.Infof("Fetching user like status for post %v", postId)

	result, err := repository.db.Run(IsPersonLikedThePost, map[string]interface{}{
		"userId": userId,
		"puid":   postId,
	})

	if err != nil {
		logger.Errorf("Error occurred while fetching like status for user %v and post %v, Error %v", userId, postId, err)
		return false, err
	}

	if result.Next() {
		isLiked, isPresent := result.Record().Get("isLiked")

		if !isPresent {
			logger.Error("unable to get key value of isLiked from like status query")
			return false, errors.New("unable to find key isLiked")
		}

		isPostLiked := isLiked.(bool)
		return isPostLiked, nil
	}

	return false, errors.New("no results found")
}

func (repository postRepository) UnlikePost(ctx context.Context, userId string, postId string) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "IsPostLikedByPerson")

	logger.Infof("unlike post of user id %v to post id %v", userId, postId)

	result, err := repository.db.Run(UnlikePost, map[string]interface{}{
		"userId": userId,
		"puid":   postId,
	})

	if err != nil {
		logger.Errorf("Error occurred while unliking the post %v by user %v", postId, userId)
		return err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error occurred while fetching result summary for unliking the post %v by user %v", postId, userId)
		return err
	}

	logger.Infof("Successfully disliked post %v by user %v", postId, userId)

	return nil
}

func (repository postRepository) CommentPost(ctx context.Context, userId string, comment string, postId string) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "CommentPost")

	logger.Infof("Commenting on post %v by user %v", postId, userId)

	result, err := repository.db.Run(CommentPost, map[string]interface{}{
		"puid":    postId,
		"userId":  userId,
		"comment": comment,
	})

	if err != nil {
		logger.Errorf("Error occurred while writing comment to post %v by user %v ,Error %v", postId, userId, err)
		return err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error occurred while getting summary for comment to post %v by user %v ,Error %v", postId, userId, err)
		return err
	}

	logger.Infof("Successfully wrote comment to post %v by user %v", postId, userId)

	return nil
}

func (repository postRepository) GetLikesCountByPostID(ctx context.Context, postId string) (int64, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "CommentPost")

	logger.Infof("Fetching likes count for post %v", postId)

	result, err := repository.db.Run(GetLikeCountForPostID, map[string]interface{}{
		"puid": postId,
	})

	if err != nil {
		logger.Errorf("Error occurred while fetching like count for post %v, Error %v", postId, err)
		return 0, err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error occurred while fetching summary for get likes count for post %v, Error %v", postId, err)
		return 0, err
	}

	if result.Next() {
		likeCount, isPresent := result.Record().Get("likeCount")
		if !isPresent {
			logger.Errorf("No like count key present either post not found for post id %v", postId)
			return 0, errors.New("PostNotFound")
		}
		totalLikes := likeCount.(int64)

		logger.Infof("Successfully fetched likes count for post id %v", postId)
		return totalLikes, nil
	}

	return 0, errors.New("no row found")
}

func (repository postRepository) FetchPost(ctx context.Context, postId string, userId string) (response.Post, error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "FetchPost")

	logger.Infof("fetching post to view for user %v of post id %v", userId, postId)

	result, err := repository.db.Run(FetchPost, map[string]interface{}{
		"postId": postId,
		"userId": userId,
	})

	if err != nil {
		logger.Errorf("Error occurred while fetching post to view for user %v of post id %v, Error %v", userId, postId, err)
		return response.Post{}, err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error occurred while fetching summary of post fetch for user id %v of post id %v ,Error %v", userId, postId, err)
		return response.Post{}, err
	}
	if result.Next() {
		var post response.DBPost
		bindDbValues, err := utils.BindDbValues(result, post)
		if err != nil {
			logger.Errorf("binding error %v", err)
			return response.Post{}, err
		}
		jsonString, _ := json.Marshal(bindDbValues)
		err = json.Unmarshal(jsonString, &post)
		return mapDBPostToPost(post), nil
	}

	return response.Post{}, errors.New(constants.NoPostFound)
}

func mapDBPostToPost(post response.DBPost) response.Post {
	return response.Post{
		PostID: post.PostID,
		PostData: models.JSONString{
			JSONText: types.JSONText(post.PostData),
		},
		LikeCount:              post.LikeCount,
		CommentCount:           post.CommentCount,
		Interests:              post.Interests,
		AuthorID:               post.AuthorID,
		AuthorName:             post.AuthorName,
		PreviewImage:           post.PreviewImage,
		PublishedAt:            post.PublishedAt,
		IsViewerLiked:          post.IsViewerLiked,
		IsViewerIsAuthor:       post.IsViewerIsAuthor,
		IsViewerFollowedAuthor: post.IsViewerFollowedAuthor,
	}
}

func NewPostsRepository(db neo4j.Session) PostsRepository {
	return postRepository{db: db}
}
