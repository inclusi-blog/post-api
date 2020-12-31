package service

//go:generate mockgen -source=draft_service.go -destination=./../mocks/mock_draft_service.go -package=mocks

import (
	"context"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/constants"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/repository"
	"post-api/utils"
)

type DraftService interface {
	SaveDraft(postData models.UpsertDraft, ctx context.Context) *golaerror.Error
	UpsertTagline(taglineRequest request.TaglineSaveRequest, ctx context.Context) *golaerror.Error
	UpsertInterests(interestRequest request.InterestsSaveRequest, ctx context.Context) *golaerror.Error
	GetDraft(draftUID, userId string, ctx context.Context) (db.DraftDB, *golaerror.Error)
	GetAllDraft(allDraftReq models.GetAllDraftRequest, ctx context.Context) ([]db.AllDraft, error)
	SavePreviewImage(imageSaveRequest request.PreviewImageSaveRequest, ctx context.Context) *golaerror.Error
	DeleteDraft(draftID, userId string, ctx context.Context) *golaerror.Error
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
	err := service.draftRepository.IsDraftPresent(ctx, postData.DraftID, postData.UserID)
	if err != nil {
		if err.Error() == constants.NoDraftFoundCode {
			logger.Infof("No draft found for draft id %v, creating one", postData.DraftID)
			err := service.draftRepository.CreateNewPostWithData(postData, ctx)
			if err != nil {
				logger.Errorf("Error occurred while creating new post %v", err)
				return constants.StoryInternalServerError(err.Error())
			}
			return nil
		}
		logger.Infof("Error occurred while fetching draft existence for draft id %v", err)
		return constants.StoryInternalServerError(err.Error())
	}
	logger.Infof("Updating draft for draft id %v", postData.DraftID)
	err = service.draftRepository.UpdateDraft(postData, ctx)

	if err != nil {
		logger.Errorf("Error occurred while updating draft for draft id %v", postData.DraftID)
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully updated draft for draft id %v", postData.DraftID)

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

func (service draftService) GetDraft(draftUID, userId string, ctx context.Context) (db.DraftDB, *golaerror.Error) {

	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "GetDraft")

	logger.Infof("Calling service to get draft using draft ID %s", draftUID)

	draftData, err := service.draftRepository.GetDraft(ctx, draftUID, userId)

	if err != nil {
		if err.Error() == constants.NoDraftFoundCode {
			logger.Errorf("Error no draft found for draft id %v", draftUID)
			return db.DraftDB{}, &constants.NoDraftFoundError
		}
		logger.Errorf("Error occurred while getting draft from repository %v", err)
		return db.DraftDB{}, &constants.PostServiceFailureError
	}

	logger.Info("Successfully stored got draft details")

	return draftData, nil
}

func (service draftService) GetAllDraft(allDraftReq models.GetAllDraftRequest, ctx context.Context) ([]db.AllDraft, error) {

	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "GetAllDraft")

	logger.Infof("Calling service to get draft using user ID %s", allDraftReq.UserID)

	var allDraftData []db.AllDraft

	draftData, err := service.draftRepository.GetAllDraft(ctx, allDraftReq)
	if err != nil {
		if err.Error() == constants.NoDraftFoundCode {
			logger.Errorf("Error occurred while getting draft data, no draft found for draft id %v .%v", allDraftReq.UserID, err)
			return allDraftData, &constants.NoDraftFoundError
		}
		logger.Errorf("Error occurred while getting all draft from repository %v", err)
		return allDraftData, &constants.PostServiceFailureError
	}

	for _, dbDraft := range draftData {
		var draft db.AllDraft
		draft.DraftID = dbDraft.DraftID
		draft.PostData = dbDraft.PostData
		draft.Tagline = &dbDraft.Tagline
		draft.Interest = dbDraft.Interest
		draft.UserID = dbDraft.UserID

		title, err := utils.GetTitleFromSlateJson(ctx, dbDraft.PostData)
		if err != nil {
			logger.Errorf("Error occurred while converting title json to string %v .%v", dbDraft.DraftID, err)
			return allDraftData, &constants.ConvertTitleToStringError
		}

		draft.TitleData = title
		allDraftData = append(allDraftData, draft)
	}

	logger.Info("Successfully stored got draft details")

	return allDraftData, nil
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

func (service draftService) DeleteDraft(draftID, userId string, ctx context.Context) *golaerror.Error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "DeleteDraft")

	logger.Infof("Deleting draft from draft repository for draft id %v", draftID)

	err := service.draftRepository.DeleteDraft(ctx, draftID, userId)

	if err != nil {
		logger.Errorf("error occurred while deleting draft from draft repository for draft %v", draftID)
		return constants.StoryInternalServerError(err.Error())
	}

	logger.Info("Successfully deleted draft from draft repository")

	return nil
}

func NewDraftService(repository repository.DraftRepository) DraftService {
	return draftService{
		draftRepository: repository,
	}
}
