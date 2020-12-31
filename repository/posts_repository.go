package repository

//go:generate mockgen -source=posts_repository.go -destination=./../mocks/mock_posts_repository.go -package=mocks

import (
	"context"
	"errors"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"post-api/models/db"

	"github.com/gola-glitch/gola-utils/logging"
)

type PostsRepository interface {
	CreatePost(ctx context.Context, post db.PublishPost) error
	LikePost(postID string, userID string, ctx context.Context) error
	UnlikePost(ctx context.Context, userId string, postId string) error
	IsPostLikedByPerson(ctx context.Context, userId string, postId string) (bool, error)
	CommentPost(ctx context.Context, userId string, comment string, postId string) error
	GetLikesCountByPostID(ctx context.Context, postId string) (int64, error)
}

type postRepository struct {
	db neo4j.Session
}

const (
	PublishPost           = "MATCH (author:Person) WHERE author.userId = $userId MATCH (interest:Interest) WHERE interest.name IN $interests MERGE (post:Post{title: $title, puid: $puid, postData: $postData, tagline: $tagline, previewImage: $previewImage, readTime: $readTime})-[audit:PUBLISHED_BY{createdAt: timestamp()}]->(author) MERGE (post)-[:FALLS_UNDER]->(interest)"
	IsPersonLikedThePost  = "MATCH (user:Person{ userId: $userId}) MATCH (post:Post{ puid: $puid}) return EXISTS((user)-[:LIKED]->(post)) as isLiked"
	CommentPost           = "MATCH (user:Person{ userId: $userId}) MATCH (post:Post) WHERE post.puid = $puid MERGE (user)-[:COMMENTED{commentId: apoc.create.uuid(), comment: $comment, createdAt: timestamp()}]->(post)"
	LikePost              = "MATCH (user:Person{ userId: $userId}) MATCH (post:Post) WHERE post.puid = $puid MERGE (user)-[:LIKED]->(post)"
	UnlikePost            = "MATCH (user:Person{ userId: $userId})-[like:LIKED]->(post: Post{ puid: $puid}) delete like"
	GetLikeCountForPostID = "MATCH (readers:Person)-[likes:LIKED]->(post:Post{ puid: $puid}) RETURN count(likes) as likeCount"
)

func (repository postRepository) CreatePost(ctx context.Context, post db.PublishPost) error {
	logger := logging.GetLogger(ctx).WithField("class", "PostsRepository").WithField("method", "CreatePost")

	logger.Infof("Publishing the post for postID %v", post.PUID)

	result, err := repository.db.Run(PublishPost, map[string]interface{}{
		"userId":       post.UserID,
		"postData":     post.PostData.String(),
		"puid":         post.PUID,
		"readTime":     post.ReadTime,
		"interests":    post.Interest,
		"title":        post.Title,
		"tagline":      post.Tagline,
		"previewImage": post.PreviewImage,
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

func NewPostsRepository(db neo4j.Session) PostsRepository {
	return postRepository{db: db}
}
