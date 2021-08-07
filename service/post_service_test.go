package service

import (
	"context"
	"database/sql"
	"errors"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
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
	suite.postService = NewPostService(suite.mockPostsRepository, suite.mockDraftsRepository, suite.mockPostValidator, suite.mockPreviewPostRepository)
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
		PreviewImage: sql.NullString{
			String: "https://www.some-url.com",
			Valid:  true,
		},
		Tagline: sql.NullString{
			String: "",
			Valid:  false,
		},
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}
	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		ReadTime:  22,
		ViewCount: 0,
	}

	previewPost := db.PreviewPost{
		PostID:       1,
		Title:        "Install apps via helm in kubernetes",
		Tagline:      *draft.Tagline,
		PreviewImage: "https://www.some-url.com",
		LikeCount:    0,
		CommentCount: 0,
		ViewTime:     0,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(int64(1), nil).Times(1)
	suite.mockPreviewPostRepository.EXPECT().SavePreview(suite.goContext, previewPost).Return(int64(1), nil).Times(1)
	suite.mockPostsRepository.EXPECT().SaveInitialLike(suite.goContext, int64(1)).Return(nil).Times(1)
	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.Nil(err)
}

func (suite *PostServiceTest) TestPublishPost_WhenSaveInitialLikeFails() {
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	draftDB := db.DraftDB{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: sql.NullString{
			String: "https://www.some-url.com",
			Valid:  true,
		},
		Tagline: sql.NullString{
			String: "",
			Valid:  false,
		},
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}
	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		ReadTime:  22,
		ViewCount: 0,
	}

	previewPost := db.PreviewPost{
		PostID:       1,
		Title:        "Install apps via helm in kubernetes",
		Tagline:      *draft.Tagline,
		PreviewImage: "https://www.some-url.com",
		LikeCount:    0,
		CommentCount: 0,
		ViewTime:     0,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(int64(1), nil).Times(1)
	suite.mockPreviewPostRepository.EXPECT().SavePreview(suite.goContext, previewPost).Return(int64(1), nil).Times(1)
	suite.mockPostsRepository.EXPECT().SaveInitialLike(suite.goContext, int64(1)).Return(errors.New("something went wrong")).Times(1)
	err := suite.postService.PublishPost(suite.goContext, "1231212")
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
		PreviewImage: sql.NullString{
			String: "",
			Valid:  false,
		},
		Tagline: sql.NullString{
			String: "",
			Valid:  false,
		},
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}
	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		ReadTime:  22,
		ViewCount: 0,
	}

	previewPost := db.PreviewPost{
		PostID:       1,
		Title:        "Install apps via helm in kubernetes",
		Tagline:      *draft.Tagline,
		PreviewImage: "https://www.some-url.com",
		LikeCount:    0,
		CommentCount: 0,
		ViewTime:     0,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "",
		ReadTime:     22,
		PreviewImage: "https://www.some-url.com",
	}, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(int64(1), nil).Times(1)
	suite.mockPreviewPostRepository.EXPECT().SavePreview(suite.goContext, previewPost).Return(int64(1), nil).Times(1)
	suite.mockPostsRepository.EXPECT().SaveInitialLike(suite.goContext, int64(1)).Return(nil).Times(1)
	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.Nil(err)
}

func (suite *PostServiceTest) TestPublishPost_WhenGetDraftReturnsError() {
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(db.DraftDB{}, errors.New("something went wrong")).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, db.PublishPost{}).Return(int64(1), nil).Times(0)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenGetDraftReturnsSqlNoRowError() {
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(db.DraftDB{}, sql.ErrNoRows).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, db.PublishPost{}).Return(int64(0), nil).Times(0)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.NotNil(err)
	suite.Equal(&constants.NoDraftFoundError, err)
}

func (suite *PostServiceTest) TestPublishPost_WhenCreatePostReturnsError() {
	draftDB := db.DraftDB{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: sql.NullString{
			String: "https://www.some-url.com",
			Valid:  true,
		},
		Tagline: sql.NullString{
			String: "this is some tag line",
			Valid:  false,
		},
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}

	draft := db.Draft{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &draftDB.PreviewImage.String,
		Tagline:      &draftDB.Tagline.String,
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		ReadTime:  22,
		ViewCount: 0,
	}
	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(int64(0), errors.New("something went wrong")).Times(1)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
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
		PreviewImage: sql.NullString{
			String: "https://www.some-url.com",
			Valid:  true,
		},
		Tagline: sql.NullString{
			String: "",
			Valid:  false,
		},
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
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
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		ReadTime:  22,
		ViewCount: 0,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{}, &constants.DraftValidationFailedError).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(int64(0), nil).Times(0)

	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.NotNil(err)
	suite.Equal(&constants.DraftValidationFailedError, err)
}

func (suite *PostServiceTest) TestPublishPost_WhenSavePreviewPostFails() {
	draftDB := db.DraftDB{
		DraftID: "1231212",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: sql.NullString{
			String: "https://www.some-url.com",
			Valid:  true,
		},
		Tagline: sql.NullString{
			String: "",
			Valid:  false,
		},
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
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
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}

	post := db.PublishPost{
		PUID:      "1231212",
		UserID:    "1",
		PostData:  draft.PostData,
		ReadTime:  22,
		ViewCount: 0,
	}

	previewPost := db.PreviewPost{
		PostID:       1,
		Title:        "Install apps via helm in kubernetes",
		Tagline:      *draft.Tagline,
		PreviewImage: "https://www.some-url.com",
		LikeCount:    0,
		CommentCount: 0,
		ViewTime:     0,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, "1231212").Return(draftDB, nil).Times(1)
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, post).Return(int64(1), nil).Times(1)
	suite.mockPreviewPostRepository.EXPECT().SavePreview(suite.goContext, previewPost).Return(int64(1), errors.New("something went wrong")).Times(1)
	err := suite.postService.PublishPost(suite.goContext, "1231212")
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestLikePost_WhenSuccess() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().GetPostID(suite.goContext, postUUID).Return(int64(1), nil).Times(1)
	suite.mockPostsRepository.EXPECT().AppendOrRemoveUserFromLikedBy(int64(1), int64(1), suite.goContext).Return(nil).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikeCountByPost(suite.goContext, int64(1)).Return(int64(1), nil).Times(1)

	expectedCount := request.LikedByCount{LikeCount: 1}

	likeCount, err := suite.postService.LikePost(int64(1), postUUID, suite.goContext)

	suite.Nil(err)
	suite.Equal(expectedCount, likeCount)
}

func (suite *PostServiceTest) TestLikePost_WhenGetPostIDFailsWithError() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().GetPostID(suite.goContext, postUUID).Return(int64(0), errors.New("something went wrong")).Times(1)
	suite.mockPostsRepository.EXPECT().AppendOrRemoveUserFromLikedBy(int64(0), int64(0), suite.goContext).Return(nil).Times(0)
	suite.mockPostsRepository.EXPECT().GetLikeCountByPost(suite.goContext, int64(0)).Return(int64(1), nil).Times(0)

	likeCount, err := suite.postService.LikePost(int64(1), postUUID, suite.goContext)

	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
	suite.Equal(request.LikedByCount{}, likeCount)
}

func (suite *PostServiceTest) TestLikePost_WhenGetPostIDFailsWithSQLNoRowsError() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().GetPostID(suite.goContext, postUUID).Return(int64(0), sql.ErrNoRows).Times(1)
	suite.mockPostsRepository.EXPECT().AppendOrRemoveUserFromLikedBy(int64(0), int64(0), suite.goContext).Return(nil).Times(0)
	suite.mockPostsRepository.EXPECT().GetLikeCountByPost(suite.goContext, int64(0)).Return(int64(1), nil).Times(0)

	likeCount, err := suite.postService.LikePost(int64(1), postUUID, suite.goContext)

	suite.NotNil(err)
	suite.Equal(&constants.PostNotFoundErr, err)
	suite.Equal(request.LikedByCount{}, likeCount)
}

func (suite *PostServiceTest) TestLikePost_WhenAppendOrRemoveUserLikeByFails() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().GetPostID(suite.goContext, postUUID).Return(int64(1), nil).Times(1)
	suite.mockPostsRepository.EXPECT().AppendOrRemoveUserFromLikedBy(int64(1), int64(1), suite.goContext).Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikeCountByPost(suite.goContext, int64(1)).Return(int64(1), nil).Times(0)

	likeCount, err := suite.postService.LikePost(int64(1), postUUID, suite.goContext)

	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), err)
	suite.Equal(request.LikedByCount{}, likeCount)
}

func (suite *PostServiceTest) TestLikePost_WhenGetCountByPostFails() {
	postUUID := "q1dsct52"

	suite.mockPostsRepository.EXPECT().GetPostID(suite.goContext, postUUID).Return(int64(1), nil).Times(1)
	suite.mockPostsRepository.EXPECT().AppendOrRemoveUserFromLikedBy(int64(1), int64(1), suite.goContext).Return(nil).Times(1)
	suite.mockPostsRepository.EXPECT().GetLikeCountByPost(suite.goContext, int64(1)).Return(int64(0), errors.New(test_helper.ErrSomethingWentWrong)).Times(1)

	likeCount, err := suite.postService.LikePost(int64(1), postUUID, suite.goContext)

	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), err)
	suite.Equal(request.LikedByCount{}, likeCount)
}
