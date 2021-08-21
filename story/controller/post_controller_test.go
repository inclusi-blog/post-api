package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	golaConstants "github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/constants"
	"github.com/gola-glitch/gola-utils/model"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/story/constants"
	"post-api/story/mocks"
	"post-api/validators"
	"testing"
)

type PostControllerTest struct {
	suite.Suite
	mockCtrl        *gomock.Controller
	recorder        *httptest.ResponseRecorder
	context         *gin.Context
	mockPostService *mocks.MockPostService
	postController  PostController
	userUUID        uuid.UUID
	idToken         model.IdToken
}

func (suite *PostControllerTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockPostService = mocks.NewMockPostService(suite.mockCtrl)
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.postController = NewPostController(suite.mockPostService)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("validPostUID", validators.ValidPostUID)
	}
	suite.userUUID = uuid.New()
	suite.idToken = model.IdToken{
		UserId:          suite.userUUID.String(),
		Username:        "dummy-user",
		Email:           "dummuser@gmail.com",
		Subject:         "gola",
		AccessTokenHash: "dummy-hash",
	}
	suite.context.Set(golaConstants.ContextDecryptedIdTokenKey, suite.idToken)
}

func (suite *PostControllerTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestPostControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PostControllerTest))
}

func (suite *PostControllerTest) TestPublishPost_WhenSuccess() {
	draftId := uuid.New()
	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId, suite.userUUID).Return(nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/publish?draft="+draftId.String(), nil)
	suite.postController.PublishPost(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(`{"status":"published"}`, string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestPublishPost_WhenTokenNotExists() {
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	draftId := uuid.New()
	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId, suite.userUUID).Return(nil).Times(0)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/publish?draft="+draftId.String(), nil)
	suite.postController.PublishPost(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	expectedErr := &constants.InternalServerError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestPublishPost_WhenBadRequest() {
	draftId := "1q2we3r"

	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId, suite.userUUID).Return(nil).Times(0)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/publish?draft", nil)
	suite.postController.PublishPost(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	expectedErr := &constants.PayloadValidationError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestPublishPost_WhenPublishPostFails() {
	draftId := uuid.New()

	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId, suite.userUUID).Return(&constants.InternalServerError).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/publish?draft="+draftId.String(), nil)
	suite.postController.PublishPost(suite.context)
	expectedErr := &constants.InternalServerError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestLikes_WhenSuccess() {
	postID := uuid.New()

	suite.mockPostService.EXPECT().LikePost(suite.context, postID, suite.userUUID).Return(nil).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post",
			Value: postID.String(),
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/like?post="+postID.String(), nil)

	suite.postController.Like(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *PostControllerTest) TestLikes_WhenIDTokenNotPresent() {
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	postID := uuid.New()

	suite.mockPostService.EXPECT().LikePost(suite.context, postID, suite.userUUID).Return(nil).Times(0)

	params := gin.Params{
		gin.Param{
			Key:   "post",
			Value: postID.String(),
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/like?post="+postID.String(), nil)

	suite.postController.Like(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	expectedErr := &constants.InternalServerError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestLikes_WhenBadRequest() {
	postID := uuid.New()

	suite.mockPostService.EXPECT().LikePost(suite.context, postID, suite.userUUID).Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/like?post=", nil)
	suite.postController.Like(suite.context)
	jsonBytes, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *PostControllerTest) TestLikes_WhenUnableToLikePost() {
	postID := uuid.New()

	suite.mockPostService.EXPECT().LikePost(suite.context, postID, suite.userUUID).Return(&constants.InternalServerError).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post",
			Value: postID.String(),
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/like?post="+postID.String(), nil)

	suite.postController.Like(suite.context)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	expectedErr := &constants.InternalServerError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}
