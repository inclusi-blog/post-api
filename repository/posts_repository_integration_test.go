package repository

import (
	"context"
	"database/sql"
	"fmt"
	"post-api/dbhelper"
	"post-api/models"
	"post-api/models/db"
	"post-api/repository/helper"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PostsRepositoryIntegrationTest struct {
	suite.Suite
	db              *sqlx.DB
	goContext       context.Context
	postsRepository PostsRepository
	dbHelper        helper.DbHelper
}

func (suite *PostsRepositoryIntegrationTest) SetupTest() {
	err := godotenv.Load("../docker-compose-test.env")
	suite.Nil(err)
	connectionString := dbhelper.BuildConnectionString()
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintln("Could not connect to test DB", err))
	}
	fmt.Print(db)
	suite.db = db
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	suite.postsRepository = NewPostsRepository(db)
	suite.dbHelper = helper.NewDbHelper(db)
}

func (suite *PostsRepositoryIntegrationTest) TearDownTest() {
	suite.ClearPostsData()
	_ = suite.db.Close()
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
		UserID: "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	_, err := suite.postsRepository.CreatePost(suite.goContext, post)

	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCreatePost_WhenDBReturnsError() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-invalid-id",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	_, err := suite.postsRepository.CreatePost(suite.goContext, post)

	suite.NotNil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestSaveInitialLike_WhenThereIsAPostAvailable() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	postID, err := suite.postsRepository.CreatePost(suite.goContext, post)

	suite.Nil(err)
	err = suite.postsRepository.SaveInitialLike(suite.goContext, postID)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestSaveInitialLike_WhenThereIsNoPostAvailable() {
	err := suite.postsRepository.SaveInitialLike(suite.goContext, 1)
	suite.NotNil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestGetPostID_WhenThereIsAPost() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	expectedPostID, err := suite.postsRepository.CreatePost(suite.goContext, post)

	suite.Nil(err)
	actualPostID, err := suite.postsRepository.GetPostID(suite.goContext, post.PUID)
	suite.Nil(err)
	suite.Equal(expectedPostID, actualPostID)
}

func (suite *PostsRepositoryIntegrationTest) TestGetPostID_WhenThereIsNoPost() {
	actualPostID, err := suite.postsRepository.GetPostID(suite.goContext, "qw23e5tsa")
	suite.NotNil(err)
	suite.Equal(sql.ErrNoRows, err)
	suite.Zero(actualPostID)
}

func (suite *PostsRepositoryIntegrationTest) TestAppendOrRemoveUserFromLikedByWhenUserLikes() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	postID, err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	saveLikeErr := suite.postsRepository.SaveInitialLike(suite.goContext, postID)
	suite.Nil(saveLikeErr)

	var userID = int64(1)

	appendOrRemoveUserErr := suite.postsRepository.AppendOrRemoveUserFromLikedBy(postID, userID, suite.goContext)
	suite.Nil(appendOrRemoveUserErr)

	likeCount, err := suite.postsRepository.GetLikeCountByPost(suite.goContext, postID)

	suite.Equal(int64(1), likeCount)
}

func (suite *PostsRepositoryIntegrationTest) TestAppendOrRemoveUserFromLikedByWhenUserDisLikes() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	postID, err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	saveLikeErr := suite.postsRepository.SaveInitialLike(suite.goContext, postID)
	suite.Nil(saveLikeErr)

	var userID = int64(1)

	appendOrRemoveUserErr := suite.postsRepository.AppendOrRemoveUserFromLikedBy(postID, userID, suite.goContext)
	suite.Nil(appendOrRemoveUserErr)

	appendOrRemoveUserErr = suite.postsRepository.AppendOrRemoveUserFromLikedBy(postID, userID, suite.goContext)
	suite.Nil(appendOrRemoveUserErr)

	likeCount, err := suite.postsRepository.GetLikeCountByPost(suite.goContext, postID)

	suite.Equal(int64(0), likeCount)
}

func (suite *PostsRepositoryIntegrationTest) TestGetLikeCountByPost_WhenValidPostID() {

	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	postID, err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	saveLikeErr := suite.postsRepository.SaveInitialLike(suite.goContext, postID)
	suite.Nil(saveLikeErr)

	likeCount, err := suite.postsRepository.GetLikeCountByPost(suite.goContext, postID)
	suite.Nil(err)

	suite.Equal(int64(0), likeCount)

}

func (suite *PostsRepositoryIntegrationTest) TestGetLikeCountByPost_WhenInValidPostID() {

	likeCount, err := suite.postsRepository.GetLikeCountByPost(suite.goContext, int64(3))
	suite.Equal(sql.ErrNoRows, err)
	suite.Equal(int64(0), likeCount)

}
