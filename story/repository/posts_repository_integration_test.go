package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"post-api/dbhelper"
	repoHelper "post-api/helper"
	"post-api/story/models"
	"post-api/story/models/db"
	"post-api/story/repository/helper"
	"post-api/story/service/test_helper"
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
	draftHelper     DraftRepository
	userHelper      helper.UserRepository
	transaction     repoHelper.TransactionManager
	dbHelper        helper.DbHelper
}

func (suite *PostsRepositoryIntegrationTest) SetupTest() {
	err := godotenv.Load("../docker-compose-test.env")
	suite.Nil(err)
	connectionString := dbhelper.BuildConnectionString()
	dbObj, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintln("Could not connect to test DB", err))
	}
	fmt.Print(dbObj)
	suite.db = dbObj
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	suite.postsRepository = NewPostsRepository(dbObj)
	suite.userHelper = helper.NewUserRepository(suite.db)
	suite.draftHelper = NewDraftRepository(suite.db)
	suite.transaction = repoHelper.NewTransactionManager(suite.db)
	suite.dbHelper = helper.NewDbHelper(dbObj)
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
	_, _, _ = suite.createPost()
}

func (suite *PostsRepositoryIntegrationTest) createPost() (uuid.UUID, uuid.UUID, uuid.UUID) {
	userRequest := helper.CreateUserRequest{
		Email: "dummy@gmail.com",
		Role:  "User",
	}
	userID, err := suite.userHelper.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	draft := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userID,
	}
	draftUUID, err := suite.draftHelper.CreateDraft(suite.goContext, draft)
	suite.Nil(err)
	post := db.PublishPost{
		UserID: userID,
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		DraftID: draftUUID,
	}
	transaction := suite.transaction.NewTransaction()
	postUUID, err := suite.postsRepository.CreatePost(suite.goContext, transaction, post)
	suite.Nil(err)
	err = transaction.Commit()
	suite.Nil(err)
	return userID, draftUUID, postUUID
}

func (suite *PostsRepositoryIntegrationTest) TestCreatePost_WhenDBReturnsError() {
	transaction := suite.transaction.NewTransaction()
	post := db.PublishPost{
		UserID: uuid.New(),
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		DraftID: uuid.New(),
	}

	_, err := suite.postsRepository.CreatePost(suite.goContext, transaction, post)

	suite.NotNil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestLike_WhenUserLikes() {
	userID, _, postID := suite.createPost()
	saveLikeErr := suite.postsRepository.Like(suite.goContext, postID, userID)
	suite.Nil(saveLikeErr)
}

func (suite *PostsRepositoryIntegrationTest) TestLike_WhenUserDisLikes() {
	saveLikeErr := suite.postsRepository.Like(suite.goContext, uuid.New(), uuid.New())
	suite.NotNil(saveLikeErr)
}

func (suite *PostsRepositoryIntegrationTest) TestUnLike_WhenUserLikes() {
	userID, _, postID := suite.createPost()
	saveLikeErr := suite.postsRepository.Like(suite.goContext, postID, userID)
	suite.Nil(saveLikeErr)
	err := suite.postsRepository.UnLike(suite.goContext, postID, userID)
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestUnLike_WhenUserDisLikes() {
	err := suite.postsRepository.UnLike(suite.goContext, uuid.New(), uuid.New())
	suite.NotNil(err)
	suite.Equal("never liked", err.Error())
}

func (suite *PostsRepositoryIntegrationTest) TestAddInterests_WhenThereIsAPost() {
	_, _, postID := suite.createPost()
	interestsIDs, err := suite.dbHelper.GetInterestsIDs([]string{"Sports", "Art", "Entertainment"})
	suite.Nil(err)
	transaction := suite.transaction.NewTransaction()
	err = suite.postsRepository.AddInterests(suite.goContext, transaction, postID, interestsIDs)
	suite.Nil(err)
	err = transaction.Commit()
	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestAddInterests_WhenThereIsANoPost() {
	interestsIDs, err := suite.dbHelper.GetInterestsIDs([]string{"Sports", "Art", "Entertainment"})
	suite.Nil(err)
	transaction := suite.transaction.NewTransaction()
	err = suite.postsRepository.AddInterests(suite.goContext, transaction, uuid.New(), interestsIDs)
	suite.NotNil(err)
}
