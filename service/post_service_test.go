package service

import (
	"context"
	"errors"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/response"
	"post-api/service/test_helper"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
)

type PostServiceTest struct {
	suite.Suite
	mockController            *gomock.Controller
	goContext                 context.Context
	mockPostsRepository       *mocks.MockPostsRepository
	mockDraftsRepository      *mocks.MockDraftRepository
	mockPreviewPostRepository *mocks.MockPreviewPostsRepository
	mockPostValidator         *mocks.MockPostValidator
	postService               PostService
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
	suite.mockPreviewPostRepository = mocks.NewMockPreviewPostsRepository(suite.mockController)
	suite.postService = NewPostService(suite.mockPostsRepository, suite.mockDraftsRepository, suite.mockPostValidator)
}

func (suite *PostServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *PostServiceTest) TestPublishPost_WhenSuccess() {
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	draftDB := db.DraftDB{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: "https://www.some-url.com",
		Tagline:      "",
		Interest:     []string{"sports", "economy"},
	}
	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interest:     []string{"sports", "economy"},
	}

	post := db.PublishPost{
		PUID:         "1231212",
		UserID:       "1",
		PostData:     draft.PostData,
		ReadTime:     22,
		Interest:     []string{"sports", "economy"},
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "",
		PreviewImage: "https://www.some-url.com",
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212", "some-user").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(nil).Times(1)
	err := suite.postService.PublishPost(suite.goContext, "1231212", "some-user")
	suite.Nil(err)
}

func (suite *PostServiceTest) TestPublishPost_WhenNoPreviewImageInDraft() {

	tmpPreviewImage := ""
	tmpTagLine := ""

	draftDB := db.DraftDB{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: "",
		Tagline:      "",
		Interest:     []string{"sports", "economy"},
	}
	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interest:     []string{"sports", "economy"},
	}

	post := db.PublishPost{
		PUID:         "1231212",
		UserID:       "1",
		PostData:     draft.PostData,
		ReadTime:     22,
		Interest:     []string{"sports", "economy"},
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "",
		PreviewImage: "https://www.some-url.com",
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212", "some-user").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "",
		ReadTime:     22,
		PreviewImage: "https://www.some-url.com",
	}, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(nil).Times(1)
	err := suite.postService.PublishPost(suite.goContext, "1231212", "some-user")
	suite.Nil(err)
}

func (suite *PostServiceTest) TestPublishPost_WhenGetDraftReturnsError() {
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212", "some-user").Return(db.DraftDB{}, errors.New("something went wrong")).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, db.PublishPost{}).Return(nil).Times(0)

	err := suite.postService.PublishPost(suite.goContext, "1231212", "some-user")
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenGetDraftReturnsNoDraftFoundError() {
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212", "some-user").Return(db.DraftDB{}, errors.New(constants.NoDraftFoundCode)).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, db.PublishPost{}).Return(nil).Times(0)

	err := suite.postService.PublishPost(suite.goContext, "1231212", "some-user")
	suite.NotNil(err)
	suite.Equal(&constants.NoDraftFoundError, err)
}

func (suite *PostServiceTest) TestPublishPost_WhenCreatePostReturnsError() {
	draftDB := db.DraftDB{
		DraftID: "1231212",
		UserID:  "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: "https://www.some-url.com",
		Tagline:      "this is some tag line",
		Interest:     []string{"sports", "economy"},
	}

	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &draftDB.PreviewImage,
		Tagline:      &draftDB.Tagline,
		Interest:     []string{"sports", "economy"},
	}

	post := db.PublishPost{
		PUID:         "1231212",
		UserID:       "some-user",
		PostData:     draft.PostData,
		ReadTime:     22,
		Interest:     []string{"sports", "economy"},
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "this is some tag line",
		PreviewImage: "https://www.some-url.com",
	}
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212", "some-user").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(errors.New("something went wrong")).Times(1)

	err := suite.postService.PublishPost(suite.goContext, "1231212", "some-user")
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenValidateDraftFails() {
	draftDB := db.DraftDB{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: "https://www.some-url.com",
		Tagline:      "",
		Interest:     []string{"sports", "economy"},
	}

	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""

	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interest:     []string{"sports", "economy"},
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212", "some-user").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{}, &constants.DraftValidationFailedError).Times(1)

	err := suite.postService.PublishPost(suite.goContext, "1231212", "some-user")
	suite.NotNil(err)
	suite.Equal(&constants.DraftValidationFailedError, err)
}

func (suite *PostServiceTest) TestLikePost_WhenSuccess() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().LikePost(postUUID, "some-user", suite.goContext).Return(nil).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikesCountByPostID(suite.goContext, postUUID).Return(int64(1), nil).Times(1)

	expectedCount := response.LikedByCount{LikeCount: 1}

	likeCount, err := suite.postService.LikePost("some-user", postUUID, suite.goContext)

	suite.Nil(err)
	suite.Equal(expectedCount, likeCount)
}

func (suite *PostServiceTest) TestLikePost_WhenRepositoryLikePostFails() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().LikePost(postUUID, "some-user", suite.goContext).Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikesCountByPostID(suite.goContext, postUUID).Return(int64(1), nil).Times(0)

	likeCount, err := suite.postService.LikePost("some-user", postUUID, suite.goContext)

	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), err)
	suite.Equal(response.LikedByCount{}, likeCount)
}

func (suite *PostServiceTest) TestLikePost_WhenGetCountByPostFails() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().LikePost(postUUID, "some-user", suite.goContext).Return(nil).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikesCountByPostID(suite.goContext, postUUID).Return(int64(0), errors.New(test_helper.ErrSomethingWentWrong)).Times(1)

	likeCount, err := suite.postService.LikePost("some-user", postUUID, suite.goContext)

	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), err)
	suite.Equal(response.LikedByCount{}, likeCount)
}

func (suite *PostServiceTest) TestUnlikePost_WhenSuccess() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().UnlikePost(suite.goContext, "some-user", postUUID).Return(nil).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikesCountByPostID(suite.goContext, postUUID).Return(int64(1), nil).Times(1)

	expectedCount := response.LikedByCount{LikeCount: 1}

	likeCount, err := suite.postService.UnlikePost("some-user", postUUID, suite.goContext)

	suite.Nil(err)
	suite.Equal(expectedCount, likeCount)
}

func (suite *PostServiceTest) TestUnlikePost_WhenRepositoryUnlikePostFails() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().UnlikePost(suite.goContext, "some-user", postUUID).Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikesCountByPostID(suite.goContext, postUUID).Return(int64(1), nil).Times(0)

	likeCount, err := suite.postService.UnlikePost("some-user", postUUID, suite.goContext)

	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), err)
	suite.Equal(response.LikedByCount{}, likeCount)
}

func (suite *PostServiceTest) TestUnlikePost_WhenGetCountByPostFails() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().UnlikePost(suite.goContext, "some-user", postUUID).Return(nil).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikesCountByPostID(suite.goContext, postUUID).Return(int64(0), errors.New(test_helper.ErrSomethingWentWrong)).Times(1)

	likeCount, err := suite.postService.UnlikePost("some-user", postUUID, suite.goContext)

	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), err)
	suite.Equal(response.LikedByCount{}, likeCount)
}
