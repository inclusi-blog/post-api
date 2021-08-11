package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/db"
	"post-api/service/test_helper"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
)

type PostServiceTest struct {
	suite.Suite
	mockController             *gomock.Controller
	goContext                  context.Context
	mockPostsRepository        *mocks.MockPostsRepository
	mockDraftsRepository       *mocks.MockDraftRepository
	mockAbstractPostRepository *mocks.MockAbstractPostRepository
	mockInterestsRepository    *mocks.MockInterestsRepository
	mockTransaction            *mocks.MockTransaction
	mockTransactionManager     *mocks.MockTransactionManager
	mockPostValidator          *mocks.MockPostValidator
	postService                PostService
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
	suite.mockAbstractPostRepository = mocks.NewMockAbstractPostRepository(suite.mockController)
	suite.mockTransaction = mocks.NewMockTransaction(suite.mockController)
	suite.mockTransactionManager = mocks.NewMockTransactionManager(suite.mockController)
	suite.mockInterestsRepository = mocks.NewMockInterestsRepository(suite.mockController)
	suite.postService = NewPostService(suite.mockPostsRepository, suite.mockDraftsRepository, suite.mockPostValidator, suite.mockAbstractPostRepository, suite.mockInterestsRepository, suite.mockTransactionManager)
}

func (suite *PostServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *PostServiceTest) TestPublishPost_WhenSuccess() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	postUUID := uuid.New()
	abstractPostUUID := uuid.New()
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}
	post := db.PublishPost{
		UserID: userUUID,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		DraftID: draftUUID,
	}
	abstractPost := db.AbstractPost{
		Model:        db.Model{},
		PostID:       postUUID,
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "",
		PreviewImage: "https://www.some-url.com",
		ViewTime:     22,
	}
	interestUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, nil).Times(1)
	draft.ConvertInterests()
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, suite.mockTransaction, post).Return(postUUID, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().GetInterestIDs(suite.goContext, draft.InterestTags).Return(interestUUIDs, nil).Times(1)
	suite.mockPostsRepository.EXPECT().AddInterests(suite.goContext, suite.mockTransaction, postUUID, interestUUIDs).Return(nil).Times(1)
	suite.mockAbstractPostRepository.EXPECT().Save(suite.goContext, suite.mockTransaction, abstractPost).Return(abstractPostUUID, nil).Times(1)
	suite.mockTransaction.EXPECT().Commit().Return(nil).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.Nil(err)
}

func (suite *PostServiceTest) TestPublishPost_WhenSuccessThereIsNoPreviewImage() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	postUUID := uuid.New()
	abstractPostUUID := uuid.New()
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: nil,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}
	post := db.PublishPost{
		UserID: userUUID,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		DraftID: draftUUID,
	}
	abstractPost := db.AbstractPost{
		Model:        db.Model{},
		PostID:       postUUID,
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "",
		PreviewImage: "https://www.some-url.com",
		ViewTime:     22,
	}
	interestUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, nil).Times(1)
	draft.ConvertInterests()
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "",
		ReadTime:     22,
		PreviewImage: "https://www.some-url.com",
	}, nil).Times(1)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, suite.mockTransaction, post).Return(postUUID, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().GetInterestIDs(suite.goContext, draft.InterestTags).Return(interestUUIDs, nil).Times(1)
	suite.mockPostsRepository.EXPECT().AddInterests(suite.goContext, suite.mockTransaction, postUUID, interestUUIDs).Return(nil).Times(1)
	suite.mockAbstractPostRepository.EXPECT().Save(suite.goContext, suite.mockTransaction, abstractPost).Return(abstractPostUUID, nil).Times(1)
	suite.mockTransaction.EXPECT().Commit().Return(nil).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.Nil(err)
}

func (suite *PostServiceTest) TestPublishPost_WhenGetDraftFails() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, errors.New("something went wrong")).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenNoDraftFound() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, sql.ErrNoRows).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(&constants.NoDraftFoundError, err)
}

func (suite *PostServiceTest) TestPublishPost_WhenValidationFails() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, nil).Times(1)
	draft.ConvertInterests()
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{}, &constants.DraftValidationFailedError).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.NotNil(&constants.DraftValidationFailedError, err)
}

func (suite *PostServiceTest) TestPublishPost_WhenCreatePostFails() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	var postUUID uuid.UUID
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}
	post := db.PublishPost{
		UserID: userUUID,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		DraftID: draftUUID,
	}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, nil).Times(1)
	draft.ConvertInterests()
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, suite.mockTransaction, post).Return(postUUID, errors.New("something went wrong")).Times(1)
	suite.mockTransaction.EXPECT().Rollback().Return(nil).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenGetInterestsFails() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	postUUID := uuid.New()
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}
	post := db.PublishPost{
		UserID: userUUID,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		DraftID: draftUUID,
	}
	var interestUUIDs []uuid.UUID

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, nil).Times(1)
	draft.ConvertInterests()
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, suite.mockTransaction, post).Return(postUUID, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().GetInterestIDs(suite.goContext, draft.InterestTags).Return(interestUUIDs, errors.New("something went wrong")).Times(1)
	suite.mockTransaction.EXPECT().Rollback().Return(nil).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenAddInterestsFails() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	postUUID := uuid.New()
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}
	post := db.PublishPost{
		UserID: userUUID,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		DraftID: draftUUID,
	}
	interestUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, nil).Times(1)
	draft.ConvertInterests()
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, suite.mockTransaction, post).Return(postUUID, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().GetInterestIDs(suite.goContext, draft.InterestTags).Return(interestUUIDs, nil).Times(1)
	suite.mockPostsRepository.EXPECT().AddInterests(suite.goContext, suite.mockTransaction, postUUID, interestUUIDs).Return(errors.New("something went wrong")).Times(1)
	suite.mockTransaction.EXPECT().Rollback().Return(nil).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestPublishPost_WhenAbstractPostSaveFails() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	postUUID := uuid.New()
	var abstractPostUUID uuid.UUID
	tmpPreviewImage := "https://www.some-url.com"
	tmpTagLine := ""
	interests := "{sports,economy}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: &tmpPreviewImage,
		Tagline:      &tmpTagLine,
		Interests:    &interests,
	}
	post := db.PublishPost{
		UserID: userUUID,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		DraftID: draftUUID,
	}
	abstractPost := db.AbstractPost{
		Model:        db.Model{},
		PostID:       postUUID,
		Title:        "Install apps via helm in kubernetes",
		Tagline:      "",
		PreviewImage: "https://www.some-url.com",
		ViewTime:     22,
	}
	interestUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	suite.mockDraftsRepository.EXPECT().GetDraft(suite.goContext, draftUUID, userUUID).Return(draft, nil).Times(1)
	draft.ConvertInterests()
	suite.mockPostValidator.EXPECT().ValidateAndGetReadTime(draft, suite.goContext).Return(models.MetaData{
		Title:    "Install apps via helm in kubernetes",
		Tagline:  "",
		ReadTime: 22,
	}, nil).Times(1)
	suite.mockTransactionManager.EXPECT().NewTransaction().Return(suite.mockTransaction).Times(1)
	suite.mockPostsRepository.EXPECT().CreatePost(suite.goContext, suite.mockTransaction, post).Return(postUUID, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().GetInterestIDs(suite.goContext, draft.InterestTags).Return(interestUUIDs, nil).Times(1)
	suite.mockPostsRepository.EXPECT().AddInterests(suite.goContext, suite.mockTransaction, postUUID, interestUUIDs).Return(nil).Times(1)
	suite.mockAbstractPostRepository.EXPECT().Save(suite.goContext, suite.mockTransaction, abstractPost).Return(abstractPostUUID, errors.New("something went wrong")).Times(1)
	suite.mockTransaction.EXPECT().Rollback().Return(nil).Times(1)

	err := suite.postService.PublishPost(suite.goContext, draftUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestLikePost_WhenSuccess() {
	postUUID := uuid.New()
	userUUID := uuid.New()
	suite.mockPostsRepository.EXPECT().Like(suite.goContext, postUUID, userUUID).Return(errors.New("something went wrong")).Times(1)

	err := suite.postService.LikePost(suite.goContext, postUUID, userUUID)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *PostServiceTest) TestLikePost_WhenRepositoryReturnsError() {
	postUUID := uuid.New()
	userUUID := uuid.New()
	suite.mockPostsRepository.EXPECT().Like(suite.goContext, postUUID, userUUID).Return(nil).Times(1)

	err := suite.postService.LikePost(suite.goContext, postUUID, userUUID)
	suite.Nil(err)
}
