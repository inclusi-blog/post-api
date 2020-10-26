package service

import (
	"context"
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
		Target: "post",
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
		Target: "post",
	}

	suite.mockDraftRepository.EXPECT().SavePostDraft(newDraft, suite.goContext).Return(errors.New("something went wrong in db")).Times(1)

	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)

	suite.NotNil(expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsNilForTitleDraft() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		TitleData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
		Target: "title",
	}

	suite.mockDraftRepository.EXPECT().SaveTitleDraft(newDraft, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsErrorOnTitleDraft() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		TitleData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
		Target: "title",
	}

	suite.mockDraftRepository.EXPECT().SaveTitleDraft(newDraft, suite.goContext).Return(errors.New("something went wrong in db")).Times(1)

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

	actualDraft := db.Draft{DraftID: "121212",
		UserID:    "12",
		PostData:  models.JSONString{},
		TitleData: models.JSONString{},
		Tagline:   "My first Data",
		Interest: models.JSONString{JSONText: types.JSONText(`[
			  {
				"id": "1",
				"name": "sports"
			  },
			  {
				"id": "2",
				"name": "economy"
			  }
			]`)}}

	expectedDraft := db.Draft{DraftID: "121212",
		UserID:    "12",
		PostData:  models.JSONString{},
		TitleData: models.JSONString{},
		Tagline:   "My first Data",
		Interest: models.JSONString{JSONText: types.JSONText(`[
			  {
				"id": "1",
				"name": "sports"
			  },
			  {
				"id": "2",
				"name": "economy"
			  }
			]`)}}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID).Return(actualDraft, nil).Times(1)

	draftData, expectedError := suite.draftService.GetDraft(draftID, suite.goContext)

	suite.Equal(expectedDraft, draftData)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsError() {
	draftID := "121212"

	actualDraft := db.Draft{}
	expectedDraft := db.Draft{}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID).Return(actualDraft, errors.New("something went wrong")).Times(1)

	draftData, expectedError := suite.draftService.GetDraft(draftID, suite.goContext)
	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}
