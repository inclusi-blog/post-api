package service

//go:generate mockgen -source=draft_service.go -destination=./../mocks/mock_draft_service.go -package=mocks

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/inclusi-blog/gola-utils/model"
	"post-api/service"
	"post-api/story/constants"
	"post-api/story/models"
	"post-api/story/models/db"
	"post-api/story/models/request"
	"post-api/story/models/response"
	"post-api/story/repository"
	"post-api/story/utils"
	"time"

	"github.com/inclusi-blog/gola-utils/golaerror"
	"github.com/inclusi-blog/gola-utils/logging"
)

type DraftService interface {
	CreateDraft(ctx context.Context, draft models.CreateDraft) (uuid.UUID, error)
	UpdateDraft(postData models.UpsertDraft, ctx context.Context) *golaerror.Error
	UpsertInterests(interestRequest request.InterestsSaveRequest, ctx context.Context) *golaerror.Error
	UpsertTagline(taglineRequest request.TaglineSaveRequest, ctx context.Context) *golaerror.Error
	GetDraft(ctx context.Context, draftUID, userUUID uuid.UUID) (db.Draft, *golaerror.Error)
	SavePreviewImage(ctx context.Context, imageSaveRequest request.PreviewImageSaveRequest) *golaerror.Error
	SaveImage(ctx context.Context, imageSaveRequest request.PreviewImageSaveRequest) (string, *golaerror.Error)
	GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.DraftPreview, error)
	DeleteDraft(ctx context.Context, draftID, userUUID uuid.UUID) *golaerror.Error
	ValidateAndGetDraft(ctx context.Context, draftId uuid.UUID, user model.IdToken) (response.PreviewDraft, *golaerror.Error)
	GetDraftImage(ctx context.Context, draftID uuid.UUID, imageID uuid.UUID) (string, *golaerror.Error)
}

type draftService struct {
	draftRepository    repository.DraftRepository
	interestRepository repository.InterestsRepository
	validator          utils.PostValidator
	awsServices        service.AwsServices
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

	apiErr := draft.ConvertInterests(func(interests []string) *golaerror.Error {
		draft.InterestTags, err = service.interestRepository.GetInterestsForName(ctx, interests)
		if err != nil {
			logger.Errorf("unable to get interests %v", err)
			return constants.StoryInternalServerError("something went wrong")
		}
		return nil
	})

	if apiErr != nil {
		logger.Error("unable to get interests")
		return db.Draft{}, apiErr
	}

	if draft.PreviewImage != nil {
		*draft.PreviewImage, err = service.awsServices.GetObjectInS3(*draft.PreviewImage, time.Hour*time.Duration(6))
		if err != nil {
			logger.Errorf("unable to fetch preview image from s3 %v", err)
			return db.Draft{}, &constants.InternalServerError
		}
	}

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

func (service draftService) SaveImage(ctx context.Context, imageSaveRequest request.PreviewImageSaveRequest) (string, *golaerror.Error) {
	id := imageSaveRequest.DraftID
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "SaveImage")
	logger.Infof("Saving draft image for draft id %v", id)

	imageID, err := service.draftRepository.UpsertImage(ctx, imageSaveRequest)

	if err != nil {
		logger.Errorf("Error occurred while saving draft image to draft %v .%v", id, err)
		return "", constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully stored draft image for draft id %v", id)
	return imageID, nil
}

func (service draftService) GetDraftImage(ctx context.Context, draftID uuid.UUID, imageID uuid.UUID) (string, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "GetDraftImage")
	logger.Infof("Fetching draft image for draft id %v and image id %v", draftID.String(), imageID.String())

	uploadKey, err := service.draftRepository.GetDraftImage(ctx, draftID, imageID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Errorf("No draft image found for draft id %v and image id %v. Error %v", draftID, imageID, err)
			return "", &constants.ObjectNotFoundError
		}
		logger.Errorf("Error occurred while fetching draft image to draft %v and image id %v. Error %v", draftID, imageID, err)
		return "", constants.StoryInternalServerError(err.Error())
	}

	logger.Infof("Successfully fetched draft image for draft id %v", draftID)
	return uploadKey, nil
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
		apiErr := draft.ConvertInterests(func(interests []string) *golaerror.Error {
			draft.InterestTags, err = service.interestRepository.GetInterestsForName(ctx, interests)
			if err != nil {
				logger.Errorf("unable to get interests %v", err)
				return constants.StoryInternalServerError("something went wrong")
			}
			return nil
		})

		if apiErr != nil {
			logger.Error("unable to get interests")
			return nil, apiErr
		}
		oldTagline := ""
		if draft.Tagline != nil {
			oldTagline = *draft.Tagline
		}
		updatedDraft := db.DraftPreview{
			DraftID:   draft.DraftID,
			UserID:    draft.UserID,
			Data:      draft.Data,
			Tagline:   oldTagline,
			Interests: draft.InterestTags,
			CreatedAt: draft.CreatedAt,
		}

		title, tagline, err := utils.GetTitleAndTaglineFromData(ctx, draft.Data)
		if err != nil {
			logger.Errorf("Error occurred while converting title json to string %v .%v", draft.DraftID, err)
			return updatedDrafts, &constants.ConvertTitleToStringError
		}

		updatedDraft.Title = title
		if updatedDraft.Tagline == "" && tagline != "" {
			updatedDraft.Tagline = tagline
		}
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

func (service draftService) ValidateAndGetDraft(ctx context.Context, draftId uuid.UUID, user model.IdToken) (response.PreviewDraft, *golaerror.Error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftService").WithField("method", "ValidateAndGetDraft")

	logger.Infof("Fetching draftDB for validation of draftDB id %v", draftId)

	userUUID, _ := uuid.Parse(user.UserId)
	draft, err := service.draftRepository.GetDraft(ctx, draftId, userUUID)

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from db %v, Error %v", draftId, err)
		return response.PreviewDraft{}, constants.StoryInternalServerError(err.Error())
	}

	var selectedInterests []string
	apiErr := draft.ConvertInterests(func(interests []string) *golaerror.Error {
		selectedInterests = interests
		draft.InterestTags, err = service.interestRepository.GetInterestsForName(ctx, interests)
		if err != nil {
			logger.Errorf("unable to get interests %v", err)
			return constants.StoryInternalServerError("something went wrong")
		}
		return nil
	})

	if apiErr != nil {
		logger.Error("unable to get interests")
		return response.PreviewDraft{}, apiErr
	}

	metaData, draftValidationErr := service.validator.ValidateAndGetReadTime(draft, ctx)

	if draftValidationErr != nil {
		logger.Errorf("Error occurred while validating draft for draft id %v, Error %v", draftId, err)
		return response.PreviewDraft{}, draftValidationErr
	}

	previewDraft := mapDraftToPreviewMetaData(draft, metaData, draftId, user.Username, selectedInterests)
	if previewDraft.PreviewImage != "" {
		previewDraft.PreviewImage, err = service.awsServices.GetObjectInS3(previewDraft.PreviewImage, time.Hour*time.Duration(6))
		if err != nil {
			logger.Errorf("unable to fetch preview image from s3 %v", err)
			return response.PreviewDraft{}, &constants.InternalServerError
		}
	}
	return previewDraft, nil
}

func InternalServerError(err error, logger logging.GolaLoggerEntry) *golaerror.Error {
	if err != nil {
		logger.Errorf("Error occurred while saving draft data into draft repository %v", err)
		return &constants.InternalServerError
	}
	return nil
}

func NewDraftService(repository repository.DraftRepository, interestsRepository repository.InterestsRepository, validator utils.PostValidator, awsServices service.AwsServices) DraftService {
	return draftService{
		draftRepository:    repository,
		interestRepository: interestsRepository,
		validator:          validator,
		awsServices:        awsServices,
	}
}

func mapDraftToPreviewMetaData(draftDB db.Draft, metaData models.MetaData, draftId uuid.UUID, username string, interests []string) response.PreviewDraft {
	var previewDraft response.PreviewDraft

	previewDraft.Tagline = *draftDB.Tagline
	previewDraft.PreviewImage = *draftDB.PreviewImage
	if draftDB.Tagline == nil {
		previewDraft.Tagline = metaData.Tagline
	}
	if draftDB.PreviewImage == nil {
		previewDraft.PreviewImage = metaData.PreviewImage
	}
	previewDraft.Interest = interests
	previewDraft.DraftID = draftId
	previewDraft.Title = metaData.Title
	previewDraft.AuthorName = username
	return previewDraft
}
