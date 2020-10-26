package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
)

type DraftControllerTest struct {
	suite.Suite
	mockCtrl         *gomock.Controller
	recorder         *httptest.ResponseRecorder
	context          *gin.Context
	mockDraftService *mocks.MockDraftService
	draftController  DraftController
}

func (suite *DraftControllerTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockDraftService = mocks.NewMockDraftService(suite.mockCtrl)
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.draftController = NewDraftController(suite.mockDraftService)
}

func (suite *DraftControllerTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestDraftControllerTestSuite(t *testing.T) {
	suite.Run(t, new(DraftControllerTest))
}

func (suite *DraftControllerTest) TestSaveDraft_WhenAPISuccess() {
	newDraft := models.UpsertDraft{
		DraftID: "qwerty1234as",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title":"hello"}`),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(`{}`),
		},
		Target: "post",
	}

	suite.mockDraftService.EXPECT().SaveDraft(newDraft, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/upsertDraft", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveDraft_WhenBadRequest() {
	newDraft := models.UpsertDraft{}

	requestBody := `{user_id:"1",post_data:{"title":"hello"}}`

	suite.mockDraftService.EXPECT().SaveDraft(newDraft, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/upsertDraft", bytes.NewBufferString(requestBody))

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveDraft_WhenServiceReturnsError() {
	newDraft := models.UpsertDraft{
		DraftID: "qwerty1234as",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{}`),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(`{"title":"this is title"}`),
		},
		Target: "title",
	}

	suite.mockDraftService.EXPECT().SaveDraft(newDraft, suite.context).Return(&constants.InternalServerError).Times(1)

	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/upsertDraft", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveDraft_WhenTargetIsNotTitleOrPostReturnsBadRequest() {
	newDraft := models.UpsertDraft{
		DraftID: "qwerty1234as",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title":"hello"}`),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(`{}`),
		},
		Target: "hello",
	}

	suite.mockDraftService.EXPECT().SaveDraft(newDraft, suite.context).Return(nil).Times(0)

	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/upsertDraft", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveTagline_WhenAPISuccess() {
	newDraft := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "some-id",
		Tagline: "this is some request",
	}

	suite.mockDraftService.EXPECT().UpsertTagline(newDraft, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/draft/tagline", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveTagline(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveTagline_WhenBadRequest() {
	newDraft := request.TaglineSaveRequest{}

	requestBody := `{user_id:"1",tagline: "some-tagline"}`

	suite.mockDraftService.EXPECT().UpsertTagline(newDraft, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/draft/tagline", bytes.NewBufferString(requestBody))

	suite.draftController.SaveTagline(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveTagline_WhenServiceReturnsError() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "dummy-id",
		Tagline: "this is some tagline",
	}

	suite.mockDraftService.EXPECT().UpsertTagline(saveRequest, suite.context).Return(&constants.PostServiceFailureError).Times(1)

	jsonBytes, err := json.Marshal(saveRequest)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/draft/tagline", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveTagline(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

// SaveInterests Test Scripts

func (suite *DraftControllerTest) TestSaveInterests_WhenAPISuccess() {
	newInterest := request.InterestsSaveRequest{
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

	suite.mockDraftService.EXPECT().UpsertInterests(newInterest, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(newInterest)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/draft/upsertInterests", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.service.UpsertInterests(newInterest, suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveInterests_WhenBadRequest() {
	newInterst := request.InterestsSaveRequest{}

	requestBody := `{
		"interests": [
		  {
			"id": "1",
			"name": "sports"
		  },
		  {
			"id": "2",
			"name": "economy"
		  }
		],
		"draft_id": "121212"
	  }`

	suite.mockDraftService.EXPECT().UpsertInterests(newInterst, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/draft/upsertInterests", bytes.NewBufferString(requestBody))

	suite.draftController.SaveInterests(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

// GetDraft Test Scripts

func (suite *DraftControllerTest) TestGetDraft_WhenAPISuccess() {
	DraftID := "121212"

	suite.mockDraftService.EXPECT().GetDraft(DraftID, suite.context).Return(db.Draft{}, nil).Times(1)

	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/draft/get-draft?draft_id=121212", nil)

	suite.draftController.service.GetDraft(DraftID, suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetDraft_WhenBadRequest() {
	DraftID := ""

	suite.mockDraftService.EXPECT().GetDraft(DraftID, suite.context).Return(db.Draft{}, nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/draft/get-draft?draft_id=", nil)

	suite.draftController.GetDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}
