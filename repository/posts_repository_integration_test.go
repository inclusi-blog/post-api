package repository

import (
	"context"
	"errors"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"post-api/constants"
	"post-api/dbhelper"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/response"
	"post-api/repository/helper"
	"testing"
)

type PostsRepositoryIntegrationTest struct {
	suite.Suite
	db              neo4j.Session
	driver          neo4j.Driver
	adminDb         neo4j.Session
	adminDriver     neo4j.Driver
	goContext       context.Context
	postsRepository PostsRepository
	draftRepository DraftRepository
	dbHelper        helper.DbHelper
}

func (suite *PostsRepositoryIntegrationTest) SetupTest() {
	err := godotenv.Load("../docker-compose-test.env")
	suite.Nil(err)
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	logger := logging.GetLogger(context.Background())
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }
	suite.driver, err = neo4j.NewDriver(dbhelper.BuildConnectionString(), neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), ""), configForNeo4j40)
	suite.Nil(err)
	suite.adminDriver, err = neo4j.NewDriver(dbhelper.BuildConnectionString(), neo4j.BasicAuth(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"), ""), configForNeo4j40)
	suite.Nil(err)
	suite.NotNil(suite.adminDriver)
	suite.NotNil(suite.driver)

	logger.Info("logging")
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: os.Getenv("DB_SERVICE_NAME")}
	suite.db, err = suite.driver.NewSession(sessionConfig)
	suite.Nil(err)
	suite.adminDb, err = suite.adminDriver.NewSession(sessionConfig)
	suite.Nil(err)

	suite.postsRepository = NewPostsRepository(suite.db)
	suite.draftRepository = NewDraftRepository(suite.db)
	suite.dbHelper = helper.NewDbHelper(suite.adminDb)
	suite.insertInterestEntries()
	suite.createSampleUser("some-user")
}

func (suite *PostsRepositoryIntegrationTest) TearDownTest() {
	suite.ClearPostsData()
	err := suite.driver.Close()
	suite.Nil(err)
	err = suite.adminDriver.Close()
	suite.Nil(err)
	err = suite.db.Close()
	suite.Nil(err)
	err = suite.adminDb.Close()
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) ClearPostsData() {
	e := suite.dbHelper.ClearAll()
	if e != nil {
		assert.Error(suite.T(), e)
	}
}

func TestPostsRepositoryIntegrationTest(t *testing.T) {
	suite.Run(t, new(PostsRepositoryIntegrationTest))
}

func (suite *PostsRepositoryIntegrationTest) TestCreatePost_WhenSuccessfullyStoredInDB() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
		PostUrl:      "this-is-some-url",
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		suite.Nil(err)
		return nil, err
	})

	suite.Nil(result)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCreatePost_WhenSamePostIsSavedDBShouldReturnError() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
		PostUrl:      "this-is-some-url",
	}
	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		suite.Nil(err)
		err = transaction.Commit()
		suite.Nil(err)
		err = suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		return nil, err
	})
	suite.Nil(result)
	suite.NotNil(err)

}

func (suite *PostsRepositoryIntegrationTest) TestLikePost_WhenUserLikes() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		if err != nil {
			err := transaction.Commit()
			suite.Nil(err)
		}
		return nil, err
	})
	suite.Nil(result)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestLikePost_WhenThereIsNoPost() {
	err := suite.postsRepository.LikePost("1q2w3e4r5t6y", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q2w3e4r5t6y")
	suite.NotNil(err)
	suite.Equal("no results found", err.Error())
	suite.False(isLiked)
}

func (suite *PostsRepositoryIntegrationTest) TestIsPostLikedByPerson_WhenThereArePostAndLikes() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		if err != nil {
			err := transaction.Commit()
			suite.Nil(err)
		}
		return nil, err
	})
	suite.Nil(result)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestIsPostLikedByPerson_WhenThereIsNoPostAndLikes() {
	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q2w3e4r5t6y")
	suite.NotNil(err)
	suite.Equal("no results found", err.Error())
	suite.False(isLiked)
}

func (suite *PostsRepositoryIntegrationTest) TestUnlikePost_WhenUserLikes() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		if err != nil {
			err := transaction.Commit()
			suite.Nil(err)
		}
		return nil, err
	})
	suite.Nil(result)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.UnlikePost(suite.goContext, "some-user", "1q323e4r4r43")
	suite.Nil(err)

	isLiked, err = suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.False(isLiked)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestUnlikePost_WhenThereIsNoPost() {
	err := suite.postsRepository.UnlikePost(suite.goContext, "some-user", "1q2w3e4r5t6y")
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q2w3e4r5t6y")
	suite.NotNil(err)
	suite.Equal("no results found", err.Error())
	suite.False(isLiked)
}

func (suite *PostsRepositoryIntegrationTest) TestCommentPost_WhenUserCommentOnAPost() {
	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		if err != nil {
			err := transaction.Commit()
			suite.Nil(err)
		}
		return nil, err
	})
	suite.Nil(result)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.CommentPost(suite.goContext, "third-user", "this is awesome", "1q323e4r4r43")
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCommentPost_WhenThereIsNoPost() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		if err != nil {
			err := transaction.Commit()
			suite.Nil(err)
		}
		return nil, err
	})
	suite.Nil(result)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.CommentPost(suite.goContext, "third-user", "this is awesome", "1q323e4r4r43")
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCommentPost_WhenThereIsUserButNoPost() {
	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")

	err := suite.postsRepository.CommentPost(suite.goContext, "third-user", "this is awesome", "1q323e4r4r43")
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCommentPost_WhenThereIsNoUserAndNoPost() {
	err := suite.postsRepository.CommentPost(suite.goContext, "third-user", "this is awesome", "1q323e4r4r43")
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestGetLikesCountByPostID_WhenThereIsAPost() {
	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		if err != nil {
			err := transaction.Commit()
			suite.Nil(err)
		}
		return nil, err
	})
	suite.Nil(result)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "second-user", suite.goContext)
	suite.Nil(err)

	isLiked, err = suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	likesCount, err := suite.postsRepository.GetLikesCountByPostID(suite.goContext, "1q323e4r4r43")
	suite.Nil(err)
	suite.Equal(int64(2), likesCount)
}

func (suite *PostsRepositoryIntegrationTest) TestGetLikesCountByPostID_WhenThereIsNoPost() {
	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	err := suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "second-user", suite.goContext)
	suite.Nil(err)

	likesCount, err := suite.postsRepository.GetLikesCountByPostID(suite.goContext, "1q323e4r4r43")
	suite.Nil(err)
	suite.Equal(int64(0), likesCount)
}

func (suite *PostsRepositoryIntegrationTest) TestGetLikesCountByPostID_WhenThereIsNoUserWhenLikingPost() {
	suite.createSampleUser("third-user")
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		if err != nil {
			err := transaction.Commit()
			suite.Nil(err)
		}
		return nil, err
	})
	suite.Nil(result)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "some-user", suite.goContext)
	suite.Nil(err)

	isLiked, err := suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	err = suite.postsRepository.LikePost("1q323e4r4r43", "second-user", suite.goContext)
	suite.Nil(err)

	isLiked, err = suite.postsRepository.IsPostLikedByPerson(suite.goContext, "some-user", "1q323e4r4r43")
	suite.True(isLiked)
	suite.Nil(err)

	likesCount, err := suite.postsRepository.GetLikesCountByPostID(suite.goContext, "1q323e4r4r43")
	suite.Nil(err)
	suite.Equal(int64(1), likesCount)
}

func (suite *PostsRepositoryIntegrationTest) TestFetchPost_WhenThereIsAPost() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
		PostUrl:      "this-is-some-url",
	}

	expectedDraft := response.Post{
		PostID: "1q323e4r4r43",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		LikeCount:              0,
		CommentCount:           0,
		Interests:              []string{"Art", "Books", "Grammar"},
		AuthorID:               "some-user",
		PreviewImage:           "some-url",
		PublishedAt:            0,
		IsViewerLiked:          false,
		IsViewIsAuthor:         true,
		IsViewerFollowedAuthor: false,
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		suite.Nil(err)
		return nil, nil
	})

	suite.Nil(result)
	suite.Nil(err)
	fetchPost, err := suite.postsRepository.FetchPost(suite.goContext, "1q323e4r4r43", "some-user")
	suite.Nil(err)
	suite.Equal(expectedDraft.PreviewImage, fetchPost.PreviewImage)
	suite.Equal(expectedDraft.PostData, fetchPost.PostData)
	suite.Equal(expectedDraft.IsViewerFollowedAuthor, fetchPost.IsViewerFollowedAuthor)
	suite.Equal(expectedDraft.IsViewIsAuthor, fetchPost.IsViewIsAuthor)
	suite.Equal(expectedDraft.IsViewerLiked, fetchPost.IsViewerLiked)
	suite.Equal(expectedDraft.AuthorID, fetchPost.AuthorID)
	suite.Equal(expectedDraft.LikeCount, fetchPost.LikeCount)
	suite.Equal(expectedDraft.CommentCount, fetchPost.CommentCount)
	suite.ElementsMatch(expectedDraft.Interests, fetchPost.Interests)
	suite.NotEmpty(fetchPost.PublishedAt)
}

func (suite *PostsRepositoryIntegrationTest) TestFetchPost_WhenThereIsAPostAndSomeOneLikedIt() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
		PostUrl:      "this-is-some-url",
	}

	expectedDraft := response.Post{
		PostID: "1q323e4r4r43",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		LikeCount:              2,
		CommentCount:           0,
		Interests:              []string{"Art", "Books", "Grammar"},
		AuthorID:               "some-user",
		PreviewImage:           "some-url",
		PublishedAt:            0,
		IsViewerLiked:          false,
		IsViewIsAuthor:         true,
		IsViewerFollowedAuthor: false,
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		suite.Nil(err)
		return nil, nil
	})

	suite.Nil(result)
	suite.Nil(err)

	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")

	err = suite.postsRepository.LikePost(expectedDraft.PostID, "second-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.LikePost(expectedDraft.PostID, "third-user", suite.goContext)
	suite.Nil(err)

	fetchPost, err := suite.postsRepository.FetchPost(suite.goContext, "1q323e4r4r43", "some-user")
	suite.Nil(err)
	suite.Equal(expectedDraft.PreviewImage, fetchPost.PreviewImage)
	suite.Equal(expectedDraft.PostData, fetchPost.PostData)
	suite.Equal(expectedDraft.IsViewerFollowedAuthor, fetchPost.IsViewerFollowedAuthor)
	suite.Equal(expectedDraft.IsViewIsAuthor, fetchPost.IsViewIsAuthor)
	suite.Equal(expectedDraft.IsViewerLiked, fetchPost.IsViewerLiked)
	suite.Equal(expectedDraft.AuthorID, fetchPost.AuthorID)
	suite.Equal(expectedDraft.LikeCount, fetchPost.LikeCount)
	suite.Equal(expectedDraft.CommentCount, fetchPost.CommentCount)
	suite.ElementsMatch(expectedDraft.Interests, fetchPost.Interests)
	suite.NotEmpty(fetchPost.PublishedAt)
}

func (suite *PostsRepositoryIntegrationTest) TestFetchPost_WhenThereIsAPostAndSomeOneLikedAndCommentedOnIt() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
		PostUrl:      "this-is-some-url",
	}

	expectedDraft := response.Post{
		PostID: "1q323e4r4r43",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		LikeCount:              2,
		CommentCount:           1,
		Interests:              []string{"Art", "Books", "Grammar"},
		AuthorID:               "some-user",
		PreviewImage:           "some-url",
		PublishedAt:            0,
		IsViewerLiked:          false,
		IsViewIsAuthor:         true,
		IsViewerFollowedAuthor: false,
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		suite.Nil(err)
		return nil, nil
	})

	suite.Nil(result)
	suite.Nil(err)

	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	suite.createSampleUser("fourth-user")

	err = suite.postsRepository.LikePost(expectedDraft.PostID, "second-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.LikePost(expectedDraft.PostID, "third-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.CommentPost(suite.goContext, "fourth-user", "this is nice post", expectedDraft.PostID)
	suite.Nil(err)

	fetchPost, err := suite.postsRepository.FetchPost(suite.goContext, "1q323e4r4r43", "some-user")
	suite.Nil(err)
	suite.Equal(expectedDraft.PreviewImage, fetchPost.PreviewImage)
	suite.Equal(expectedDraft.PostData, fetchPost.PostData)
	suite.Equal(expectedDraft.IsViewerFollowedAuthor, fetchPost.IsViewerFollowedAuthor)
	suite.Equal(expectedDraft.IsViewIsAuthor, fetchPost.IsViewIsAuthor)
	suite.Equal(expectedDraft.IsViewerLiked, fetchPost.IsViewerLiked)
	suite.Equal(expectedDraft.AuthorID, fetchPost.AuthorID)
	suite.Equal(expectedDraft.LikeCount, fetchPost.LikeCount)
	suite.Equal(expectedDraft.CommentCount, fetchPost.CommentCount)
	suite.ElementsMatch(expectedDraft.Interests, fetchPost.Interests)
	suite.NotEmpty(fetchPost.PublishedAt)
}

func (suite *PostsRepositoryIntegrationTest) TestFetchPost_WhenThereIsAPostAndDifferentUserViewingThePost() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
		PostUrl:      "this-is-some-url",
	}

	expectedDraft := response.Post{
		PostID: "1q323e4r4r43",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		LikeCount:              2,
		CommentCount:           1,
		Interests:              []string{"Art", "Books", "Grammar"},
		AuthorID:               "some-user",
		PreviewImage:           "some-url",
		PublishedAt:            0,
		IsViewerLiked:          false,
		IsViewIsAuthor:         false,
		IsViewerFollowedAuthor: false,
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		suite.Nil(err)
		return nil, nil
	})

	suite.Nil(result)
	suite.Nil(err)

	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	suite.createSampleUser("fourth-user")

	err = suite.postsRepository.LikePost(expectedDraft.PostID, "second-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.LikePost(expectedDraft.PostID, "third-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.CommentPost(suite.goContext, "fourth-user", "this is nice post", expectedDraft.PostID)
	suite.Nil(err)

	fetchPost, err := suite.postsRepository.FetchPost(suite.goContext, "1q323e4r4r43", "fourth-user")
	suite.Nil(err)
	suite.Equal(expectedDraft.PreviewImage, fetchPost.PreviewImage)
	suite.Equal(expectedDraft.PostData, fetchPost.PostData)
	suite.Equal(expectedDraft.IsViewerFollowedAuthor, fetchPost.IsViewerFollowedAuthor)
	suite.Equal(expectedDraft.IsViewIsAuthor, fetchPost.IsViewIsAuthor)
	suite.Equal(expectedDraft.IsViewerLiked, fetchPost.IsViewerLiked)
	suite.Equal(expectedDraft.AuthorID, fetchPost.AuthorID)
	suite.Equal(expectedDraft.LikeCount, fetchPost.LikeCount)
	suite.Equal(expectedDraft.CommentCount, fetchPost.CommentCount)
	suite.ElementsMatch(expectedDraft.Interests, fetchPost.Interests)
	suite.NotEmpty(fetchPost.PublishedAt)
}

func (suite *PostsRepositoryIntegrationTest) TestFetchPost_WhenThereIsAPostViewerLikedThePost() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:     73,
		Interest:     []string{"Art", "Books", "Grammar"},
		Title:        "this is some title",
		Tagline:      "this is some tagline",
		PreviewImage: "some-url",
		PostUrl:      "this-is-some-url",
	}

	expectedDraft := response.Post{
		PostID: "1q323e4r4r43",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		LikeCount:              2,
		CommentCount:           1,
		Interests:              []string{"Art", "Books", "Grammar"},
		AuthorID:               "some-user",
		PreviewImage:           "some-url",
		PublishedAt:            0,
		IsViewerLiked:          true,
		IsViewIsAuthor:         false,
		IsViewerFollowedAuthor: false,
	}

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.postsRepository.CreatePost(suite.goContext, post, transaction)
		suite.Nil(err)
		return nil, nil
	})

	suite.Nil(result)
	suite.Nil(err)

	suite.createSampleUser("second-user")
	suite.createSampleUser("third-user")
	suite.createSampleUser("fourth-user")

	err = suite.postsRepository.LikePost(expectedDraft.PostID, "second-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.LikePost(expectedDraft.PostID, "third-user", suite.goContext)
	suite.Nil(err)

	err = suite.postsRepository.CommentPost(suite.goContext, "fourth-user", "this is nice post", expectedDraft.PostID)
	suite.Nil(err)

	fetchPost, err := suite.postsRepository.FetchPost(suite.goContext, "1q323e4r4r43", "second-user")
	suite.Nil(err)
	suite.Equal(expectedDraft.PreviewImage, fetchPost.PreviewImage)
	suite.Equal(expectedDraft.PostData, fetchPost.PostData)
	suite.Equal(expectedDraft.IsViewerFollowedAuthor, fetchPost.IsViewerFollowedAuthor)
	suite.Equal(expectedDraft.IsViewIsAuthor, fetchPost.IsViewIsAuthor)
	suite.Equal(expectedDraft.IsViewerLiked, fetchPost.IsViewerLiked)
	suite.Equal(expectedDraft.AuthorID, fetchPost.AuthorID)
	suite.Equal(expectedDraft.LikeCount, fetchPost.LikeCount)
	suite.Equal(expectedDraft.CommentCount, fetchPost.CommentCount)
	suite.ElementsMatch(expectedDraft.Interests, fetchPost.Interests)
	suite.NotEmpty(fetchPost.PublishedAt)
}

func (suite *PostsRepositoryIntegrationTest) TestFetchPost_WhenThereIsNoPost() {
	fetchPost, err := suite.postsRepository.FetchPost(suite.goContext, "1q323e4r4r43", "second-user")
	suite.NotNil(err)
	suite.Equal(errors.New(constants.NoPostFound), err)
	suite.Empty(fetchPost)
}

func (suite *PostsRepositoryIntegrationTest) insertInterestEntries() {
	interests := []string{
		"CREATE (interest:Category:Interest{name: 'Art'})",
		"CREATE (interest:Category:Interest{name: 'Entertainment'})",
		"CREATE (interest:Category:Interest{name: 'Culture'})",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Poem'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Short stories'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Books'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Literature'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Grammar'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Comics'})-[:BELONGS_TO]->(art)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Movies'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Series'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Anime'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Cartoon'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Animation'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Cooking'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Food'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Agriculture'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Festival'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Language'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Philosophy'})-[:BELONGS_TO]->(culture)",
	}
	for _, query := range interests {
		_, err := suite.adminDb.Run(query, nil)
		suite.Nil(err)
	}
}

func (suite *PostsRepositoryIntegrationTest) createSampleUser(userId string) {
	_, err := suite.adminDb.Run("CREATE (person:Person{ userId: $userId})", map[string]interface{}{
		"userId": userId,
	})
	suite.Nil(err)
}
