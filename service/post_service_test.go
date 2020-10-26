package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/db"
	"post-api/service/test_helper"
	"testing"
)

type PostServiceTest struct {
	suite.Suite
	mockController       *gomock.Controller
	goContext            context.Context
	mockPostsRepository  *mocks.MockPostsRepository
	mockDraftsRepository *mocks.MockDraftRepository
	mockPostValidator    *mocks.MockPostValidator
	postService          PostService
}

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTest))
}

func (suite *PostServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.goContext = context.WithValue(context.Background(), "someKey", "someValue")
	suite.mockPostsRepository = mocks.NewMockPostsRepository(suite.mockController)
	suite.mockDraftsRepository = mocks.NewMockDraftRepository(suite.mockController)
	suite.mockPostValidator = mocks.NewMockPostValidator(suite.mockController)
	suite.postService = NewPostService(suite.mockPostsRepository, suite.mockDraftsRepository, suite.mockPostValidator)
}

func (suite *PostServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *PostServiceTest) TestPublishPost_WhenSuccess() {
	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline:  "",
		Interest: models.JSONString{},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		TitleData: draft.TitleData,
		ReadTime:  22,
		ViewCount: 0,
	}
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draft, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(&draft, suite.goContext).Return(22, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(nil).Times(1)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.Nil(err)
}

func (suite *PostServiceTest) TestPublishPost_WhenGetDraftReturnsError() {
	draft := db.Draft{}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draft, errors.New("something went wrong")).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, db.PublishPost{}).Return(nil).Times(0)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenGetDraftReturnsSqlNoRowError() {
	draft := db.Draft{}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draft, sql.ErrNoRows).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, db.PublishPost{}).Return(nil).Times(0)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.NotNil(err)
	suite.Equal(&constants.NoDraftFoundError, err)
}

func (suite *PostServiceTest) TestPublishPost_WhenCreatePostReturnsError() {
	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline:  "this is some tag line",
		Interest: models.JSONString{},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		TitleData: draft.TitleData,
		ReadTime:  22,
		ViewCount: 0,
	}
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draft, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(&draft, suite.goContext).Return(22, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(errors.New("something went wrong")).Times(1)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenValidateDraftFails() {
	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline:  "",
		Interest: models.JSONString{},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		TitleData: draft.TitleData,
		ReadTime:  22,
		ViewCount: 0,
	}

	expectedErr := constants.DraftValidationFailedError
	expectedErr.AdditionalData = "something went wrong"
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draft, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(&draft, suite.goContext).Return(0, errors.New("something went wrong")).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(nil).Times(0)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.NotNil(err)
	suite.Equal(&expectedErr, err)
}
