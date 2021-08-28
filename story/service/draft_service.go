package service

//go:generate mockgen -source=draft_service.go -destination=./../mocks/mock_draft_service.go -package=mocks

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"post-api/story/constants"
	"post-api/story/models"
	"post-api/story/models/db"
	"post-api/story/models/request"
	"post-api/story/repository"
	"post-api/story/utils"

	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
)

type DraftService interface {
	CreateDraft(ctx context.Context, draft models.CreateDraft) (uuid.UUID, error)
	UpdateDraft(postData models.UpsertDraft, ctx context.Context) *golaerror.Error
	UpsertInterests(interestRequest request.InterestsSaveRequest, ctx context.Context) *golaerror.Error
	UpsertTagline(taglineRequest request.TaglineSaveRequest, ctx context.Context) *golaerror.Error
	GetDraft(ctx context.Context, draftUID, userUUID uuid.UUID) (db.Draft, *golaerror.Error)
	SavePreviewImage(ctx context.Context, imageSaveRequest request.PreviewImageSaveRequest) *golaerror.Error
	GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.DraftPreview, error)
	DeleteDraft(ctx context.Context, draftID, userUUID uuid.UUID) *golaerror.Error
}

type draftService struct {
	draftRepository repository.DraftRepository
}

func (service draftService) CreateDraft(ctx context.Context, draft models.CreateDraft) (uuid.UUID, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "CreateDraft")
	logger.Info("creating draft")
	return service.draftRepository.CreateDraft(ctx, draft)
}

func (service draftService) UpdateDraft(postData models.UpsertDraft, ctx context.Context) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "UpdateDraft")
	logger.Infof("Saving post data to draft repository")
	err := service.draftRepository.SavePostDraft(postData, ctx)
	return InternalServerError(err, logger)
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

func (service draftService) GetDraft(ctx context.Context, draftUID, userUUID uuid.UUID) (db.Draft, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "GetDraft")

	logger.Infof("Calling service to get draft using draft Id %s", draftUID)

	draft, err := service.draftRepository.GetDraft(ctx, draftUID, userUUID)

	if err != nil {
		logger.Errorf("Error occurred while getting draft from repository %v", err)
		if err == sql.ErrNoRows {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", draftUID, err)
			return db.Draft{}, &constants.NoDraftFoundError
		}
		return db.Draft{}, &constants.PostServiceFailureError
	}

	logger.Info("Successfully stored got draft details")

	draft.ConvertInterests()

	return draft, nil
}

func (service draftService) SavePreviewImage(ctx context.Context, imageSaveRequest request.PreviewImageSaveRequest) *golaerror.Error {
	id := imageSaveRequest.DraftID
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "SavePreviewImage")
	logger.Infof("Saving preview image for draft id %v", id)

	err := service.draftRepository.UpsertPreviewImage(ctx, imageSaveRequest)

	if err != nil {
		logger.Errorf("Error occurred while saving preview image to draft %v .%v", id, err)
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully stored preview image for draft id %v", id)
	return nil
}

func (service draftService) GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.DraftPreview, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "GetAllDraft")

	logger.Infof("Calling service to get draft using user Id %s", allDraftReq.UserID)

	var updatedDrafts []db.DraftPreview

	drafts, err := service.draftRepository.GetAllDraft(ctx, allDraftReq)
	if err != nil {
		logger.Errorf("Error occurred while getting all draft from repository %v", err)
		if err == sql.ErrNoRows {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", allDraftReq.UserID, err)
			return updatedDrafts, &constants.NoDraftFoundError
		}
		return updatedDrafts, &constants.PostServiceFailureError
	}

	for _, draft := range drafts {
		draft.ConvertInterests()
		updatedDraft := db.DraftPreview{
			DraftID:   draft.DraftID,
			UserID:    draft.UserID,
			Data:      draft.Data,
			Tagline:   *draft.Tagline,
			Interests: draft.InterestTags,
			CreatedAt: draft.CreatedAt,
		}

		title, err := utils.GetTitleFromSlateJson(ctx, draft.Data)
		if err != nil {
			logger.Errorf("Error occurred while converting title json to string %v .%v", draft.DraftID, err)
			return updatedDrafts, &constants.ConvertTitleToStringError
		}

		updatedDraft.Title = title
		updatedDrafts = append(updatedDrafts, updatedDraft)
	}

	logger.Info("Successfully stored got draft details")

	return updatedDrafts, nil
}

func (service draftService) DeleteDraft(ctx context.Context, draftID, userUUID uuid.UUID) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "DeleteDraft")

	logger.Infof("Deleting draft from draft repository for draft id %v", draftID)

	err := service.draftRepository.DeleteDraft(ctx, draftID, userUUID)

	if err != nil {
		if err == sql.ErrNoRows {
			logger.Errorf("no draft found for draft id %v .Error %v", draftID, err)
			return &constants.NoDraftFoundError
		}
		logger.Errorf("error occurred while deleting draft from draft repository for draft %v", draftID)
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Info("Successfully deleted draft from draft repository")

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
