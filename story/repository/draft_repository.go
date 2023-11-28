package repository

//go:generate mockgen -source=draft_repository.go -destination=./../mocks/mock_draft_repository.go -package=mocks

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	transaction "post-api/helper"
	"post-api/story/models"
	"post-api/story/models/db"
	"post-api/story/models/request"

	"github.com/inclusi-blog/gola-utils/logging"
	"github.com/jmoiron/sqlx"
)

type DraftRepository interface {
	SavePostDraft(draft models.UpsertDraft, ctx context.Context) error
	CreateDraft(ctx context.Context, draft models.CreateDraft) (uuid.UUID, error)
	SaveTaglineToDraft(taglineSaveRequest request.TaglineSaveRequest, ctx context.Context) error
	SaveInterestsToDraft(interestsSaveRequest request.InterestsSaveRequest, ctx context.Context) error
	GetDraftByUser(ctx context.Context, draftUID, userID uuid.UUID) (db.Draft, error)
	GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.Draft, error)
	UpsertPreviewImage(ctx context.Context, saveRequest request.PreviewImageSaveRequest) error
	UpsertImage(ctx context.Context, saveRequest request.PreviewImageSaveRequest) (string, error)
	DeleteDraft(ctx context.Context, draftUID, userUUID uuid.UUID) error
	UpdatePublishStatus(ctx context.Context, txn transaction.Transaction, draftUID, userID uuid.UUID, status bool) error
	GetDraftImage(ctx context.Context, draftID, imageID uuid.UUID) (string, error)
	DeleteDraftImages(ctx context.Context, draftUID uuid.UUID) error
	GetDraft(ctx context.Context, draftID uuid.UUID) (*db.Draft, error)
}

const (
	CreateDraft         = "insert into drafts (id, user_id, data) values(uuid_generate_v4(), $1, $2) returning id"
	SavePostDraft       = "update drafts set data = $1, updated_at = current_timestamp where id = $2 and user_id = $3"
	SaveTagline         = "update drafts set tagline = $1, updated_at = current_timestamp where id = $2 and user_id = $3"
	SaveInterests       = "update drafts set interests = $1, updated_at = current_timestamp where id = $2 and user_id = $3"
	FetchDraftByUser    = "select id, user_id, data, preview_image, tagline, interests from drafts where id = $1 and user_id = $2"
	FetchDraft          = "select id, user_id, data, preview_image, tagline, interests from drafts where id = $1"
	SavePreviewImage    = "update drafts set preview_image = $1, updated_at = current_timestamp where id = $2 and user_id = $3"
	FetchAllDraft       = "select id, user_id, data, preview_image, tagline, interests, created_at from drafts where user_id = $1 and is_published is false order by created_at desc limit $2 offset $3"
	DeleteDraft         = "delete from drafts where id = $1 and user_id = $2"
	DeleteDraftImages   = "delete from draft_images where draft_id = $1"
	UpdatePublishStatus = "update drafts set is_published = $1 where id = $2 and user_id = $3"
	InsertDraftImage    = "insert into draft_images(id, draft_id, upload_id) values (uuid_generate_v4(), $1, $2) returning id"
	GetDraftImage       = "select id, draft_id, upload_id from draft_images where id = $1 and draft_id = $2"
)

type draftRepository struct {
	db *sqlx.DB
}

func (repository draftRepository) CreateDraft(ctx context.Context, draft models.CreateDraft) (uuid.UUID, error) {
	var draftUUID uuid.UUID
	err := repository.db.GetContext(ctx, &draftUUID, CreateDraft, draft.UserID, draft.Data)

	if err != nil {
		return draftUUID, err
	}

	return draftUUID, nil
}

func (repository draftRepository) UpsertPreviewImage(ctx context.Context, saveRequest request.PreviewImageSaveRequest) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("class", "UpsertPreviewImage")

	logger.Infof("Storing preview image for draft id %v", saveRequest.DraftID)

	result, err := repository.db.ExecContext(ctx, SavePreviewImage, saveRequest.UploadID, saveRequest.DraftID, saveRequest.UserID)

	if err != nil {
		logger.Errorf("Error occurred while saving preview image of draft %v", saveRequest.DraftID)
		return err
	}

	affectedNoRows, err := result.RowsAffected()

	if err != nil {
		logger.Errorf("Error occurred while fetching affected rows for preview image update of draft %v", err)
		return err
	}

	if affectedNoRows == 0 {
		logger.Errorf("Error occurred while updating preview image for draft .%v", "More than one entry got updated")
		return errors.New("draft not found")
	}

	logger.Infof("Successfully updated the preview image for draft id %v", saveRequest.DraftID)

	return nil
}

func (repository draftRepository) GetDraftByUser(ctx context.Context, draftUID, userID uuid.UUID) (db.Draft, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "GetDraftByUser")

	logger.Infof("Fetching draft from draft repository for the given draft id %v", draftUID)

	var draft db.Draft

	err := repository.db.GetContext(ctx, &draft, FetchDraftByUser, draftUID, userID)

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from draft repository %v", err)
		return db.Draft{}, err
	}

	logger.Infof("Successfully fetching draft from draft repository for given draft id %v", draftUID)

	return draft, nil
}

func (repository draftRepository) GetDraft(ctx context.Context, draftUID uuid.UUID) (*db.Draft, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "GetDraft")

	logger.Infof("Fetching draft from draft repository for the given draft id %v", draftUID)

	var draft db.Draft

	err := repository.db.GetContext(ctx, &draft, FetchDraft, draftUID)

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from draft repository %v", err)
		return nil, err
	}

	logger.Infof("Successfully fetching draft from draft repository for given draft id %v", draftUID)

	return &draft, nil
}

func (repository draftRepository) SavePostDraft(draft models.UpsertDraft, ctx context.Context) error {
	logger := logging.GetLogger(ctx)
	log := logger.WithField("class", "DraftRepository").WithField("method", "SavePostDraft")

	log.Infof("Inserting or updating the existing post in draft for user %v", draft.UserID)

	result, err := repository.db.ExecContext(ctx, SavePostDraft, draft.Data, draft.DraftID, draft.UserID)

	if err != nil {
		log.Errorf("Error occurred while updating post in draft for user %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Infof("unable to get affected rows")
		return err
	}

	if rowsAffected == 0 {
		log.Error("no such draft")
		return errors.New("draft not found")
	}

	return nil
}

func (repository draftRepository) SaveTaglineToDraft(taglineSaveRequest request.TaglineSaveRequest, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "SaveTaglineToDraft")
	logger.Info("Inserting tagline or upserting to the draft for the given draft id")

	result, err := repository.db.ExecContext(ctx, SaveTagline, taglineSaveRequest.Tagline, taglineSaveRequest.DraftID, taglineSaveRequest.UserID)

	if err != nil {
		logger.Errorf("Error occurred while updating post tagline in draft for user %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.Errorf("unable to get affected rows %v", err)
		return err
	}

	if rowsAffected == 0 {
		logger.Errorf("draft not found")
		return errors.New("draft not found")
	}

	logger.Infof("Successfully saved the tagline for draft id %v", taglineSaveRequest.Tagline)
	return nil
}

func (repository draftRepository) SaveInterestsToDraft(interestsSaveRequest request.InterestsSaveRequest, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "SaveInterestsToDraft")
	logger.Info("Inserting interests or upserting to the draft for the given draft id")

	// TODO Need to check affected rows
	result, err := repository.db.ExecContext(ctx, SaveInterests, pq.Array(interestsSaveRequest.Interests), interestsSaveRequest.DraftID, interestsSaveRequest.UserID)

	if err != nil {
		logger.Errorf("Error occurred while updating post tagline in draft for user %v", err)
		return err
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		logger.Errorf("unable to get affected rows %v", err)
		return err
	}
	if rowAffected == 0 {
		logger.Errorf("no row affected, draft not found")
		return errors.New("draft not found")
	}

	logger.Infof("Successfully saved the Interests for draft id %v", interestsSaveRequest.Interests)
	return nil
}

func (repository draftRepository) GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.Draft, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "GetAllDraft")

	logger.Infof("Fetching draft from draft repository for the given user id %v", allDraftReq.UserID)

	var drafts []db.Draft

	err := repository.db.SelectContext(ctx, &drafts, FetchAllDraft, allDraftReq.UserID, allDraftReq.Limit, allDraftReq.StartValue)

	if err != nil {
		logger.Errorf("Error occurred while fetching all draft from draft repository %v", err)
		return drafts, err
	}

	logger.Infof("Successfully fetching draft from draft repository for given user id %v", allDraftReq.UserID)

	return drafts, nil
}

func (repository draftRepository) DeleteDraft(ctx context.Context, draftUID, userUUID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "").WithField("method", "DeleteDraft")

	logger.Infof("Deleting draft for the given draft uid %v", draftUID)

	result, err := repository.db.ExecContext(ctx, DeleteDraft, draftUID, userUUID)

	if err != nil {
		logger.Errorf("error occurred while deleting draft from draft repository for draft id %v", draftUID)
		return err
	}
	affectedRows, err := result.RowsAffected()

	if err != nil {
		logger.Errorf("error occurred while fetchin affected rows for delete draft for draft id %v", draftUID)
		return err
	}

	if affectedRows == 0 {
		logger.Errorf("draft not found or more than on draft deleted for draft id %v", draftUID)
		return sql.ErrNoRows
	}

	logger.Info("Successfully deleted draft")
	return nil
}

func (repository draftRepository) DeleteDraftImages(ctx context.Context, draftUID uuid.UUID) error {
	logger := logging.GetLogger(ctx).WithField("class", "").WithField("method", "DeleteDraft")

	logger.Infof("Deleting draft images for the given draft uid %v", draftUID)

	_, err := repository.db.ExecContext(ctx, DeleteDraftImages, draftUID)

	if err != nil {
		logger.Errorf("error occurred while deleting draft from draft repository for draft id %v. Error %v", draftUID, err)
		return err
	}

	logger.Info("Successfully deleted draft images")
	return nil
}

func (repository draftRepository) UpdatePublishStatus(ctx context.Context, txn transaction.Transaction, draftUID, userID uuid.UUID, status bool) error {
	logger := logging.GetLogger(ctx).WithField("class", "").WithField("method", "DeleteDraft")

	logger.Infof("Deleting draft for the given draft uid %v", draftUID)

	result, err := txn.ExecContext(ctx, UpdatePublishStatus, status, draftUID, userID)

	if err != nil {
		logger.Errorf("error occurred while deleting draft from draft repository for draft id %v", draftUID)
		return err
	}
	affectedRows, err := result.RowsAffected()

	if err != nil {
		logger.Errorf("error occurred while fetchin affected rows for delete draft for draft id %v", draftUID)
		return err
	}

	if affectedRows == 0 || affectedRows > 1 {
		logger.Errorf("draft not found or more than on draft deleted for draft id %v", draftUID)
		return errors.New("more than one row affected or no row affected")
	}

	logger.Info("Successfully deleted draft")
	return nil
}

func (repository draftRepository) UpsertImage(ctx context.Context, saveRequest request.PreviewImageSaveRequest) (string, error) {
	var imageID uuid.UUID
	logger := logging.GetLogger(ctx)
	logger.Infof("entering draft repository to update draft image for draft id %v", saveRequest.DraftID.String())
	err := repository.db.GetContext(ctx, &imageID, InsertDraftImage, saveRequest.DraftID, saveRequest.UploadID)

	if err != nil {
		logger.Errorf("Unable to insert draft image for draft id %v. Error %v", saveRequest.DraftID.String(), err)
		return "", err
	}

	logger.Infof("successfully inserted draft image for draft id %v", saveRequest.DraftID.String())

	return imageID.String(), nil
}

func (repository draftRepository) GetDraftImage(ctx context.Context, draftID, imageID uuid.UUID) (string, error) {
	type draftImage struct {
		ID       string `db:"id"`
		DraftID  string `db:"draft_id"`
		UploadID string `db:"upload_id"`
	}
	var image draftImage
	logger := logging.GetLogger(ctx)
	logger.Infof("entering draft repository to get draft image for draft id %v", draftID.String())
	err := repository.db.GetContext(ctx, &image, GetDraftImage, imageID, draftID)

	if err != nil {
		logger.Errorf("Unable to get draft image for draft id %v. Error %v", draftID.String(), err)
		return "", err
	}

	logger.Infof("successfully fetched the draft image for draft id %v", draftID.String())

	return image.UploadID, nil
}

func NewDraftRepository(db *sqlx.DB) DraftRepository {
	return draftRepository{
		db: db,
	}
}
