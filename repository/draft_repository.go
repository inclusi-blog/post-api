package repository

//go:generate mockgen -source=draft_repository.go -destination=./../mocks/mock_draft_repository.go -package=mocks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx/types"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"post-api/constants"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/utils"
)

type DraftRepository interface {
	IsDraftPresent(ctx context.Context, draftId, userId string) error
	CreateNewPostWithData(draft models.UpsertDraft, ctx context.Context) error
	UpdateDraft(draft models.UpsertDraft, ctx context.Context) error
	SaveTaglineToDraft(taglineSaveRequest request.TaglineSaveRequest, ctx context.Context) error
	SaveInterestsToDraft(interestsSaveRequest request.InterestsSaveRequest, ctx context.Context) error
	DeleteInterest(ctx context.Context, saveRequest request.InterestsSaveRequest) error
	GetDraft(ctx context.Context, draftUID string, userId string) (db.DraftDB, error)
	UpsertPreviewImage(ctx context.Context, saveRequest request.PreviewImageSaveRequest) error
	DeleteDraft(ctx context.Context, draftUID string, userId string) error
	GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.DraftDB, error)
	UpdatePublishedStatus(ctx context.Context, draftId, userId string, transaction neo4j.Transaction) error
}

const (
	CreateNewDraft   = "MATCH (author:Person) where author.userId = $userId MERGE (draft:Draft{draftId: $draftId, data: $data})-[audit:CREATED_BY{createdAt: timestamp()}]->(author)"
	IsDraftPresent   = "MATCH (draft:Draft)-[:CREATED_BY]->(person:Person) where person.userId = $userId and draft.draftId = $draftId return draft.draftId"
	SavePostDraft    = "MATCH (draft:Draft)-[audit:CREATED_BY]->(author:Person) where author.userId = $userId and draft.draftId = $draftId set draft.data = $data"
	SaveTagline      = "MATCH (draft:Draft)-[audit:CREATED_BY]->(author:Person) where author.userId = $userId and draft.draftId = $draftId set draft.tagline = $tagline"
	SaveInterests    = "MATCH (draft:Draft{ draftId: $draftId})-[audit:CREATED_BY]->(author:Person{userId: $userId}) MATCH (interest:Interest { name: $interestName}) MERGE (draft)-[:FALLS_UNDER]->(interest) SET audit.updatedAt = timestamp()"
	DeleteInterest   = "MATCH (author:Person{userId: $userId})<-[:CREATED_BY]-(draft:Draft{draftId: $draftId})-[interest:FALLS_UNDER]->(tag:Interest{name: $interestName}) DELETE interest"
	FetchDraft       = "MATCH (draft:Draft)-[audit:CREATED_BY]->(author:Person) where draft.draftId = $draftId and author.userId = $userId OPTIONAL MATCH (tags:Interest)<-[:FALLS_UNDER]-(draft) return draft.draftId as draftId, author.userId as userId, draft.tagline as tagline, draft.previewImage as previewImage, collect(tags.name) as interests, draft.data as postData, CASE WHEN exists(draft.isPublished) THEN draft.isPublished ELSE false END AS isPublished, audit.createdAt as createdAt"
	SavePreviewImage = "MATCH (draft:Draft)-[audit:CREATED_BY]->(author:Person) where draft.draftId = $draftId and author.userId = $userId set draft.previewImage = $previewImage, audit.updatedAt = timestamp()"
	DeleteDraft      = "MATCH (draft:Draft)-[audit:CREATED_BY]->(author:Person) where draft.draftId = $draftId and author.userId = $userId detach delete audit, draft"
	FetchAllDraft    = "MATCH (draft:Draft)-[audit:CREATED_BY]->(author:Person) OPTIONAL match (draft)-[:FALLS_UNDER]->(tags:Interest) where author.userId = $userId return draft.draftId as draftId, author.userId as userId, draft.tagline as tagline, draft.previewImage as previewImage, collect(tags.name) as interests, draft.data as postData, audit.createdAt as createdAt order by createdAt desc skip $offset limit $limit"
	SetPublishStatus = "MATCH (draft:Draft)-[:CREATED_BY]->(author:Person) where author.userId = $userId and draft.draftId = $draftId set draft.isPublished = true return draft.draftId as draftId"
)

type draftRepository struct {
	db neo4j.Session
}

func (repository draftRepository) CreateNewPostWithData(draft models.UpsertDraft, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "CreateNewPostWithData")

	logger.Infof("Inserting or updating the existing post in draft for user %v", draft.UserID)

	result, err := repository.db.Run(CreateNewDraft, map[string]interface{}{
		"draftId": draft.DraftID,
		"userId":  draft.UserID,
		"data":    draft.PostData.String(),
	})

	if err != nil {
		logger.Errorf("Error occurred while updating post in draft for user %v", err)
		return err
	}

	_, err = result.Consume()
	if err != nil {
		logger.Errorf("Error occurred while updating post in draft for user %v", err)
		return err
	}

	return nil
}

func (repository draftRepository) IsDraftPresent(ctx context.Context, draftId string, userId string) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "IsDraftPresent")
	result, err := repository.db.Run(IsDraftPresent, map[string]interface{}{
		"userId":  userId,
		"draftId": draftId,
	})

	if err != nil {
		logger.Errorf("error occurred while fetching draft existence %v", err)
		return err
	}

	if result.Next() {
		return nil
	}
	return errors.New(constants.NoDraftFoundCode)
}

func (repository draftRepository) UpdateDraft(draft models.UpsertDraft, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "UpdateDraft")
	logger.Infof("Inserting or updating the existing post in draft for user %v", draft.UserID)

	result, err := repository.db.Run(SavePostDraft, map[string]interface{}{
		"draftId": draft.DraftID,
		"userId":  draft.UserID,
		"data":    draft.PostData.String(),
	})

	if err != nil {
		logger.Errorf("Error occurred while updating post in draft for user %v", err)
		return err
	}

	_, err = result.Consume()

	if err != nil {
		logger.Errorf("Error occurred while updating post in draft for user %v", err)
		return err
	}
	return nil
}

func (repository draftRepository) SaveTaglineToDraft(taglineSaveRequest request.TaglineSaveRequest, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "SaveTaglineToDraft")
	logger.Info("Inserting tagline or upserting to the draft for the given draft id")

	result, err := repository.db.Run(SaveTagline, map[string]interface{}{
		"tagline": taglineSaveRequest.Tagline,
		"userId":  taglineSaveRequest.UserID,
		"draftId": taglineSaveRequest.DraftID,
	})

	if err != nil {
		logger.Errorf("Error occurred while updating post tagline in draft for user %v", err)
		return err
	}

	_, err = result.Consume()

	if err != nil {
		logger.Errorf("Error occurred while updating post in draft for user %v", err)
		return err
	}

	logger.Infof("Successfully saved the tagline for draft id %v", taglineSaveRequest.Tagline)
	return nil
}

func (repository draftRepository) GetDraft(ctx context.Context, draftUID string, userId string) (db.DraftDB, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "GetDraft")

	logger.Infof("Fetching draft from draft repository for the given draft id %v", draftUID)

	result, err := repository.db.Run(FetchDraft, map[string]interface{}{
		"draftId": draftUID,
		"userId":  userId,
	})

	if err != nil {
		logger.Errorf("Error occurred while fetching draft from draft repository %v", err)
		return db.DraftDB{}, err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error occurred while fetching summary draft from draft repository %v", err)
		return db.DraftDB{}, err
	}

	if result.Next() {
		var draft db.DraftForm
		bindDbValues, err := utils.BindDbValues(result, draft)
		if err != nil {
			logger.Errorf("binding error %v", err)
			return db.DraftDB{}, err
		}

		jsonString, _ := json.Marshal(bindDbValues)
		err = json.Unmarshal(jsonString, &draft)
		logger.Infof("this is response %v", draft)
		draftDB := db.DraftDB{
			DraftID: draft.DraftID,
			UserID:  draft.UserID,
			PostData: models.JSONString{
				JSONText: types.JSONText(draft.PostData),
			},
			PreviewImage: draft.PreviewImage,
			Tagline:      draft.Tagline,
			Interest:     draft.Interest,
			IsPublished:  draft.IsPublished,
			CreatedAt:    draft.CreatedAt,
		}
		return draftDB, nil
	}

	logger.Infof("Successfully fetching draft from draft repository for given draft id %v", draftUID)

	return db.DraftDB{}, errors.New(constants.NoDraftFoundCode)
}

func (repository draftRepository) SaveInterestsToDraft(interestsSaveRequest request.InterestsSaveRequest, ctx context.Context) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "SaveInterestsToDraft")
	logger.Info("Inserting interests or upserting to the draft for the given draft id")

	result, err := repository.db.Run(SaveInterests, map[string]interface{}{
		"draftId":      interestsSaveRequest.DraftID,
		"userId":       interestsSaveRequest.UserID,
		"interestName": interestsSaveRequest.Interest,
	})

	if err != nil {
		logger.Errorf("Error occurred while updating post tagline in draft for user %v", err)
		return err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("error occured while fetching result summary %v", err)
		return err
	}

	logger.Infof("Successfully saved the Interests for draft id %v", interestsSaveRequest.Interest)
	return nil
}

func (repository draftRepository) DeleteInterest(ctx context.Context, deleteInterestRequest request.InterestsSaveRequest) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "DeleteInterest")
	logger.Infof("Deleting interest for draft %v", deleteInterestRequest.DraftID)

	result, err := repository.db.Run(DeleteInterest, map[string]interface{}{
		"draftId":      deleteInterestRequest.DraftID,
		"interestName": deleteInterestRequest.Interest,
		"userId":       deleteInterestRequest.UserID,
	})

	if err != nil {
		logger.Errorf("Error occurred while deleting draft interest %v", err)
		return err
	}

	_, err = result.Summary()
	if err != nil {
		logger.Errorf("Error occurred while fetching summary while deleting draft interest %v", err)
		return err
	}

	logger.Info("Successfully deleted draft interest")

	return nil
}

func (repository draftRepository) UpsertPreviewImage(ctx context.Context, saveRequest request.PreviewImageSaveRequest) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("class", "UpsertPreviewImage")

	logger.Infof("Storing preview image for draft id %v", saveRequest.DraftID)

	result, err := repository.db.Run(SavePreviewImage, map[string]interface{}{
		"draftId":      saveRequest.DraftID,
		"userId":       saveRequest.UserID,
		"previewImage": saveRequest.PreviewImageUrl,
	})

	if err != nil {
		logger.Errorf("Error occurred while saving preview image of draft %v", saveRequest.DraftID)
		return err
	}

	_, err = result.Summary()
	if err != nil {
		logger.Errorf("Error occurred while fetching affected result summary for preview image update of draft %v", err)
		return err
	}

	logger.Infof("Successfully updated the preview image for draft id %v", saveRequest.DraftID)

	return nil
}

func (repository draftRepository) DeleteDraft(ctx context.Context, draftId string, userId string) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("class", "DeleteDraft")

	result, err := repository.db.Run(DeleteDraft, map[string]interface{}{
		"draftId": draftId,
		"userId":  userId,
	})

	if err != nil {
		logger.Errorf("Error occurred while deleting draft of id %v for user %v", draftId, userId)
		return err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error occurred while validating result summary for deleting draft with id %v", draftId)
		return err
	}

	logger.Infof("Successfully deleted draft for the given draft id %v", draftId)
	return nil
}

func (repository draftRepository) GetAllDraft(ctx context.Context, allDraftReq models.GetAllDraftRequest) ([]db.DraftDB, error) {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "GetAllDraft")

	logger.Infof("Fetching draft from draft repository for the given user id %v", allDraftReq.UserID)

	var drafts []db.DraftDB
	result, err := repository.db.Run(FetchAllDraft, map[string]interface{}{
		"userId": allDraftReq.UserID,
		"offset": allDraftReq.StartValue,
		"limit":  allDraftReq.Limit,
	})

	if err != nil {
		logger.Errorf("Error occurred while fetching all draft from draft repository %v", err)
		return drafts, err
	}

	_, err = result.Summary()

	if err != nil {
		logger.Errorf("Error occurred while fetching result summary for fetch all draft %v", err)
		return drafts, err
	}

	for result.Next() {
		var draft db.DraftDB
		draftId := result.Record().GetByIndex(0)
		if draftId != nil {
			draft.DraftID = draftId.(string)
		}
		userID := result.Record().GetByIndex(1)
		if userID != nil {
			draft.UserID = userID.(string)
		}
		tagline := result.Record().GetByIndex(2)
		if tagline != nil {
			draft.Tagline = tagline.(string)
		}
		previewImage := result.Record().GetByIndex(3)
		if previewImage != nil {
			draft.PreviewImage = previewImage.(string)
		}
		interests := result.Record().GetByIndex(4)
		if interests != nil {
			for _, interest := range interests.([]interface{}) {
				draft.Interest = append(draft.Interest, interest.(string))
			}
		}
		postData := result.Record().GetByIndex(5)
		if postData != nil {
			postString := postData.(string)
			marshalBytes, err := json.Marshal(postString)
			err = json.Unmarshal(marshalBytes, &draft.PostData)
			if err != nil {
				logger.Errorf("invalid struct format while destructuring the post data %v", err)
				return []db.DraftDB{}, err
			}
		}

		createdTime := result.Record().GetByIndex(6)
		if createdTime != nil {
			draft.CreatedAt = createdTime.(int64)
		}

		drafts = append(drafts, draft)
	}

	if len(drafts) == 0 {
		logger.Errorf("Error no drafts found for user %v", allDraftReq.UserID)
		return []db.DraftDB{}, errors.New(constants.NoDraftFoundCode)
	}

	logger.Infof("Successfully fetching draft from draft repository for given user id %v", allDraftReq.UserID)

	return drafts, nil
}

func (repository draftRepository) UpdatePublishedStatus(ctx context.Context, draftId, userId string, transaction neo4j.Transaction) error {
	logger := logging.GetLogger(ctx).WithField("class", "DraftRepository").WithField("method", "UpdatePublishedStatus")
	logger.Infof("Writing draft %v status for post publish", draftId)

	result, err := transaction.Run(SetPublishStatus, map[string]interface{}{
		"draftId": draftId,
		"userId":  userId,
	})

	if err != nil {
		logger.Errorf("Error occurred while updating draft publish status for draft %v, Error %v", draftId, err)
		return err
	}

	_, err = result.Summary()
	if err != nil {
		logger.Errorf("Error occurred while fetching summary for draft publish status update for draft %v, Error %v", draftId, err)
		return err
	}

	if !result.Next() {
		return errors.New(constants.NoDraftFoundCode)
	}

	logger.Infof("Successfully updated publish status for draft %v", draftId)

	return nil
}

func NewDraftRepository(db neo4j.Session) DraftRepository {
	return draftRepository{
		db: db,
	}
}
