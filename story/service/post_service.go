package service

//go:generate mockgen -source=post_service.go -destination=./../mocks/mock_post_service.go -package=mocks

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"post-api/helper"
	"post-api/story/constants"
	"post-api/story/models/db"
	"post-api/story/models/request"
	"post-api/story/models/response"
	"post-api/story/repository"
	"post-api/story/utils"
	"strings"

	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
)

type PostService interface {
	GetPost(ctx context.Context, postId, userId uuid.UUID) (response.Post, *golaerror.Error)
	PublishPost(ctx context.Context, draftUID, userUUID uuid.UUID) (string, *golaerror.Error)
	LikePost(ctx context.Context, postID, userID uuid.UUID) *golaerror.Error
	UnLikePost(ctx context.Context, postUID, userID uuid.UUID) *golaerror.Error
	GetPublishedPostByUser(ctx context.Context, request request.GetPublishedPostRequest) ([]response.PublishedPost, *golaerror.Error)
	Comment(ctx context.Context, comment request.Comment) *golaerror.Error
	GetComments(ctx context.Context, commentsRequest request.FetchComments) ([]response.Comment, *golaerror.Error)
	SavePost(ctx context.Context, postID, userID uuid.UUID) *golaerror.Error
	MarkAsViewed(ctx context.Context, postID, userID uuid.UUID) *golaerror.Error
	FetchSavedPosts(ctx context.Context, postRequest request.PostRequest) ([]response.PostView, *golaerror.Error)
	FetchViewedPosts(ctx context.Context, postRequest request.PostRequest) ([]response.PostView, *golaerror.Error)
}

type postService struct {
	transactionManager     helper.TransactionManager
	repository             repository.PostsRepository
	interestRepository     repository.InterestsRepository
	draftRepository        repository.DraftRepository
	abstractPostRepository repository.AbstractPostRepository
	validator              utils.PostValidator
}

func (service postService) PublishPost(ctx context.Context, draftUID, userUUID uuid.UUID) (string, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "PublishPost")
	draft, err := service.draftRepository.GetDraft(ctx, draftUID, userUUID)
	if err != nil {
		logger.Errorf("error occurred while fetching draft from draft repository %v", err)
		if err == sql.ErrNoRows {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", draftUID, err)
			return "", &constants.NoDraftFoundError
		}
		return "", constants.StoryInternalServerError(err.Error())
	}
	apiErr := draft.ConvertInterests(func(interests []string) *golaerror.Error {
		draft.InterestTags, err = service.interestRepository.GetInterestsForName(ctx, interests)
		if err != nil {
			logger.Errorf("unable to get interests %v", err)
			return constants.StoryInternalServerError("something went wrong")
		}
		return nil
	})

	if apiErr != nil {
		logger.Error("unable to get interests")
		return "", apiErr
	}

	logger.Infof("validating draft and generated read time for post %v", draftUID)
	metaData, validationErr := service.validator.ValidateAndGetReadTime(draft, ctx)

	if validationErr != nil {
		logger.Errorf("Error occurred while validating draft of id %v .%v", draftUID, validationErr)
		return "", validationErr
	}

	if *draft.Tagline == "" {
		draft.Tagline = &metaData.Tagline
	}

	if draft.PreviewImage == nil || *draft.PreviewImage == "" {
		draft.PreviewImage = &metaData.PreviewImage
	}

	url := utils.GenerateUrl(metaData.Title)

	post := db.PublishPost{
		DraftID:  draftUID,
		UserID:   userUUID,
		PostData: draft.Data,
	}
	txn := service.transactionManager.NewTransaction()
	logger.Infof("Saving post in post repository for post id %v", draftUID)
	postID, err := service.repository.CreatePost(ctx, txn, post)

	if err != nil {
		_ = txn.Rollback()
		logger.Errorf("Error occurred while publishing post in post repository %v", err)
		return "", constants.StoryInternalServerError(err.Error())
	}
	logger.Infof("Successfully saved story for post id %v", draftUID)
	var interests []uuid.UUID
	for _, interest := range draft.InterestTags {
		interests = append(interests, interest.ID)
	}
	err = service.repository.AddInterests(ctx, txn, postID, interests)
	if err != nil {
		_ = txn.Rollback()
		logger.Errorf("unable to add interests to posts %v", err)
		return "", constants.StoryInternalServerError(err.Error())
	}
	abstractPost := db.AbstractPost{
		PostID:       postID,
		Title:        metaData.Title,
		Tagline:      *draft.Tagline,
		PreviewImage: *draft.PreviewImage,
		ViewTime:     int64(metaData.ReadTime),
	}

	_, err = service.abstractPostRepository.Save(ctx, txn, abstractPost)
	if err != nil {
		_ = txn.Rollback()
		logger.Errorf("Error occurred while saving abstract post for post id %v .%v", abstractPost, err)
		return "", constants.StoryInternalServerError(err.Error())
	}

	err = service.draftRepository.UpdatePublishStatus(ctx, txn, draftUID, userUUID, true)
	if err != nil {
		logger.Errorf("unable to update the publish status %", err)
		_ = txn.Rollback()
		return "", constants.StoryInternalServerError(err.Error())
	}

	_ = txn.Commit()

	finalPostUrl := strings.Join([]string{url, postID.String()}, "-")

	logger.Infof("Successfully stored the preview post in preview post repository")
	return finalPostUrl, nil
}

func (service postService) LikePost(ctx context.Context, postUID, userID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "LikePost")
	logger.Infof("Saving post data to draft repository")

	err := service.repository.Like(ctx, postUID, userID)
	if err != nil {
		logger.Errorf("Error occurred while Updating likedby in likes repository %v", err)
		return constants.StoryInternalServerError(err.Error())
	}
	return nil
}

func (service postService) UnLikePost(ctx context.Context, postUID, userID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "LikePost")
	logger.Infof("Saving post data to draft repository")

	err := service.repository.UnLike(ctx, postUID, userID)
	if err != nil {
		logger.Errorf("Error occurred while Updating likedby in likes repository %v", err)
		return constants.StoryInternalServerError(err.Error())
	}
	return nil
}

func (service postService) GetPost(ctx context.Context, postId, userId uuid.UUID) (response.Post, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "GetPost")
	logger.Infof("Fetching post for post id %v", postId)

	post, err := service.repository.FetchPost(ctx, postId, userId)

	if err != nil {
		logger.Errorf("Error occurred while fetching post for given post id %v, Error %v", postId, err)
		if err == sql.ErrNoRows {
			logger.Errorf("Error No post found for given post id %v", postId)
			return response.Post{}, &constants.PostNotFoundErr
		}
		return response.Post{}, constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully fetching post from post repository for given post id %v", postId)

	return post, nil
}

func (service postService) GetPublishedPostByUser(ctx context.Context, request request.GetPublishedPostRequest) ([]response.PublishedPost, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "GetPost")
	logger.Infof("Fetching posts for user id %v", request.UserID)

	posts, err := service.repository.GetPublishedPostByUser(ctx, request)

	if err != nil {
		logger.Errorf("Error occurred while fetching posts for given user id %v, Error %v", request.UserID, err)
		return nil, constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully fetching posts from post repository for given user id %v", request.UserID)

	return posts, nil
}

func (service postService) Comment(ctx context.Context, comment request.Comment) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "Comment")

	err := service.repository.Comment(ctx, comment)
	if err != nil {
		logger.Infof("unable to comment %v", err)
		return constants.StoryInternalServerError(err.Error())
	}
	logger.Info("comment successfully posted")

	return nil
}

func (service postService) GetComments(ctx context.Context, commentsRequest request.FetchComments) ([]response.Comment, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "FetchComments")
	logger.Infof("fetching post for post id %v", commentsRequest.PostID)

	comments, err := service.repository.FetchComments(ctx, commentsRequest)
	if err != nil {
		logger.Errorf("unable to fetch comments from repository %v", err)
		return nil, constants.StoryInternalServerError(err.Error())
	}

	logger.Info("successfully fetched comments")

	return comments, nil
}

func (service postService) SavePost(ctx context.Context, postID, userID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "FetchComments")
	logger.Info("marking post as read later")

	err := service.repository.BookmarkPost(ctx, postID, userID)
	if err != nil {
		logger.Errorf("unable to mark post as read later %v", err)
		return &constants.InternalServerError
	}
	logger.Info("successfully marked post as saved")

	return nil
}

func (service postService) MarkAsViewed(ctx context.Context, postID, userID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "FetchComments")
	logger.Info("marking post as viewed")

	err := service.repository.MarkAsViewed(ctx, postID, userID)
	if err != nil {
		logger.Errorf("unable to mark post as viewed %v", err)
		return &constants.InternalServerError
	}
	logger.Info("successfully marked post as viewed")

	return nil
}

func (service postService) FetchSavedPosts(ctx context.Context, postRequest request.PostRequest) ([]response.PostView, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "FetchSavedPosts")

	posts, err := service.repository.FetchReadLater(ctx, postRequest)
	if err != nil {
		logger.Errorf("unable to fetch read later posts %v", err)
		return nil, &constants.InternalServerError
	}

	return posts, nil
}

func (service postService) FetchViewedPosts(ctx context.Context, postRequest request.PostRequest) ([]response.PostView, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "FetchViewedPosts")

	posts, err := service.repository.FetchViewedPosts(ctx, postRequest)
	if err != nil {
		logger.Errorf("unable to fetch viewed posts %v", err)
		return nil, &constants.InternalServerError
	}

	return posts, nil
}

func NewPostService(postsRepository repository.PostsRepository, draftRepository repository.DraftRepository, validator utils.PostValidator, previewPostsRepository repository.AbstractPostRepository, interestsRepository repository.InterestsRepository, manager helper.TransactionManager) PostService {
	return postService{
		transactionManager:     manager,
		repository:             postsRepository,
		interestRepository:     interestsRepository,
		draftRepository:        draftRepository,
		abstractPostRepository: previewPostsRepository,
		validator:              validator,
	}
}
