package service

// mockgen -source=service/draft_service.go -destination=mocks/mock_draft_service.go -package=mocks
import (
	"context"
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
	return service.draftRepository.SaveDraft(postData, ctx)
}

func NewDraftService(repository repository.DraftRepository) DraftService {
	return draftService{
		draftRepository: repository,
	}
}
