package service

//go:generate mockgen -source=post_service.go -destination=./../mocks/mock_post_service.go -package=mocks

import (
	"context"
	"post-api/constants"
	"post-api/models/db"
	"post-api/models/response"
	"post-api/repository"
	"post-api/utils"

	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
)

type PostService interface {
	PublishPost(ctx context.Context, draftUID, userId string) *golaerror.Error
	LikePost(userID string, postID string, ctx context.Context) (response.LikedByCount, *golaerror.Error)
	UnlikePost(userId string, postId string, ctx context.Context) (response.LikedByCount, *golaerror.Error)
	CommentPost(ctx context.Context, userId, postId, comment string) *golaerror.Error
}

type postService struct {
	repository      repository.PostsRepository
	draftRepository repository.DraftRepository
	validator       utils.PostValidator
}

func (service postService) PublishPost(ctx context.Context, draftUID string, userId string) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "PublishPost")
	dbDraft, err := service.draftRepository.GetDraft(ctx, draftUID, userId)

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from draft repository %v", err)
		if err.Error() == constants.NoDraftFoundCode {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", draftUID, err)
			return &constants.NoDraftFoundError
		}
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Validating draft and Generated read time for post %v", draftUID)

	draft := db.Draft{
		DraftID:      dbDraft.DraftID,
		UserID:       dbDraft.UserID,
		PostData:     dbDraft.PostData,
		PreviewImage: &dbDraft.PreviewImage,
		Tagline:      &dbDraft.Tagline,
		Interest:     dbDraft.Interest,
	}

	metaData, validationErr := service.validator.ValidateAndGetReadTime(draft, ctx)

	if validationErr != nil {
		logger.Errorf("Error occurred while validating draft of id %v .%v", draftUID, validationErr)
		return validationErr
	}

	if *draft.Tagline == "" {
		draft.Tagline = &metaData.Tagline
	}

	if *draft.PreviewImage == "" {
		draft.PreviewImage = &metaData.PreviewImage
	}

	post := db.PublishPost{
		PUID:         draftUID,
		UserID:       draft.UserID,
		PostData:     draft.PostData,
		ReadTime:     metaData.ReadTime,
		Interest:     draft.Interest,
		Title:        metaData.Title,
		Tagline:      *draft.Tagline,
		PreviewImage: *draft.PreviewImage,
	}

	logger.Infof("Saving post in post repository for post id %v", draftUID)

	err = service.repository.CreatePost(ctx, post)

	if err != nil {
		logger.Errorf("Error occurred while publishing post in post repository %v", err)
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully saved story for post id %v", draftUID)

	return nil
}

func (service postService) LikePost(userID string, postUID string, ctx context.Context) (response.LikedByCount, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "LikePost")
	logger.Infof("Saving post data to draft repository")

	var likedByCount response.LikedByCount

	err := service.repository.LikePost(postUID, userID, ctx)
	if err != nil {
		logger.Errorf("Error occurred while Updating likedby in likes repository %v", err)
		return likedByCount, constants.StoryInternalServerError(err.Error())
	}

	likeCount, err := service.repository.GetLikesCountByPostID(ctx, postUID)
	if err != nil {
		logger.Errorf("Error occurred while Getting like id in likes repository %v", err)
		return likedByCount, constants.StoryInternalServerError(err.Error())
	}

	likedByCount.LikeCount = likeCount

	return likedByCount, nil
}

func (service postService) UnlikePost(userId string, postId string, ctx context.Context) (response.LikedByCount, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "LikePost")
	logger.Infof("Saving post data to draft repository")

	var likedByCount response.LikedByCount

	err := service.repository.UnlikePost(ctx, userId, postId)
	if err != nil {
		logger.Errorf("Error occurred while Updating likedby in likes repository %v", err)
		return likedByCount, constants.StoryInternalServerError(err.Error())
	}

	likeCount, err := service.repository.GetLikesCountByPostID(ctx, postId)
	if err != nil {
		logger.Errorf("Error occurred while Getting like id in likes repository %v", err)
		return likedByCount, constants.StoryInternalServerError(err.Error())
	}

	likedByCount.LikeCount = likeCount

	return likedByCount, nil
}

func (service postService) CommentPost(ctx context.Context, userId, postId, comment string) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "CommentPost")
	logger.Infof("Commenting on post %v by user %v", postId, userId)

	err := service.repository.CommentPost(ctx, userId, comment, postId)

	if err != nil {
		logger.Errorf("Error occurred while user %v commenting on post %v", userId, postId)
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully commented on post %v by user %v", postId, userId)
	return nil
}

func NewPostService(postsRepository repository.PostsRepository, draftRepository repository.DraftRepository, validator utils.PostValidator) PostService {
	return postService{
		repository:      postsRepository,
		draftRepository: draftRepository,
		validator:       validator,
	}
}
