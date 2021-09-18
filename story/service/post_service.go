package service

//go:generate mockgen -source=post_service.go -destination=./../mocks/mock_post_service.go -package=mocks

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"post-api/helper"
	"post-api/story/constants"
	"post-api/story/models/db"
	"post-api/story/repository"
	"post-api/story/utils"

	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
)

type PostService interface {
	PublishPost(ctx context.Context, draftUID, userUUID uuid.UUID) *golaerror.Error
	LikePost(ctx context.Context, postID, userID uuid.UUID) *golaerror.Error
}

type postService struct {
	transactionManager     helper.TransactionManager
	repository             repository.PostsRepository
	interestRepository     repository.InterestsRepository
	draftRepository        repository.DraftRepository
	abstractPostRepository repository.AbstractPostRepository
	validator              utils.PostValidator
}

func (service postService) PublishPost(ctx context.Context, draftUID, userUUID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "PublishPost")
	draft, err := service.draftRepository.GetDraft(ctx, draftUID, userUUID)
	if err != nil {
		logger.Errorf("error occurred while fetching draft from draft repository %v", err)
		if err == sql.ErrNoRows {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", draftUID, err)
			return &constants.NoDraftFoundError
		}
		return constants.StoryInternalServerError(err.Error())
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
		return apiErr
	}

	logger.Infof("validating draft and generated read time for post %v", draftUID)
	metaData, validationErr := service.validator.ValidateAndGetReadTime(draft, ctx)

	if validationErr != nil {
		logger.Errorf("Error occurred while validating draft of id %v .%v", draftUID, validationErr)
		return validationErr
	}

	if *draft.Tagline == "" {
		draft.Tagline = &metaData.Tagline
	}

	if draft.PreviewImage == nil || *draft.PreviewImage == "" {
		draft.PreviewImage = &metaData.PreviewImage
	}

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
		return constants.StoryInternalServerError(err.Error())
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
		return constants.StoryInternalServerError(err.Error())
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
		return constants.StoryInternalServerError(err.Error())
	}
	_ = txn.Commit()

	logger.Infof("Successfully stored the preview post in preview post repository")
	return nil
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
