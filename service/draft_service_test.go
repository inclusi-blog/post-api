package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/request"
	"testing"
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
