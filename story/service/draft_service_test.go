package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"post-api/story/constants"
	"post-api/story/mocks"
	"post-api/story/models"
	"post-api/story/models/db"
	"post-api/story/models/request"
	"post-api/story/service/test_helper"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
)

type DraftServiceTest struct {
	suite.Suite
	mockController      *gomock.Controller
	goContext           context.Context
	mockDraftRepository *mocks.MockDraftRepository
	draftService        DraftService
}

func TestDraftServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DraftServiceTest))
}

func (suite *DraftServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockDraftRepository = mocks.NewMockDraftRepository(suite.mockController)
	suite.draftService = NewDraftService(suite.mockDraftRepository)
	suite.goContext = context.WithValue(context.Background(), "someKey", "someValue")
}

func (suite *DraftServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *DraftServiceTest) TestCreateDraft_WhenDBReturnsError() {
	userUUID := uuid.New()
	draft := models.CreateDraft{
		Data: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
		UserID: userUUID,
	}
	var newDraftUUID uuid.UUID
	suite.mockDraftRepository.EXPECT().CreateDraft(suite.goContext, draft).Return(newDraftUUID, errors.New("something went wrong")).Times(1)

	_, err := suite.draftService.CreateDraft(suite.goContext, draft)
	suite.NotNil(err)
	suite.Equal(errors.New("something went wrong"), err)
}

func (suite *DraftServiceTest) TestCreateDraft_WhenSuccessfullyCreated() {
	userUUID := uuid.New()
	draft := models.CreateDraft{
		Data: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
		UserID: userUUID,
	}
	newDraftUUID := uuid.New()
	suite.mockDraftRepository.EXPECT().CreateDraft(suite.goContext, draft).Return(newDraftUUID, nil).Times(1)

	draftUUID, err := suite.draftService.CreateDraft(suite.goContext, draft)
	suite.Nil(err)
	suite.Equal(newDraftUUID, draftUUID)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsNil() {
	newDraft := models.UpsertDraft{
		DraftID: uuid.New(),
		UserID:  uuid.New(),
		Data: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().SavePostDraft(newDraft, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.UpdateDraft(newDraft, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsError() {
	newDraft := models.UpsertDraft{
		DraftID: uuid.New(),
		UserID:  uuid.New(),
		Data: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().SavePostDraft(newDraft, suite.goContext).Return(errors.New("something went wrong in db")).Times(1)

	expectedError := suite.draftService.UpdateDraft(newDraft, suite.goContext)

	suite.NotNil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertInterests_WhenDraftRepositoryReturnsNoError() {
	saveRequest := request.InterestsSaveRequest{
		Interests: []string{"Sports", "Economy"},
		DraftID:   uuid.New(),
		UserID:    uuid.New(),
	}
	suite.mockDraftRepository.EXPECT().SaveInterestsToDraft(saveRequest, suite.goContext).Return(nil).Times(1)
	expectedError := suite.draftService.UpsertInterests(saveRequest, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertInterests_WhenDraftRepositoryReturnsError() {
	saveRequest := request.InterestsSaveRequest{
		Interests: []string{"Sports", "Economy"},
		DraftID:   uuid.New(),
		UserID:    uuid.New(),
	}

	suite.mockDraftRepository.EXPECT().SaveInterestsToDraft(saveRequest, suite.goContext).Return(errors.New("something went wrong")).Times(1)

	expectedError := suite.draftService.UpsertInterests(saveRequest, suite.goContext)

	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestUpsertTagline_WhenDraftRepositoryReturnsNoError() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  uuid.New(),
		DraftID: uuid.New(),
		Tagline: "this is some tagline",
	}

	suite.mockDraftRepository.EXPECT().SaveTaglineToDraft(saveRequest, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.UpsertTagline(saveRequest, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertTagline_WhenDraftRepositoryReturnsError() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  uuid.New(),
		DraftID: uuid.New(),
		Tagline: "this is some tagline",
	}

	suite.mockDraftRepository.EXPECT().SaveTaglineToDraft(saveRequest, suite.goContext).Return(errors.New("something went wrong")).Times(1)

	expectedError := suite.draftService.UpsertTagline(saveRequest, suite.goContext)

	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsNoError() {
	tmpTagline := "My first Data"
	tmpPreviewPost := "https://some-url.com"
	interests := "{Culture,Sports}"

	userUUID := uuid.New()
	draftUUID := uuid.New()
	expectedDraft := db.Draft{
		DraftID:      draftUUID,
		UserID:       userUUID,
		Data:         models.JSONString{JSONText: types.JSONText(`[{"title": "hello"}]`)},
		PreviewImage: &tmpPreviewPost,
		Tagline:      &tmpTagline,
		Interests:    &interests,
	}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(expectedDraft, nil).Times(1)

	actualDraft, expectedError := suite.draftService.GetDraft(suite.goContext, draftUUID, userUUID)
	expectedDraft.ConvertInterests(nil)
	suite.Equal(expectedDraft, actualDraft)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsError() {
	draftID := uuid.New()
	userUUID := uuid.New()

	expectedDraft := db.Draft{}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID, userUUID).Return(db.Draft{}, errors.New("something went wrong")).Times(1)

	draftData, expectedError := suite.draftService.GetDraft(suite.goContext, draftID, userUUID)
	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsNoRowError() {
	draftID := uuid.New()
	userUUID := uuid.New()

	expectedDraft := db.Draft{}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID, userUUID).Return(db.Draft{}, sql.ErrNoRows).Times(1)

	draftData, expectedError := suite.draftService.GetDraft(suite.goContext, draftID, userUUID)
	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.NoDraftFoundError, expectedError)
}

func (suite *DraftServiceTest) TestSavePreviewImage_WhenSuccess() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:   uuid.New(),
		DraftID:  uuid.New(),
		UploadID: "https://some-url",
	}

	suite.mockDraftRepository.EXPECT().UpsertPreviewImage(suite.goContext, imageSaveRequest).Return(nil).Times(1)

	err := suite.draftService.SavePreviewImage(suite.goContext, imageSaveRequest)
	suite.Nil(err)
}

func (suite *DraftServiceTest) TestSavePreviewImage_WhenDbReturnsError() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:   uuid.New(),
		DraftID:  uuid.New(),
		UploadID: "https://some-url",
	}

	suite.mockDraftRepository.EXPECT().UpsertPreviewImage(suite.goContext, imageSaveRequest).Return(errors.New("something went wrong")).Times(1)

	err := suite.draftService.SavePreviewImage(suite.goContext, imageSaveRequest)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsNoError() {
	userUUID := uuid.New()
	draftRequest := models.GetAllDraftRequest{
		UserID:     userUUID,
		StartValue: 1,
		Limit:      5,
	}

	tagline := "My first Data"
	previewImage := "some preview image"
	interests := "{Culture,Technology}"
	draftOneUUID := uuid.New()
	draftTwoUUID := uuid.New()
	draftThreeUUID := uuid.New()
	now := time.Now()
	drafts := []db.Draft{
		{
			DraftID: draftOneUUID,
			UserID:  userUUID,
			Data: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			PreviewImage: &previewImage,
			Tagline:      &tagline,
			Interests:    &interests,
			CreatedAt:    &now,
		},
		{
			DraftID: draftTwoUUID,
			UserID:  userUUID,
			Data: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			PreviewImage: &previewImage,
			Tagline:      &tagline,
			Interests:    &interests,
			CreatedAt:    &now,
		},
		{
			DraftID: draftThreeUUID,
			UserID:  userUUID,
			Data: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			PreviewImage: &previewImage,
			Tagline:      &tagline,
			Interests:    &interests,
			CreatedAt:    &now,
		},
	}

	expectedDrafts := []db.DraftPreview{
		{
			DraftID: draftOneUUID,
			UserID:  userUUID,
			Data: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			Title:     "தமிழ்நாட்டில் கொரோனா தொற்று பரவத் தொடங்கியபோது நோயாளிகளின் எண்ணிக்கை படிப்படியாக அதிகரித்து வந்தது. ",
			Tagline:   tagline,
			Interests: []string{"Culture", "Technology"},
			CreatedAt: &now,
		},
		{
			DraftID: draftTwoUUID,
			UserID:  userUUID,
			Data: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			Title:     "தமிழ்நாட்டில் கொரோனா தொற்று பரவத் தொடங்கியபோது நோயாளிகளின் எண்ணிக்கை படிப்படியாக அதிகரித்து வந்தது. ",
			Tagline:   tagline,
			Interests: []string{"Culture", "Technology"},
			CreatedAt: &now,
		},
		{
			DraftID: draftThreeUUID,
			UserID:  userUUID,
			Data: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			Title:     "தமிழ்நாட்டில் கொரோனா தொற்று பரவத் தொடங்கியபோது நோயாளிகளின் எண்ணிக்கை படிப்படியாக அதிகரித்து வந்தது. ",
			Tagline:   tagline,
			Interests: []string{"Culture", "Technology"},
			CreatedAt: &now,
		},
	}

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, draftRequest).Return(drafts, nil).Times(1)

	actualDrafts, expectedError := suite.draftService.GetAllDraft(suite.goContext, draftRequest)
	suite.Equal(expectedDrafts, actualDrafts)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsError() {
	userUUID := uuid.New()
	draftRequest := models.GetAllDraftRequest{
		UserID:     userUUID,
		StartValue: 1,
		Limit:      5,
	}
	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, draftRequest).Return([]db.Draft{}, errors.New("something went wrong")).Times(1)

	draftData, expectedError := suite.draftService.GetAllDraft(suite.goContext, draftRequest)
	suite.Nil(draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsInvalidData() {
	userUUID := uuid.New()
	draftOneUUID := uuid.New()
	draftRequest := models.GetAllDraftRequest{
		UserID:     userUUID,
		StartValue: 1,
		Limit:      5,
	}

	tagline := "My first Data"
	previewImage := "some preview image"
	interests := "{Culture,Technology}"

	draft := []db.Draft{
		{
			DraftID: draftOneUUID,
			UserID:  userUUID,
			Data: models.JSONString{
				JSONText: types.JSONText(`{`),
			},
			PreviewImage: &previewImage,
			Tagline:      &tagline,
			Interests:    &interests,
		},
	}

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, draftRequest).Return(draft, nil).Times(1)

	allDraftActual, expectedError := suite.draftService.GetAllDraft(suite.goContext, draftRequest)
	suite.Nil(allDraftActual)
	suite.NotNil(expectedError)
	suite.Equal(&constants.ConvertTitleToStringError, expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsNoRowError() {
	userUUID := uuid.New()
	draftRequest := models.GetAllDraftRequest{
		UserID:     userUUID,
		StartValue: 1,
		Limit:      5,
	}

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, draftRequest).Return([]db.Draft{}, sql.ErrNoRows).Times(1)

	actualDrafts, expectedError := suite.draftService.GetAllDraft(suite.goContext, draftRequest)
	suite.Nil(actualDrafts)
	suite.NotNil(expectedError)
	suite.Equal(&constants.NoDraftFoundError, expectedError)
}

func (suite *DraftServiceTest) TestDeleteDraft_WhenDraftRepositoryReturnNoError() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	suite.mockDraftRepository.EXPECT().DeleteDraft(suite.goContext, draftUUID, userUUID).Return(nil).Times(1)

	err := suite.draftService.DeleteDraft(suite.goContext, draftUUID, userUUID)
	suite.Nil(err)
}

func (suite *DraftServiceTest) TestDeleteDraft_WhenDraftRepositoryReturnsNotFoundError() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	suite.mockDraftRepository.EXPECT().DeleteDraft(suite.goContext, draftUUID, userUUID).Return(sql.ErrNoRows).Times(1)

	err := suite.draftService.DeleteDraft(suite.goContext, draftUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(&constants.NoDraftFoundError, err)
}

func (suite *DraftServiceTest) TestDeleteDraft_WhenDraftRepositoryReturnsGenericError() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	suite.mockDraftRepository.EXPECT().DeleteDraft(suite.goContext, draftUUID, userUUID).Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)

	err := suite.draftService.DeleteDraft(suite.goContext, draftUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), err)
}
