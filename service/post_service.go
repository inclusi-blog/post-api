package service

//go:generate mockgen -source=post_service.go -destination=./../mocks/mock_post_service.go -package=mocks

import (
	"context"
	"database/sql"
	"fmt"
	"post-api/constants"
	"post-api/models/db"
	"post-api/repository"
	"post-api/utils"

	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
)

type PostService interface {
	PublishPost(ctx context.Context, draftUID string) *golaerror.Error
}

type postService struct {
	repository      repository.PostsRepository
	draftRepository repository.DraftRepository
	validator       utils.PostValidator
}

func (service postService) PublishPost(ctx context.Context, draftUID string) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "PostService").WithField("method", "PublishPost")
	draft, err := service.draftRepository.GetDraft(ctx, draftUID)

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from draft repository %v", err)
		if err == sql.ErrNoRows {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", draftUID, err)
			return &constants.NoDraftFoundError
		}
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Validating draft and Generated read time for post %v", draftUID)

	readTime, err := service.validator.ValidateAndGetReadTime(&draft, ctx)
	fmt.Println(draft.Tagline)
	if err != nil {
		failedError := constants.DraftValidationFailedError
		failedError.AdditionalData = err.Error()
		logger.Errorf("Error occurred while validating draft of id %v .%v", draftUID, err)
		return &failedError
	}

	post := db.PublishPost{
		PUID:      draftUID,
		UserID:    "1",
		PostData:  draft.PostData,
		TitleData: draft.TitleData,
		ReadTime:  readTime,
		ViewCount: 0,
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

func NewPostService(postsRepository repository.PostsRepository, draftRepository repository.DraftRepository, validator utils.PostValidator) PostService {
	return postService{
		repository:      postsRepository,
		draftRepository: draftRepository,
		validator:       validator,
	}
}
