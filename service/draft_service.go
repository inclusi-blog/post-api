package service

//go:generate mockgen -source=draft_service.go -destination=./../mocks/mock_draft_service.go -package=mocks

import (
	"context"
	"database/sql"
	"post-api/constants"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/repository"
	"post-api/utils"

	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
)

type DraftService interface {
	SaveDraft(postData models.UpsertDraft, ctx context.Context) *golaerror.Error
	UpsertTagline(taglineRequest request.TaglineSaveRequest, ctx context.Context) *golaerror.Error
	UpsertInterests(interestRequest request.InterestsSaveRequest, ctx context.Context) *golaerror.Error
	GetDraft(draftUID string, ctx context.Context) (db.Draft, *golaerror.Error)
	GetAllDraft(allDraftReq models.GetAllDraftRequest, ctx context.Context) ([]db.AllDraft, error)
	SavePreviewImage(imageSaveRequest request.PreviewImageSaveRequest, ctx context.Context) *golaerror.Error
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
	logger.Infof("Saving post data to draft repository")
	err := service.draftRepository.SavePostDraft(postData, ctx)
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

func (service draftService) GetDraft(draftUID string, ctx context.Context) (db.Draft, *golaerror.Error) {

	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "GetDraft")

	logger.Infof("Calling service to get draft using draft ID %s", draftUID)

	draftData, err := service.draftRepository.GetDraft(ctx, draftUID)

	if err != nil {
		logger.Errorf("Error occurred while getting draft from repository %v", err)
		if err == sql.ErrNoRows {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", draftUID, err)
			return db.Draft{}, &constants.NoDraftFoundError
		}
		return db.Draft{}, &constants.PostServiceFailureError
	}

	logger.Info("Successfully stored got draft details")

	draft := db.Draft{
		DraftID:      draftData.DraftID,
		UserID:       draftData.UserID,
		PostData:     draftData.PostData,
		PreviewImage: draftData.PreviewImage.String,
		Tagline:      draftData.Tagline.String,
		Interest:     draftData.Interest,
	}

	return draft, nil
}

func (service draftService) SavePreviewImage(imageSaveRequest request.PreviewImageSaveRequest, ctx context.Context) *golaerror.Error {
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

func (service draftService) GetAllDraft(allDraftReq models.GetAllDraftRequest, ctx context.Context) ([]db.AllDraft, error) {

	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "GetAllDraft")

	logger.Infof("Calling service to get draft using user ID %s", allDraftReq.UserID)

	var allDraftData []db.AllDraft

	DraftData, err := service.draftRepository.GetAllDraft(ctx, allDraftReq)
	if err != nil {
		logger.Errorf("Error occurred while getting all draft from repository %v", err)
		if err == sql.ErrNoRows {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", allDraftReq.UserID, err)
			return allDraftData, &constants.NoDraftFoundError
		}
		return allDraftData, &constants.PostServiceFailureError
	}

	for _, val := range DraftData {
		var draft db.AllDraft
		draft.DraftID = val.DraftID
		draft.PostData = val.PostData
		draft.Tagline = val.Tagline
		draft.Interest = val.Interest
		draft.UserID = val.UserID

		title, err := utils.GetTitleFromSlateJson(ctx, val.PostData)
		if err != nil {
			logger.Errorf("Error occurred while converting title json to string %v .%v", val.DraftID, err)
			return allDraftData, &constants.ConnvertTitleToStringError
		}

		draft.TitleData = title

	}

	logger.Info("Successfully stored got draft details")

	return allDraftData, nil
}

func NewDraftService(repository repository.DraftRepository) DraftService {
	return draftService{
		draftRepository: repository,
	}
}
