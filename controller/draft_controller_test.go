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
	}

	suite.mockDraftService.EXPECT().SaveDraft(newDraft, suite.context).Return(&constants.InternalServerError).Times(1)

	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/upsertDraft", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
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
			JSONText: types.JSONText(`[{"id":"1","name":"sports"},{"id":"2","name":"economy"}]`),
		},
		DraftID: "121212",
		UserID:  "1",
	}

	suite.mockDraftService.EXPECT().UpsertInterests(newInterest, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(newInterest)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/draft/upsert-interests", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveInterests(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveInterests_WhenBadRequest() {
	newInterst := request.InterestsSaveRequest{}

	requestBody := `{
		"interests": [{"id":"1","name":"sports"},{"id":"2","name":"economy"}],
		"user_id": "1"
	  }`

	suite.mockDraftService.EXPECT().UpsertInterests(newInterst, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/draft/upsert-interests", bytes.NewBufferString(requestBody))

	marshal, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)

	suite.draftController.SaveInterests(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

// GetDraft Test Scripts

func (suite *DraftControllerTest) TestGetDraft_WhenAPISuccess() {
	DraftID := "121212"
	params := gin.Params{
		gin.Param{
			Key:   "draft_id",
			Value: "121212",
		},
	}
	suite.context.Params = params

	suite.mockDraftService.EXPECT().GetDraft(DraftID, suite.context).Return(db.Draft{}, nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/draft/get-draft/121212", nil)
	suite.draftController.GetDraft(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetDraft_WhenBadRequest() {
	DraftID := ""

	suite.mockDraftService.EXPECT().GetDraft(DraftID, suite.context).Return(db.Draft{}, nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/draft/get-draft", nil)

	suite.draftController.GetDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSavePreviewImage_WhenAPISuccess() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "q12w3e",
		PreviewImageUrl: "http://www.some-url.com",
	}
	suite.mockDraftService.EXPECT().SavePreviewImage(imageSaveRequest, suite.context).Return(nil).Times(1)

	requestBytes, err := json.Marshal(imageSaveRequest)
	suite.Nil(err)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "api/v1/post/draft/update-preview-image", bytes.NewBufferString(string(requestBytes)))

	suite.draftController.SavePreviewImage(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	expectedResponse := `{"status":"success"}`
	suite.Equal(expectedResponse, string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestSavePreviewImage_WhenInvalidRequest() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "q12w3e",
		PreviewImageUrl: "http://www.some-url.com",
	}

	invalidRequest := `{"user_id": "1"}`

	suite.mockDraftService.EXPECT().SavePreviewImage(imageSaveRequest, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "api/v1/post/draft/update-preview-image", bytes.NewBufferString(invalidRequest))

	suite.draftController.SavePreviewImage(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	bytesData, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(string(bytesData), string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestSavePreviewImage_WhenServiceFails() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "q12w3e",
		PreviewImageUrl: "http://www.some-url.com",
	}

	requestBytes, err := json.Marshal(imageSaveRequest)
	suite.Nil(err)

	suite.mockDraftService.EXPECT().SavePreviewImage(imageSaveRequest, suite.context).Return(&constants.PostServiceFailureError).Times(1)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "api/v1/post/draft/update-preview-image", bytes.NewBufferString(string(requestBytes)))

	suite.draftController.SavePreviewImage(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	bytesData, err := json.Marshal(constants.PostServiceFailureError)
	suite.Nil(err)
	suite.Equal(string(bytesData), string(suite.recorder.Body.Bytes()))
}

// GetAllDraft Test Scripts

func (suite *DraftControllerTest) TestGetAllDraft_WhenAPISuccess() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: "1",
		Limit:      "5",
	}

	jsonBytes, err := json.Marshal(allDraftReq)
	suite.Nil(err)

	suite.mockDraftService.EXPECT().GetAllDraft(allDraftReq, suite.context).Return([]db.AllDraft{}, nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/draft/get-all-draft", bytes.NewBufferString(string(jsonBytes)))
	suite.draftController.GetAllDraft(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetAllDraft_WhenBadRequest() {

	requestBody := `{user_id:"1",start_value:"1",limit:1}`

	suite.mockDraftService.EXPECT().GetAllDraft(requestBody, suite.context).Return([]db.AllDraft{}, nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/draft/get-all-draft", bytes.NewBufferString(requestBody))

	suite.draftController.GetAllDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}
