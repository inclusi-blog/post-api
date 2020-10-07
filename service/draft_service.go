package service

// mockgen -source=service/draft_service.go -destination=mocks/mock_draft_service.go -package=mocks
import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/constants"
	"post-api/models"
	"post-api/repository"
)

type DraftService interface {
	SaveDraft(postData models.UpsertDraft, ctx context.Context) error
}

type draftService struct {
	draftRepository repository.DraftRepository
}

func (service draftService) SaveDraft(postData models.UpsertDraft, ctx context.Context) error {
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
