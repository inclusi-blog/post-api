package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/constants"
	"post-api/mocks"
	"testing"
)

type PostControllerTest struct {
	suite.Suite
	mockCtrl        *gomock.Controller
	recorder        *httptest.ResponseRecorder
	context         *gin.Context
	mockPostService *mocks.MockPostService
	postController  PostController
}

func (suite *PostControllerTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockPostService = mocks.NewMockPostService(suite.mockCtrl)
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.postController = NewPostController(suite.mockPostService)
}

func (suite *PostControllerTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestPostControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PostControllerTest))
}

func (suite *PostControllerTest) TestPublishPost_WhenSuccess() {
	draftId := "1q2we3r"

	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId).Return(nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/publish", bytes.NewBufferString(`{ "draft_id": "1q2we3r"}`))
	suite.postController.PublishPost(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(`{"status":"published"}`, string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestPublishPost_WhenBadRequest() {
	draftId := "1q2we3r"

	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId).Return(nil).Times(0)
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/publish", bytes.NewBufferString(`{ "draft_id": ""}`))
	suite.postController.PublishPost(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	expectedErr := &constants.PayloadValidationError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestPublishPost_WhenPublishPostFails() {
	draftId := "1q2we3r"

	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId).Return(&constants.InternalServerError).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/publish", bytes.NewBufferString(`{ "draft_id": "1q2we3r"}`))
	suite.postController.PublishPost(suite.context)
	expectedErr := &constants.InternalServerError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}
