package repository

//go:generate mockgen -source=draft_repository.go -destination=./../mocks/mock_draft_repository.go -package=mocks

import (
	"context"
	"errors"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"

	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx"
)

type DraftRepository interface {
	SavePostDraft(draft models.UpsertDraft, ctx context.Context) error
	SaveTaglineToDraft(taglineSaveRequest request.TaglineSaveRequest, ctx context.Context) error
	SaveInterestsToDraft(interestsSaveRequest request.InterestsSaveRequest, ctx context.Context) error
	GetDraft(ctx context.Context, draftUID string) (db.DraftDB, error)
	GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.Draft, error)
	UpsertPreviewImage(ctx context.Context, saveRequest request.PreviewImageSaveRequest) error
}

const (
	SavePostDraft    = "INSERT INTO drafts (draft_id, user_id, post_data) VALUES($1, $2, $3) ON CONFLICT(draft_id) DO UPDATE SET POST_DATA = $4, UPDATED_AT = current_timestamp"
	SaveTagline      = "INSERT INTO drafts (draft_id, user_id, tagline) VALUES($1, $2, $3) ON CONFLICT(draft_id) DO UPDATE SET tagline = $4, UPDATED_AT = current_timestamp"
	SaveInterests    = "INSERT INTO drafts (draft_id, user_id, interest) VALUES($1, $2, $3) ON CONFLICT(draft_id) DO UPDATE SET interest = $4, UPDATED_AT = current_timestamp"
	FetchDraft       = "SELECT draft_id, user_id, tagline, preview_image, interest, post_data FROM DRAFTS WHERE draft_id = $1"
	SavePreviewImage = "INSERT INTO drafts (draft_id, user_id, preview_image) VALUES($1, $2, $3) ON CONFLICT(draft_id) DO UPDATE SET preview_image = $4, UPDATED_AT = current_timestamp"
	FetchAllDraft    = "SELECT draft_id, user_id, tagline, interest, post_data FROM DRAFTS WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
)

type draftRepository struct {
	db *sqlx.DB
}

func (repository draftRepository) UpsertPreviewImage(ctx context.Context, saveRequest request.PreviewImageSaveRequest) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("class", "UpsertPreviewImage")

	logger.Infof("Storing preview image for draft id %v", saveRequest.DraftID)

	result, err := repository.db.ExecContext(ctx, SavePreviewImage, saveRequest.DraftID, saveRequest.UserID, saveRequest.PreviewImageUrl, saveRequest.PreviewImageUrl)

	if err != nil {
		logger.Errorf("Error occurred while saving preview image of draft %v", saveRequest.DraftID)
		return err
	}

	affectedNoRows, err := result.RowsAffected()

	if err != nil {
		logger.Errorf("Error occurred while fetching affected rows for preview image update of draft %v", err)
		return err
	}

	if affectedNoRows != 1 {
		logger.Errorf("Error occurred while updating preview image for draft .%v", "More than one entry got updated")
		return errors.New("more than one entry got updated")
	}

	logger.Infof("Successfully updated the preview image for draft id %v", saveRequest.DraftID)

	return nil
}

func (repository draftRepository) GetDraft(ctx context.Context, draftUID string) (db.DraftDB, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "GetDraft")

	logger.Infof("Fetching draft from draft repository for the given draft id %v", draftUID)

	var draft db.DraftDB

	err := repository.db.GetContext(ctx, &draft, FetchDraft, draftUID)

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from draft repository %v", err)
		return db.DraftDB{}, err
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

// TODO Need test for this
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

// This method is for getting all the drafts of specific user
func (repository draftRepository) GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.Draft, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "GetAllDraft")

	logger.Infof("Fetching draft from draft repository for the given user id %v", allDraftReq.UserID)

	var allDraft []db.Draft

	rows, err := repository.db.QueryContext(ctx, FetchAllDraft, allDraftReq.UserID, allDraftReq.Limit, allDraftReq.StartValue)

	if err != nil {
		logger.Errorf("Error occurred while fetching all draft from draft repository %v", err)
		return allDraft, err
	}

	for rows.Next() {
		var draft db.Draft
		if scanErr := rows.Scan(&draft.DraftID, &draft.UserID, &draft.Tagline, &draft.Interest, &draft.PostData); scanErr != nil {
			return []db.Draft{}, scanErr
		}
		allDraft = append(allDraft, draft)
	}
	logger.Infof("Successfully fetching draft from draft repository for given user id %v", allDraftReq.UserID)

	return allDraft, nil
}

func NewDraftRepository(db *sqlx.DB) DraftRepository {
	return draftRepository{
		db: db,
	}
}
