package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	golaConstants "github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/constants"
	"github.com/gola-glitch/gola-utils/model"
	"github.com/google/uuid"
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
	emptyContext     *gin.Context
	mockDraftService *mocks.MockDraftService
	draftController  DraftController
	userUUID         uuid.UUID
	idToken          model.IdToken
}

func (suite *DraftControllerTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockDraftService = mocks.NewMockDraftService(suite.mockCtrl)
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.emptyContext, _ = gin.CreateTestContext(suite.recorder)
	suite.userUUID = uuid.New()
	suite.idToken = model.IdToken{
		UserId:          suite.userUUID.String(),
		Username:        "dummy-user",
		Email:           "dummuser@gmail.com",
		Subject:         "gola",
		AccessTokenHash: "dummy-hash",
	}
	suite.context.Set(golaConstants.ContextDecryptedIdTokenKey, suite.idToken)
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
	draftUUID := uuid.New()
	newDraft := models.UpsertDraft{
		DraftID: draftUUID,
		UserID:  suite.userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(`{"title":"hello"}`),
		},
	}

	suite.mockDraftService.EXPECT().UpdateDraft(newDraft, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)
	suite.context.Request, err = http.NewRequest(http.MethodPut, "/api/v1/draft?draft="+draftUUID.String(), bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveDraft_WhenIDTokenNotPresent() {
	draftUUID := uuid.New()
	newDraft := models.UpsertDraft{
		DraftID: draftUUID,
		UserID:  suite.userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(`{"title":"hello"}`),
		},
	}
	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)
	suite.emptyContext.Request, err = http.NewRequest(http.MethodPut, "/api/v1/draft?draft="+draftUUID.String(), bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)
	suite.draftController.SaveDraft(suite.emptyContext)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveDraft_WhenBadRequest() {
	newDraft := models.UpsertDraft{}

	requestBody := `{user_id:"1",post_data:{"title":"hello"}}`

	suite.mockDraftService.EXPECT().UpdateDraft(newDraft, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/draft", bytes.NewBufferString(requestBody))

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveDraft_WhenBadRequestNoPostData() {
	draftUUID := uuid.New()
	newDraft := models.UpsertDraft{}

	requestBody := `{post_data:{"title":"hello"}}`

	suite.mockDraftService.EXPECT().UpdateDraft(newDraft, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/draft?draft="+draftUUID.String(), bytes.NewBufferString(requestBody))

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveDraft_WhenServiceReturnsError() {
	draftUUID := uuid.New()
	newDraft := models.UpsertDraft{
		DraftID: draftUUID,
		UserID:  suite.userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(`{}`),
		},
	}

	suite.mockDraftService.EXPECT().UpdateDraft(newDraft, suite.context).Return(&constants.InternalServerError).Times(1)

	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPut, "/api/v1/draft?draft="+draftUUID.String(), bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveDraft(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveTagline_WhenAPISuccess() {
	draftUUID := uuid.New()
	newDraft := request.TaglineSaveRequest{
		UserID:  suite.userUUID,
		DraftID: draftUUID,
		Tagline: "this is some request",
	}

	suite.mockDraftService.EXPECT().UpsertTagline(newDraft, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(newDraft)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPut, "/api/v1/draft/tagline?draft="+draftUUID.String(), bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveTagline(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveTagline_WhenBadRequest() {
	newDraft := request.TaglineSaveRequest{}

	requestBody := `{user_id:"1",tagline: "some-tagline"}`

	suite.mockDraftService.EXPECT().UpsertTagline(newDraft, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/draft/tagline", bytes.NewBufferString(requestBody))

	suite.draftController.SaveTagline(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveTagline_WhenBadRequestInvalidRequestBody() {
	draftUUID := uuid.New()
	newDraft := request.TaglineSaveRequest{}
	requestBody := `{agline: "some-tagline"}`

	suite.mockDraftService.EXPECT().UpsertTagline(newDraft, suite.context).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/draft/tagline?draft="+draftUUID.String(), bytes.NewBufferString(requestBody))

	suite.draftController.SaveTagline(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveTagline_WhenIDTokenNotPresent() {
	newDraft := request.TaglineSaveRequest{}

	requestBody := `{user_id:"1",tagline: "some-tagline"}`

	suite.mockDraftService.EXPECT().UpsertTagline(newDraft, suite.emptyContext).Return(nil).Times(0)

	suite.emptyContext.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/draft/tagline", bytes.NewBufferString(requestBody))

	suite.draftController.SaveTagline(suite.emptyContext)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveTagline_WhenServiceReturnsError() {
	draftUUID := uuid.New()
	saveRequest := request.TaglineSaveRequest{
		UserID:  suite.userUUID,
		DraftID: draftUUID,
		Tagline: "this is some tagline",
	}

	suite.mockDraftService.EXPECT().UpsertTagline(saveRequest, suite.context).Return(&constants.PostServiceFailureError).Times(1)

	jsonBytes, err := json.Marshal(saveRequest)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPost, "/api/v1/draft/tagline?draft="+draftUUID.String(), bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveTagline(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveInterests_WhenAPISuccess() {
	draftUUID := uuid.New()
	newInterest := request.InterestsSaveRequest{
		Interests: []string{"Sports", "Economy"},
		DraftID:   draftUUID,
		UserID:    suite.userUUID,
	}

	suite.mockDraftService.EXPECT().UpsertInterests(newInterest, suite.context).Return(nil).Times(1)

	jsonBytes, err := json.Marshal(newInterest)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPut, "/api/v1/draft/interests?draft="+draftUUID.String(), bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveInterests(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveInterests_WhenUpsertInterestsFails() {
	draftUUID := uuid.New()
	newInterest := request.InterestsSaveRequest{
		Interests: []string{"Sports", "Economy"},
		DraftID:   draftUUID,
		UserID:    suite.userUUID,
	}

	suite.mockDraftService.EXPECT().UpsertInterests(newInterest, suite.context).Return(&constants.InternalServerError).Times(1)

	jsonBytes, err := json.Marshal(newInterest)
	suite.Nil(err)

	suite.context.Request, err = http.NewRequest(http.MethodPut, "/api/v1/draft/interests?draft="+draftUUID.String(), bytes.NewBufferString(string(jsonBytes)))
	suite.Nil(err)

	suite.draftController.SaveInterests(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSaveInterests_WhenBadRequest() {
	draftUUID := uuid.New()
	newInterest := request.InterestsSaveRequest{}

	requestBody := `{
		"interests": [{"id":"1","name":"sports"},{"id":"2","name":"economy"}],
		"user_id": "1"
	  }`
	suite.mockDraftService.EXPECT().UpsertInterests(newInterest, suite.context).Return(nil).Times(0)
	suite.context.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/draft/interests?draft="+draftUUID.String(), bytes.NewBufferString(requestBody))

	marshal, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)

	suite.draftController.SaveInterests(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestSaveInterests_WhenDraftIDIsInvalid() {
	newInterest := request.InterestsSaveRequest{}

	requestBody := `{
		"interests": [{"id":"1","name":"sports"},{"id":"2","name":"economy"}],
		"user_id": "1"
	  }`
	suite.mockDraftService.EXPECT().UpsertInterests(newInterest, suite.context).Return(nil).Times(0)
	suite.context.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/draft/interests?draft="+"some-invalid-id", bytes.NewBufferString(requestBody))

	marshal, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)

	suite.draftController.SaveInterests(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestSaveInterests_WhenIDTokenNotPresent() {
	draftUUID := uuid.New()
	newInterest := request.InterestsSaveRequest{}

	requestBody := `{
		"interests": [{"id":"1","name":"sports"},{"id":"2","name":"economy"}],
		"user_id": "1"
	  }`
	suite.mockDraftService.EXPECT().UpsertInterests(newInterest, suite.emptyContext).Return(nil).Times(0)
	suite.emptyContext.Request, _ = http.NewRequest(http.MethodPut, "/api/v1/draft/interests?draft="+draftUUID.String(), bytes.NewBufferString(requestBody))

	suite.draftController.SaveInterests(suite.emptyContext)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetDraft_WhenAPISuccess() {
	draftUUID := uuid.New()
	suite.mockDraftService.EXPECT().GetDraft(suite.context, draftUUID, suite.userUUID).Return(db.Draft{}, nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/draft?draft="+draftUUID.String(), nil)
	suite.draftController.GetDraft(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetDraft_WhenBadRequest() {
	draftUUID := uuid.New()

	suite.mockDraftService.EXPECT().GetDraft(suite.context, draftUUID, suite.userUUID).Return(db.Draft{}, nil).Times(0)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/draft", nil)

	suite.draftController.GetDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetDraft_WhenIDTokenNotPresent() {
	draftUUID := uuid.New()
	suite.mockDraftService.EXPECT().GetDraft(suite.context, draftUUID, suite.emptyContext).Return(db.Draft{}, nil).Times(0)
	suite.emptyContext.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/draft", nil)

	suite.draftController.GetDraft(suite.emptyContext)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetDraft_WhenServiceFails() {
	draftUUID := uuid.New()
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/draft?draft="+draftUUID.String(), nil)
	jsonBytes, err := json.Marshal(&constants.NoDraftFoundError)
	suite.Nil(err)

	suite.mockDraftService.EXPECT().GetDraft(suite.context, draftUUID, suite.userUUID).Return(db.Draft{}, &constants.NoDraftFoundError).Times(1)

	suite.draftController.GetDraft(suite.context)

	suite.Equal(http.StatusNotFound, suite.recorder.Code)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
}

func (suite *DraftControllerTest) TestSavePreviewImage_WhenAPISuccess() {
	draftUUID := uuid.New()
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          suite.userUUID,
		DraftID:         draftUUID,
		PreviewImageUrl: "http://www.some-url.com",
	}
	suite.mockDraftService.EXPECT().SavePreviewImage(suite.context, imageSaveRequest).Return(nil).Times(1)
	requestBytes, err := json.Marshal(imageSaveRequest)
	suite.Nil(err)
	suite.context.Request, _ = http.NewRequest(http.MethodPut, "api/v1/draft/preview-image?draft="+draftUUID.String(), bytes.NewBufferString(string(requestBytes)))

	suite.draftController.SavePreviewImage(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	expectedResponse := `{"status":"success"}`
	suite.Equal(expectedResponse, string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestSavePreviewImage_WhenInvalidRequest() {
	draftUUID := uuid.New()
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          suite.userUUID,
		DraftID:         draftUUID,
		PreviewImageUrl: "http://www.some-url.com",
	}
	suite.mockDraftService.EXPECT().SavePreviewImage(suite.context, imageSaveRequest).Return(nil).Times(0)
	suite.context.Request, _ = http.NewRequest(http.MethodPut, "api/v1/draft/preview-image?draft"+draftUUID.String(), nil)

	suite.draftController.SavePreviewImage(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	bytesData, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(string(bytesData), string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestSavePreviewImage_WhenInvalidRequestBody() {
	draftUUID := uuid.New()
	suite.context.Request, _ = http.NewRequest(http.MethodPut, "api/v1/draft/preview-image?draft="+draftUUID.String(), bytes.NewBufferString(`{}`))

	suite.draftController.SavePreviewImage(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	bytesData, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(string(bytesData), string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestSavePreviewImage_WhenIDTokenNotPresent() {
	draftUUID := uuid.New()
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          suite.userUUID,
		DraftID:         draftUUID,
		PreviewImageUrl: "http://www.some-url.com",
	}
	suite.mockDraftService.EXPECT().SavePreviewImage(suite.emptyContext, imageSaveRequest).Return(nil).Times(0)
	suite.emptyContext.Request, _ = http.NewRequest(http.MethodPut, "api/v1/draft/preview-image?draft"+draftUUID.String(), nil)

	suite.draftController.SavePreviewImage(suite.emptyContext)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestSavePreviewImage_WhenServiceFails() {
	draftUUID := uuid.New()
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          suite.userUUID,
		DraftID:         draftUUID,
		PreviewImageUrl: "http://www.some-url.com",
	}
	requestBytes, err := json.Marshal(imageSaveRequest)
	suite.Nil(err)
	suite.mockDraftService.EXPECT().SavePreviewImage(suite.context, imageSaveRequest).Return(&constants.PostServiceFailureError).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodPut, "api/v1/draft/preview-image?draft="+draftUUID.String(), bytes.NewBufferString(string(requestBytes)))

	suite.draftController.SavePreviewImage(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	bytesData, err := json.Marshal(constants.PostServiceFailureError)
	suite.Nil(err)
	suite.Equal(string(bytesData), string(suite.recorder.Body.Bytes()))
}

func (suite *DraftControllerTest) TestGetAllDraft_WhenAPISuccess() {
	draftRequest := models.GetAllDraftRequest{
		UserID:     suite.userUUID,
		StartValue: 1,
		Limit:      5,
	}

	jsonBytes, err := json.Marshal(draftRequest)
	suite.Nil(err)

	suite.mockDraftService.EXPECT().GetAllDraft(suite.context, draftRequest).Return([]db.DraftPreview{}, nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/drafts", bytes.NewBufferString(string(jsonBytes)))
	suite.draftController.GetAllDraft(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetAllDraft_WhenServiceFails() {
	draftRequest := models.GetAllDraftRequest{
		UserID:     suite.userUUID,
		StartValue: 1,
		Limit:      5,
	}

	jsonBytes, err := json.Marshal(draftRequest)
	suite.Nil(err)

	suite.mockDraftService.EXPECT().GetAllDraft(suite.context, draftRequest).Return([]db.DraftPreview{}, &constants.InternalServerError).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/drafts", bytes.NewBufferString(string(jsonBytes)))
	suite.draftController.GetAllDraft(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetAllDraft_WhenBadRequest() {
	requestBody := `{user_id:"1",start_value:"1",limit:1}`
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/drafts", bytes.NewBufferString(requestBody))

	suite.draftController.GetAllDraft(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestGetAllDraft_WhenIDTokenNotPresent() {
	requestBody := `{user_id:"1",start_value:"1",limit:1}`
	suite.emptyContext.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/drafts", bytes.NewBufferString(requestBody))

	suite.draftController.GetAllDraft(suite.emptyContext)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestDeleteDraft_WhenSuccess() {
	draftUUID := uuid.New()
	suite.context.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/draft?draft="+draftUUID.String(), nil)
	suite.mockDraftService.EXPECT().DeleteDraft(suite.context, draftUUID, suite.userUUID).Return(nil).Times(1)

	suite.draftController.DeleteDraft(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(`{"status":"deleted"}`, suite.recorder.Body.String())
}

func (suite *DraftControllerTest) TestDeleteDraft_WhenBadRequest() {
	suite.context.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/drafts", nil)
	jsonBytes, err := json.Marshal(&constants.PayloadValidationError)
	suite.Nil(err)

	suite.draftController.DeleteDraft(suite.context)

	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
}

func (suite *DraftControllerTest) TestDeleteDraft_WhenIDTokenNotPresent() {
	suite.emptyContext.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/drafts", nil)
	suite.draftController.DeleteDraft(suite.emptyContext)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *DraftControllerTest) TestDeleteDraft_WhenBadServiceFailsWithNotFound() {
	draftUUID := uuid.New()
	suite.context.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/draft?draft="+draftUUID.String(), nil)
	suite.mockDraftService.EXPECT().DeleteDraft(suite.context, draftUUID, suite.userUUID).Return(&constants.NoDraftFoundError).Times(1)

	jsonBytes, err := json.Marshal(&constants.NoDraftFoundError)
	suite.Nil(err)

	suite.draftController.DeleteDraft(suite.context)

	suite.Equal(http.StatusNotFound, suite.recorder.Code)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
}

func (suite *DraftControllerTest) TestDeleteDraft_WhenBadServiceFailsWithGenericError() {
	draftUUID := uuid.New()
	suite.context.Request, _ = http.NewRequest(http.MethodDelete, "/api/v1/draft?draft="+draftUUID.String(), nil)
	suite.mockDraftService.EXPECT().DeleteDraft(suite.context, draftUUID, suite.userUUID).Return(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong)).Times(1)
	jsonBytes, err := json.Marshal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong))
	suite.Nil(err)

	suite.draftController.DeleteDraft(suite.context)

	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
}
