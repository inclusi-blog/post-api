package service

//go:generate mockgen -source=post_service.go -destination=./../mocks/mock_post_service.go -package=mocks

import (
	"context"
	"database/sql"
	"post-api/constants"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/repository"
	"post-api/utils"

	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
)

type PostService interface {
	PublishPost(ctx context.Context, draftUID string) *golaerror.Error
	UpdateLikedBy(userID, postID string, ctx context.Context) (request.LikedByCount, *golaerror.Error)
}

type postService struct {
	repository             repository.PostsRepository
	draftRepository        repository.DraftRepository
	previewPostsRepository repository.PreviewPostsRepository
	validator              utils.PostValidator
}

func (service postService) PublishPost(ctx context.Context, draftUID string) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "PublishPost")
	dbDraft, err := service.draftRepository.GetDraft(ctx, draftUID)

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from draft repository %v", err)
		if err == sql.ErrNoRows {
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
		PreviewImage: &dbDraft.PreviewImage.String,
		Tagline:      &dbDraft.Tagline.String,
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
		PUID:      draftUID,
		UserID:    "1",
		PostData:  draft.PostData,
		ReadTime:  metaData.ReadTime,
		ViewCount: 0,
	}

	logger.Infof("Saving post in post repository for post id %v", draftUID)

	postID, err := service.repository.CreatePost(ctx, post)

	if err != nil {
		logger.Errorf("Error occurred while publishing post in post repository %v", err)
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully saved story for post id %v", draftUID)

	previewPost := db.PreviewPost{
		PostID:       postID,
		Title:        metaData.Title,
		Tagline:      *draft.Tagline,
		PreviewImage: *draft.PreviewImage,
		LikeCount:    0,
		CommentCount: 0,
		ViewTime:     0,
	}

	_, err = service.previewPostsRepository.SavePreview(ctx, previewPost)

	if err != nil {
		logger.Errorf("Error occurred while saving preview post for post id %v .%v", post.PUID, err)
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully stored the preview post in preview post repository")

	return nil
}

func (service postService) UpdateLikedBy(userID, postID string, ctx context.Context) (request.LikedByCount, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "UpdateLikedBy")
	logger.Infof("Saving post data to draft repository")

	var likedByCount request.LikedByCount

	// Native Way of Liked By Implementation

	// id, err := service.repository.GetLikesIdByPost(ctx, postID, userID)
	// if err != nil {
	// 	logger.Errorf("Error occurred while Getting like id in likes repository %v", err)
	// 	return likedByCount, constants.StoryInternalServerError(err.Error())
	// }
	// if id == "" {
	// 	err := service.repository.SaveUserToLikedBy(postID, userID, ctx)
	// 	if err != nil {
	// 		logger.Errorf("Error occurred while Getting like id in likes repository %v", err)
	// 		return likedByCount, constants.StoryInternalServerError(err.Error())
	// 	}
	// } else {
	// 	err := service.repository.RemoveUserFromLikedBy(postID, userID, ctx)
	// 	if err != nil {
	// 		logger.Errorf("Error occurred while Getting like id in likes repository %v", err)
	// 		return likedByCount, constants.StoryInternalServerError(err.Error())
	// 	}
	// }

	err := service.repository.AppendOrRemoveUserFromLikedBy(postID, userID, ctx)
	if err != nil {
		logger.Errorf("Error occurred while Updating likedby in likes repository %v", err)
		return likedByCount, constants.StoryInternalServerError(err.Error())
	}

	likeCount, err := service.repository.GetLikeCountByPost(ctx, postID)
	if err != nil {
		logger.Errorf("Error occurred while Getting like id in likes repository %v", err)
		return likedByCount, constants.StoryInternalServerError(err.Error())
	}

	likedByCount.LikeCount = likeCount

	return likedByCount, nil
}

func NewPostService(postsRepository repository.PostsRepository, draftRepository repository.DraftRepository, validator utils.PostValidator, previewPostsRepository repository.PreviewPostsRepository) PostService {
	return postService{
		repository:             postsRepository,
		draftRepository:        draftRepository,
		previewPostsRepository: previewPostsRepository,
		validator:              validator,
	}
}
