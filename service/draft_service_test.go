package service

import (
	"context"
	"database/sql"
	"errors"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"testing"

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

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsNil() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().SavePostDraft(newDraft, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsError() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().SavePostDraft(newDraft, suite.goContext).Return(errors.New("something went wrong in db")).Times(1)

	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)

	suite.NotNil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertTagline_WhenDraftRepositoryReturnsNoError() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "dummy-id",
		Tagline: "this is some tagline",
	}

	suite.mockDraftRepository.EXPECT().SaveTaglineToDraft(saveRequest, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.UpsertTagline(saveRequest, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertTagline_WhenDraftRepositoryReturnsError() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "dummy-id",
		Tagline: "this is some tagline",
	}

	suite.mockDraftRepository.EXPECT().SaveTaglineToDraft(saveRequest, suite.goContext).Return(errors.New("something went wrong")).Times(1)

	expectedError := suite.draftService.UpsertTagline(saveRequest, suite.goContext)

	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

//UpsertInterests Tests
func (suite *DraftServiceTest) TestUpsertInterests_WhenDraftRepositoryReturnsNoError() {
	saveRequest := request.InterestsSaveRequest{
		Interests: models.JSONString{
			JSONText: types.JSONText(`[
				{
				  "id": "1",
				  "name": "sports"
				},
				{
				  "id": "2",
				  "name": "economy"
				}
			  ]`),
		},
		DraftID: "121212",
		UserID:  "1",
	}

	suite.mockDraftRepository.EXPECT().SaveInterestsToDraft(saveRequest, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.UpsertInterests(saveRequest, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertInterests_WhenDraftRepositoryReturnsError() {
	saveRequest := request.InterestsSaveRequest{
		Interests: models.JSONString{
			JSONText: types.JSONText(`[
				{
				  "id": "1",
				  "name": "sports"
				},
				{
				  "id": "2",
				  "name": "economy"
				}
			  ]`),
		},
		DraftID: "121212",
		UserID:  "1",
	}

	suite.mockDraftRepository.EXPECT().SaveInterestsToDraft(saveRequest, suite.goContext).Return(errors.New("something went wrong")).Times(1)

	expectedError := suite.draftService.UpsertInterests(saveRequest, suite.goContext)

	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

//GetDraft Tests

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsNoError() {
	draftID := "121212"
	tmpTagline := "My first Data"
	tmpPreviewPost := "https://some-url.com"

	draft := db.DraftDB{
		DraftID:  "121212",
		UserID:   "12",
		PostData: models.JSONString{},
		PreviewImage: sql.NullString{
			String: "https://some-url.com",
			Valid:  true,
		},
		Tagline: sql.NullString{
			String: "My first Data",
			Valid:  true,
		},
		Interest: models.JSONString{JSONText: types.JSONText(`[{"id": "1","name":"sports"},{"id":"2","name":"economy"}]`)},
	}

	expectedDraft := db.Draft{
		DraftID:      "121212",
		UserID:       "12",
		PostData:     models.JSONString{},
		PreviewImage: &tmpPreviewPost,
		Tagline:      &tmpTagline,
		Interest:     models.JSONString{JSONText: types.JSONText(`[{"id": "1","name":"sports"},{"id":"2","name":"economy"}]`)},
	}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID).Return(draft, nil).Times(1)

	actualDraft, expectedError := suite.draftService.GetDraft(draftID, suite.goContext)
	suite.Equal(expectedDraft, actualDraft)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsError() {
	draftID := "121212"

	expectedDraft := db.Draft{}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID).Return(db.DraftDB{}, errors.New("something went wrong")).Times(1)

	draftData, expectedError := suite.draftService.GetDraft(draftID, suite.goContext)
	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsNoRowError() {
	draftID := "121212"

	expectedDraft := db.Draft{}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID).Return(db.DraftDB{}, sql.ErrNoRows).Times(1)

	draftData, expectedError := suite.draftService.GetDraft(draftID, suite.goContext)
	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.NoDraftFoundError, expectedError)
}

func (suite *DraftServiceTest) TestSavePreviewImage_WhenSuccess() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "1q2w3e4",
		PreviewImageUrl: "https://some-url",
	}

	suite.mockDraftRepository.EXPECT().UpsertPreviewImage(suite.goContext, imageSaveRequest).Return(nil).Times(1)

	err := suite.draftService.SavePreviewImage(imageSaveRequest, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftServiceTest) TestSavePreviewImage_WhenDbReturnsError() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "1q2w3e4",
		PreviewImageUrl: "https://some-url",
	}

	suite.mockDraftRepository.EXPECT().UpsertPreviewImage(suite.goContext, imageSaveRequest).Return(errors.New("something went wrong")).Times(1)

	err := suite.draftService.SavePreviewImage(imageSaveRequest, suite.goContext)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

//GetAllDraft Tests

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsNoError() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: 1,
		Limit:      5,
	}

	var allDraftActual []db.AllDraft
	var allDraft []db.Draft

	var allDraftExpected []db.AllDraft

	tmpTagline := "My first Data"

	draft := db.AllDraft{
		DraftID:   "121212",
		UserID:    "12",
		PostData:  models.JSONString{},
		TitleData: "https://some-url.com",
		Tagline:   &tmpTagline,
		Interest:  models.JSONString{JSONText: types.JSONText(`[{"id": "1","name":"sports"},{"id":"2","name":"economy"}]`)},
	}

	allDraftActual = append(allDraftActual, draft)

	expectedDraft := db.AllDraft{
		DraftID:   "121212",
		UserID:    "12",
		PostData:  models.JSONString{},
		TitleData: "https://some-url.com",
		Tagline:   &tmpTagline,
		Interest:  models.JSONString{JSONText: types.JSONText(`[{"id": "1","name":"sports"},{"id":"2","name":"economy"}]`)},
	}

	allDraftExpected = append(allDraftExpected, expectedDraft)

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, allDraftReq).Return(allDraft, nil).Times(1)

	allDraftExpected, expectedError := suite.draftService.GetAllDraft(allDraftReq, suite.goContext)
	suite.Equal(allDraftExpected, allDraftExpected)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsError() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: 1,
		Limit:      5,
	}
	var expectedDraft []db.AllDraft
	var allDraft []db.Draft

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, allDraftReq).Return(allDraft, errors.New("something went wrong")).Times(1)

	draftData, expectedError := suite.draftService.GetAllDraft(allDraftReq, suite.goContext)
	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsNoRowError() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: 1,
		Limit:      5,
	}

	var expectedDraft []db.AllDraft
	var allDraft []db.Draft

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, allDraftReq).Return(allDraft, sql.ErrNoRows).Times(1)

	draftData, expectedError := suite.draftService.GetAllDraft(allDraftReq, suite.goContext)
	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.NoDraftFoundError, expectedError)
}
