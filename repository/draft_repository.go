package repository

// mockgen -source=repository/draft_repository.go -destination=mocks/mock_draft_repository.go -package=mocks
import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
	"post-api/models"
)

type DraftRepository interface {
	SaveDraft(draft models.UpsertDraft, ctx context.Context) error
}

const (
	SAVE_DRAFT = "INSERT INTO DRAFTS (DRAFT_ID, USER_ID, POST_DATA) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE POST_DATA = ?, UPDATED_AT = current_timestamp"
)

type draftRepository struct {
	db *sqlx.DB
}

func (repository draftRepository) SaveDraft(draft models.UpsertDraft, ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("class", "DraftRepository").WithField("method", "SaveDraft")

	log.Infof("Inserting or updating the existing post in draft for user %v", draft.UserID)

	_, err := repository.db.ExecContext(ctx, SAVE_DRAFT, draft.DraftID, draft.UserID, draft.PostData, draft.PostData)

	if err != nil {
		log.Errorf("Error occurred while updating post in draft for user %v", err)
		return err
	}

	return nil
}

func NewDraftRepository(db *sqlx.DB) DraftRepository {
	return draftRepository {
		db: db,
	}
}