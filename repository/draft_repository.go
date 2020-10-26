package repository

//go:generate mockgen -source=draft_repository.go -destination=./../mocks/mock_draft_repository.go -package=mocks

import (
	"context"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"

	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
)

type DraftRepository interface {
	SavePostDraft(draft models.UpsertDraft, ctx context.Context) error
	SaveTitleDraft(draft models.UpsertDraft, ctx context.Context) error
	SaveTaglineToDraft(taglineSaveRequest request.TaglineSaveRequest, ctx context.Context) error
	SaveInterestsToDraft(interestsSaveRequest request.InterestsSaveRequest, ctx context.Context) error
	GetDraft(ctx context.Context, draftUID string) (db.Draft, error)
}

const (
	SavePostDraft  = "INSERT INTO drafts (draft_id, user_id, post_data) VALUES($1, $2, $3) ON CONFLICT(draft_id) DO UPDATE SET POST_DATA = $4, UPDATED_AT = current_timestamp"
	SaveTitleDraft = "INSERT INTO drafts (draft_id, user_id, title_data) VALUES($1, $2, $3) ON CONFLICT(draft_id) DO UPDATE SET TITLE_DATA = $4, UPDATED_AT = current_timestamp"
	SaveTagline    = "INSERT INTO drafts (draft_id, user_id, tagline) VALUES($1, $2, $3) ON CONFLICT(draft_id) DO UPDATE SET tagline = $4, UPDATED_AT = current_timestamp"
	SaveInterests  = "INSERT INTO drafts (draft_id, user_id, interest) VALUES($1, $2, $3) ON CONFLICT(draft_id) DO UPDATE SET interest = $4, UPDATED_AT = current_timestamp"
	FetchDraft     = "SELECT draft_id, user_id, tagline, interest, post_data, title_data FROM DRAFTS WHERE draft_id = $1"
)

type draftRepository struct {
	db *sqlx.DB
}

func (repository draftRepository) GetDraft(ctx context.Context, draftUID string) (db.Draft, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "GetDraft")

	logger.Infof("Fetching draft from draft repository for the given draft id %v", draftUID)

	var draft db.Draft

	err := repository.db.GetContext(ctx, &draft, FetchDraft, draftUID)

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from draft repository %v", err)
		return db.Draft{}, err
	}

	logger.Infof("Successfully fetching draft from draft repository for given draft id %v", draftUID)

	return draft, nil
}

func (repository draftRepository) SavePostDraft(draft models.UpsertDraft, ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("class", "DraftRepository").WithField("method", "SavePostDraft")

	log.Infof("Inserting or updating the existing post in draft for user %v", draft.UserID)

	_, err := repository.db.ExecContext(ctx, SavePostDraft, draft.DraftID, draft.UserID, draft.PostData, draft.PostData)

	if err != nil {
		log.Errorf("Error occurred while updating post in draft for user %v", err)
		return err
	}

	return nil
}

func (repository draftRepository) SaveTitleDraft(draft models.UpsertDraft, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "SaveTitleDraft")
	logger.Infof("Inserting or updating the existing post title for user %v", draft.UserID)

	_, err := repository.db.ExecContext(ctx, SaveTitleDraft, draft.DraftID, draft.UserID, draft.TitleData, draft.TitleData)

	if err != nil {
		logger.Errorf("Error occurred while updating post title in draft for user %v", err)
		return err
	}

	return nil
}

func (repository draftRepository) SaveTaglineToDraft(taglineSaveRequest request.TaglineSaveRequest, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "SaveTaglineToDraft")
	logger.Info("Inserting tagline or upserting to the draft for the given draft id")

	//TODO Need to check affected rows
	_, err := repository.db.ExecContext(ctx, SaveTagline, taglineSaveRequest.DraftID, taglineSaveRequest.UserID, taglineSaveRequest.Tagline, taglineSaveRequest.Tagline)

	if err != nil {
		logger.Errorf("Error occurred while updating post tagline in draft for user %v", err)
		return err
	}

	logger.Infof("Successfully saved the tagline for draft id %v", taglineSaveRequest.Tagline)
	return nil
}

func (repository draftRepository) SaveInterestsToDraft(interestsSaveRequest request.InterestsSaveRequest, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "SaveInterestsToDraft")
	logger.Info("Inserting interests or upserting to the draft for the given draft id")

	// TODO Need to check affected rows
	_, err := repository.db.ExecContext(ctx, SaveInterests, interestsSaveRequest.DraftID, interestsSaveRequest.UserID, interestsSaveRequest.Interests, interestsSaveRequest.Interests)

	if err != nil {
		logger.Errorf("Error occurred while updating post tagline in draft for user %v", err)
		return err
	}

	logger.Infof("Successfully saved the Interests for draft id %v", interestsSaveRequest.Interests)
	return nil
}

func NewDraftRepository(db *sqlx.DB) DraftRepository {
	return draftRepository{
		db: db,
	}
}
