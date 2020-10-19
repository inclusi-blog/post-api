package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"testing"
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
