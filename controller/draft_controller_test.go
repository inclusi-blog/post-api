package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/http/httptest"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/service/test_helper"
	"post-api/validators"
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
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("validPostUID", validators.ValidPostUID)
	}
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

func (suite *DraftControllerTest) TestSaveInterests_WhenAPISuccess() {
	newInterest := request.InterestsSaveRequest{
		Interest: "sports",
		DraftID:  "121212",
		UserID:   "1",
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

func (suite *DraftControllerTest) TestDeleteInterest_WhenAPISuccess() {
	newInterest := request.InterestsSaveRequest{
		Interest: "sports",
		DraftID:  "121212",
		UserID:   "1",
	}

	suite.mockDraftService.EXPECT().DeleteInterest(suite.context, newInterest).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(newInterest)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/post/draft/delete-interests", bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.DeleteInterest(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestDeleteInterest_WhenBadRequest() {
	newInterst := request.InterestsSaveRequest{}

	requestBody := `{
		"interests": "sports",
		"user_id": "1"
	  }`

	suite.mockDraftService.EXPECT().DeleteInterest(suite.context, newInterst).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/draft/upsert-interests", bytes.NewBufferString(requestBody))

	marshal, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)

	suite.draftController.DeleteInterest(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestGetDraft_WhenAPISuccess() {
	DraftID := "q3w4e5r5t6y7"
	params := gin.Params{
		gin.Param{
			Key:   "draft_id",
			Value: "q3w4e5r5t6y7",
		},
	}
	suite.context.Params = params

	suite.mockDraftService.EXPECT().GetDraft(DraftID, "some-user", suite.context).Return(db.DraftDB{}, nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/draft/get-draft/q3w4e5r5t6y7", nil)
	suite.draftController.GetDraft(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetDraft_WhenBadRequest() {
	DraftID := "q3w4e5r5t6y"

	suite.mockDraftService.EXPECT().GetDraft(DraftID, "some-user", suite.context).Return(db.DraftDB{}, nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/draft/get-draft/q3w4e5r5t6y", nil)

	suite.draftController.GetDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetDraft_WhenServiceFails() {
	DraftID := "q3w4e5r5t6y7"
	params := gin.Params{
		gin.Param{
			Key:   "draft_id",
			Value: "q3w4e5r5t6y7",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/draft/get-draft/q3w4e5r5t6y7", nil)
	jsonBytes, err := json.Marshal(&constants.NoDraftFoundError)
	suite.Nil(err)

	suite.mockDraftService.EXPECT().GetDraft(DraftID, "some-user", suite.context).Return(db.DraftDB{}, &constants.NoDraftFoundError).Times(1)

	suite.draftController.GetDraft(suite.context)

	suite.Equal(http.StatusNotFound, suite.recorder.Code)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
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

func (suite *DraftControllerTest) TestGetAllDraft_WhenAPISuccess() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: 1,
		Limit:      5,
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

func (suite *DraftControllerTest) TestDeleteDraft_WhenSuccess() {
	suite.context.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/post/draft/q2w3e4r5t6y7", nil)
	params := gin.Params{
		gin.Param{
			Key:   "draft_id",
			Value: "q2w3e4r5t6y7",
		},
	}
	suite.context.Params = params

	suite.mockDraftService.EXPECT().DeleteDraft("q2w3e4r5t6y7", "some-user", suite.context).Return(nil).Times(1)

	suite.draftController.DeleteDraft(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(`{"status":"deleted"}`, suite.recorder.Body.String())
}

func (suite *DraftControllerTest) TestDeleteDraft_WhenBadRequest() {
	suite.context.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/post/draft/1", nil)

	suite.mockDraftService.EXPECT().DeleteDraft("1", "some-user", suite.context).Return(nil).Times(0)

	jsonBytes, err := json.Marshal(&constants.PayloadValidationError)
	suite.Nil(err)

	suite.draftController.DeleteDraft(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
}

func (suite *DraftControllerTest) TestDeleteDraft_WhenBadServiceFailsWithNotFound() {
	suite.context.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/post/draft/q2w3e4r5t6y7", nil)
	params := gin.Params{
		gin.Param{
			Key:   "draft_id",
			Value: "q2w3e4r5t6y7",
		},
	}
	suite.context.Params = params

	suite.mockDraftService.EXPECT().DeleteDraft("q2w3e4r5t6y7", "some-user", suite.context).Return(&constants.NoDraftFoundError).Times(1)

	jsonBytes, err := json.Marshal(&constants.NoDraftFoundError)
	suite.Nil(err)

	suite.draftController.DeleteDraft(suite.context)

	suite.Equal(http.StatusNotFound, suite.recorder.Code)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
}

func (suite *DraftControllerTest) TestDeleteDraft_WhenBadServiceFailsWithGenericError() {
	suite.context.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/post/draft/q2w3e4r5t6y7", nil)
	params := gin.Params{
		gin.Param{
			Key:   "draft_id",
			Value: "q2w3e4r5t6y7",
		},
	}
	suite.context.Params = params

	suite.mockDraftService.EXPECT().DeleteDraft("q2w3e4r5t6y7", "some-user", suite.context).Return(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong)).Times(1)

	jsonBytes, err := json.Marshal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong))
	suite.Nil(err)

	suite.draftController.DeleteDraft(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
}
