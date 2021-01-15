package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/response"
	"post-api/service/test_helper"
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
}

func (suite *PostControllerTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestPostControllerTestSuite(t *testing.T) {
	suite.Run(t, new(PostControllerTest))
}

func (suite *PostControllerTest) TestPublishPost_WhenSuccess() {
	draftId := "1q2we3r"

	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId, "some-user").Return("some-url", nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/publish", bytes.NewBufferString(`{ "draft_id": "1q2we3r"}`))
	suite.postController.PublishPost(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(`{"status":"published","url":"some-url"}`, string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestPublishPost_WhenBadRequest() {
	draftId := "1q2we3r"

	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId, "some-user").Return("", nil).Times(0)
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

	suite.mockPostService.EXPECT().PublishPost(suite.context, draftId, "some-user").Return("", &constants.InternalServerError).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/publish", bytes.NewBufferString(`{ "draft_id": "1q2we3r"}`))
	suite.postController.PublishPost(suite.context)
	expectedErr := &constants.InternalServerError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestLike_WhenSuccess() {
	postID := "q2w3e4r5tqaz"

	likeCount := response.LikedByCount{LikeCount: 1}
	suite.mockPostService.EXPECT().LikePost("some-user", postID, suite.context).Return(likeCount, nil).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "q2w3e4r5tqaz",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/update-likes/q2w3e4r5tqaz", nil)
	suite.postController.Like(suite.context)
	jsonBytes, err := json.Marshal(likeCount)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *PostControllerTest) TestLike_WhenBadRequest() {
	postID := "1"

	likeCount := response.LikedByCount{LikeCount: 1}
	suite.mockPostService.EXPECT().LikePost("some-user", postID, suite.context).Return(likeCount, nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/update-likes/1", nil)
	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "1",
		},
	}
	suite.context.Params = params
	suite.postController.Like(suite.context)
	jsonBytes, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *PostControllerTest) TestLike_WhenLikeUpdateServiceFails() {
	postID := "q2w3e4r5tqaz"

	suite.mockPostService.EXPECT().LikePost("some-user", postID, suite.context).Return(response.LikedByCount{}, constants.StoryInternalServerError("something went wrong")).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "q2w3e4r5tqaz",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/update-likes/q2w3e4r5tqaz", nil)
	suite.postController.Like(suite.context)
	jsonBytes, err := json.Marshal(constants.StoryInternalServerError("something went wrong"))
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *PostControllerTest) TestLike_WhenLikeUpdateServiceFailsWithNotFoundPostForGivenPostID() {
	postID := "q2w3e4r5tqaz"

	suite.mockPostService.EXPECT().LikePost("some-user", postID, suite.context).Return(response.LikedByCount{}, &constants.PostNotFoundErr).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "q2w3e4r5tqaz",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/update-likes/q2w3e4r5tqaz", nil)
	suite.postController.Like(suite.context)
	jsonBytes, err := json.Marshal(&constants.PostNotFoundErr)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusNotFound, suite.recorder.Code)
}

func (suite *PostControllerTest) TestUnlike_WhenSuccess() {
	postID := "q2w3e4r5tqaz"

	likeCount := response.LikedByCount{LikeCount: 1}
	suite.mockPostService.EXPECT().UnlikePost("some-user", postID, suite.context).Return(likeCount, nil).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "q2w3e4r5tqaz",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/q2w3e4r5tqaz/unlike", nil)
	suite.postController.Unlike(suite.context)
	jsonBytes, err := json.Marshal(likeCount)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *PostControllerTest) TestUnlike_WhenBadRequest() {
	postID := "1"

	likeCount := response.LikedByCount{LikeCount: 1}
	suite.mockPostService.EXPECT().UnlikePost("some-user", postID, suite.context).Return(likeCount, nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/1/unlike", nil)
	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "1",
		},
	}
	suite.context.Params = params
	suite.postController.Unlike(suite.context)
	jsonBytes, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *PostControllerTest) TestUnlike_WhenUnlikeUpdateServiceFails() {
	postID := "q2w3e4r5tqaz"

	suite.mockPostService.EXPECT().UnlikePost("some-user", postID, suite.context).Return(response.LikedByCount{}, constants.StoryInternalServerError("something went wrong")).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "q2w3e4r5tqaz",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/q2w3e4r5tqaz/unlike", nil)
	suite.postController.Unlike(suite.context)
	jsonBytes, err := json.Marshal(constants.StoryInternalServerError("something went wrong"))
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *PostControllerTest) TestUnlike_WhenUnlikeUpdateServiceFailsWithNotFoundPostForGivenPostID() {
	postID := "q2w3e4r5tqaz"

	suite.mockPostService.EXPECT().UnlikePost("some-user", postID, suite.context).Return(response.LikedByCount{}, &constants.PostNotFoundErr).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "q2w3e4r5tqaz",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/q2w3e4r5tqaz/unlike", nil)
	suite.postController.Unlike(suite.context)
	jsonBytes, err := json.Marshal(&constants.PostNotFoundErr)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusNotFound, suite.recorder.Code)
}

func (suite *PostControllerTest) TestComment_WhenSuccess() {
	commentRequest := `{"post_uid": "1q2w3e4r5t6y","comment": "this is some comment"}`
	suite.mockPostService.EXPECT().CommentPost(suite.context, "some-user", "1q2w3e4r5t6y", "this is some comment").Return(nil).Times(1)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/comment", bytes.NewBufferString(commentRequest))

	suite.postController.Comment(suite.context)

	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *PostControllerTest) TestComment_WhenBadRequest() {
	commentRequest := `{"comment": "this is some comment"}`
	suite.mockPostService.EXPECT().CommentPost(suite.context, "some-user", "1q2w3e4r5t6y", "this is some comment").Return(nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/comment", bytes.NewBufferString(commentRequest))

	suite.postController.Comment(suite.context)
	expectedErr := &constants.PayloadValidationError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	suite.Equal(string(marshal), suite.recorder.Body.String())
}

func (suite *PostControllerTest) TestComment_WhenPostServiceCommentFails() {
	commentRequest := `{"post_uid": "1q2w3e4r5t6y","comment": "this is some comment"}`
	suite.mockPostService.EXPECT().CommentPost(suite.context, "some-user", "1q2w3e4r5t6y", "this is some comment").Return(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong)).Times(1)

	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/comment", bytes.NewBufferString(commentRequest))

	suite.postController.Comment(suite.context)
	marshal, err := json.Marshal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong))
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(marshal), suite.recorder.Body.String())
}

func (suite *PostControllerTest) TestGetPost_WhenSuccess() {
	postID := "1q2w3e4r5t6y"

	post := response.Post{
		PostID:                 "1q2w3e4r5t6y",
		PostData:               models.JSONString{},
		LikeCount:              1,
		CommentCount:           1,
		Interests:              []string{"Art", "Books", "Grammar"},
		AuthorID:               "some-user",
		PreviewImage:           "some-url",
		PublishedAt:            1234567890,
		IsViewerLiked:          true,
		IsViewerIsAuthor:       false,
		IsViewerFollowedAuthor: false,
	}
	suite.mockPostService.EXPECT().GetPost(suite.context, postID, "some-user").Return(post, nil).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "1q2w3e4r5t6y",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/1q2w3e4r5t6y", nil)
	suite.postController.GetPost(suite.context)
	jsonBytes, err := json.Marshal(post)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusOK, suite.recorder.Code)
}

func (suite *PostControllerTest) TestGetPost_WhenBadRequest() {
	postID := "1"

	post := response.Post{
		PostID:                 "1q2w3e4r5t6y",
		PostData:               models.JSONString{},
		LikeCount:              1,
		CommentCount:           1,
		Interests:              []string{"Art", "Books", "Grammar"},
		AuthorID:               "some-user",
		PreviewImage:           "some-url",
		PublishedAt:            1234567890,
		IsViewerLiked:          true,
		IsViewerIsAuthor:       false,
		IsViewerFollowedAuthor: false,
	}
	suite.mockPostService.EXPECT().GetPost(suite.context, postID, "some-user").Return(post, nil).Times(0)

	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/1q2w3e4r5t6y", nil)
	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "1",
		},
	}
	suite.context.Params = params
	suite.postController.GetPost(suite.context)
	jsonBytes, err := json.Marshal(constants.PayloadValidationError)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
}

func (suite *PostControllerTest) TestGetPost_WhenGetPostServiceFailsWithPostNotFoundErr() {
	postID := "q2w3e4r5tqaz"

	suite.mockPostService.EXPECT().GetPost(suite.context, postID, "some-user").Return(response.Post{}, &constants.PostNotFoundErr).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "q2w3e4r5tqaz",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/1q2w3e4r5t6y", nil)
	suite.postController.GetPost(suite.context)
	jsonBytes, err := json.Marshal(&constants.PostNotFoundErr)
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusNotFound, suite.recorder.Code)
}

func (suite *PostControllerTest) TestGetPost_WhenGetPostServiceFailsWithCommonError() {
	postID := "q2w3e4r5tqaz"

	suite.mockPostService.EXPECT().GetPost(suite.context, postID, "some-user").Return(response.Post{}, constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong)).Times(1)

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "q2w3e4r5tqaz",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/1q2w3e4r5t6y", nil)
	suite.postController.GetPost(suite.context)
	jsonBytes, err := json.Marshal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong))
	suite.Nil(err)
	suite.Equal(string(jsonBytes), suite.recorder.Body.String())
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
}

func (suite *PostControllerTest) TestMarkReadLater_WhenSuccess() {
	postId := "1q2w3e4r5t6y"

	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "1q2w3e4r5t6y",
		},
	}
	suite.context.Params = params
	suite.mockPostService.EXPECT().MarkReadLater(suite.context, postId, "some-user").Return(nil).Times(1)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/1q2w3e4r5t6y/read-later", nil)
	suite.postController.MarkReadLater(suite.context)
	suite.Equal(http.StatusOK, suite.recorder.Code)
	suite.Equal(`{"status":"success"}`, string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestMarkReadLater_WhenBadRequest() {
	draftId := "1q2w3e4r5t6y"

	suite.mockPostService.EXPECT().MarkReadLater(suite.context, draftId, "some-user").Return(nil).Times(0)
	suite.context.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/post/1q2w3e4r5t6y/read-later", nil)
	suite.postController.MarkReadLater(suite.context)
	suite.Equal(http.StatusBadRequest, suite.recorder.Code)
	expectedErr := &constants.PayloadValidationError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestMarkReadLater_WhenServiceFails() {
	draftId := "1q2w3e4r5t6y"

	suite.mockPostService.EXPECT().MarkReadLater(suite.context, draftId, "some-user").Return(&constants.InternalServerError).Times(1)
	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "1q2w3e4r5t6y",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/1q2w3e4r5t6y/read-later", nil)
	suite.postController.MarkReadLater(suite.context)
	expectedErr := &constants.InternalServerError
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(http.StatusInternalServerError, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *PostControllerTest) TestMarkReadLater_WhenServiceReturnsPostNotFound() {
	draftId := "1q2w3e4r5t6y"

	suite.mockPostService.EXPECT().MarkReadLater(suite.context, draftId, "some-user").Return(&constants.PostNotFoundErr).Times(1)
	params := gin.Params{
		gin.Param{
			Key:   "post_id",
			Value: "1q2w3e4r5t6y",
		},
	}
	suite.context.Params = params
	suite.context.Request, _ = http.NewRequest(http.MethodPost, "/api/v1/post/1q2w3e4r5t6y/read-later", nil)
	suite.postController.MarkReadLater(suite.context)
	expectedErr := &constants.PostNotFoundErr
	marshal, err := json.Marshal(expectedErr)
	suite.Nil(err)
	suite.Equal(http.StatusNotFound, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}
