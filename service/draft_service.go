package service

//go:generate mockgen -source=draft_service.go -destination=./../mocks/mock_draft_service.go -package=mocks

import (
	"context"
	"post-api/constants"
	"post-api/models"
	"post-api/models/request"
	"post-api/repository"

	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
)

type DraftService interface {
	SaveDraft(postData models.UpsertDraft, ctx context.Context) *golaerror.Error
	UpsertTagline(taglineRequest request.TaglineSaveRequest, ctx context.Context) *golaerror.Error
	UpsertInterests(interestRequest request.InterestsSaveRequest, ctx context.Context) *golaerror.Error
}

type draftService struct {
	draftRepository repository.DraftRepository
}

func (service draftService) UpsertInterests(interestRequest request.InterestsSaveRequest, ctx context.Context) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "UpsertInterests")

	logger.Info("Calling service to save interests for draft")

	err := service.draftRepository.SaveInterestsToDraft(interestRequest, ctx)

	if err != nil {
		logger.Errorf("Error occurred while inserting interests in draft repository %v", err)
		return &constants.PostServiceFailureError
	}

	logger.Info("Successfully stored draft interests")

	return nil
}

func (service draftService) SaveDraft(postData models.UpsertDraft, ctx context.Context) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "SaveDraft")
	if postData.Target == "post" {
		logger.Infof("Saving post data to draft repository")
		err := service.draftRepository.SavePostDraft(postData, ctx)
		return InternalServerError(err, logger)
	}
	logger.Infof("Saving title data to draft repository")
	err := service.draftRepository.SaveTitleDraft(postData, ctx)
	return InternalServerError(err, logger)
}

func (service draftService) UpsertTagline(taglineRequest request.TaglineSaveRequest, ctx context.Context) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "UpsertTagline")

	logger.Info("Calling service to save tagline for draft")

	err := service.draftRepository.SaveTaglineToDraft(taglineRequest, ctx)

	if err != nil {
		logger.Errorf("Error occurred while inserting tagline in draft repository %v", err)
		return &constants.PostServiceFailureError
	}

	logger.Info("Successfully stored draft tagline")

	return nil
}

func InternalServerError(err error, logger logging.GolaLoggerEntry) *golaerror.Error {
	if err != nil {
		logger.Errorf("Error occurred while saving draft data into draft repository %v", err)
		return &constants.InternalServerError
	}
	return nil
}

func NewDraftService(repository repository.DraftRepository) DraftService {
	return draftService{
		draftRepository: repository,
	}
}
