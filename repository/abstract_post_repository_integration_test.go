package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"post-api/dbhelper"
	repoHelper "post-api/helper"
	"post-api/models"
	"post-api/models/db"
	"post-api/repository/helper"
	"post-api/service/test_helper"
	"testing"
)

type PreviewPostsRepositoryIntegrationTest struct {
	suite.Suite
	db                     *sqlx.DB
	goContext              context.Context
	postsRepository        PostsRepository
	abstractPostRepository AbstractPostRepository
	dbHelper               helper.DbHelper
	userHelper             helper.UserRepository
	draftHelper            DraftRepository
	transaction            repoHelper.TransactionManager
}

func (suite *PreviewPostsRepositoryIntegrationTest) SetupTest() {
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
	suite.abstractPostRepository = NewAbstractPostRepository(dbObj)
	suite.postsRepository = NewPostsRepository(dbObj)
	suite.dbHelper = helper.NewDbHelper(dbObj)
	suite.userHelper = helper.NewUserRepository(dbObj)
	suite.draftHelper = NewDraftRepository(dbObj)
	suite.transaction = repoHelper.NewTransactionManager(suite.db)
}

func (suite *PreviewPostsRepositoryIntegrationTest) TearDownTest() {
	suite.ClearDraftData()
	_ = suite.db.Close()
}

func (suite *PreviewPostsRepositoryIntegrationTest) ClearDraftData() {
	e := suite.dbHelper.ClearAll()
	if e != nil {
		assert.Error(suite.T(), e)
	}
}

func TestPreviewPostsRepositoryIntegrationTest(t *testing.T) {
	suite.Run(t, new(PreviewPostsRepositoryIntegrationTest))
}

func (suite *PreviewPostsRepositoryIntegrationTest) TestSavePreview_WhenSuccess() {
	transaction := suite.transaction.NewTransaction()
	_, _, postID := suite.createPost()
	previewPost := db.AbstractPost{
		PostID:       postID,
		Title:        "Some title",
		Tagline:      "some tagline",
		PreviewImage: "some image url",
		ViewTime:     500,
	}

	_, err := suite.abstractPostRepository.Save(suite.goContext, transaction, previewPost)
	_ = transaction.Commit()
	suite.Nil(err)
}

func (suite *PreviewPostsRepositoryIntegrationTest) TestSavePreview_WhenNoPostSavedInPostsTable() {
	transaction := suite.transaction.NewTransaction()
	previewPost := db.AbstractPost{
		PostID:       uuid.New(),
		Title:        "Some title",
		Tagline:      "some tagline",
		PreviewImage: "some image url",
		ViewTime:     50,
	}

	_, err := suite.abstractPostRepository.Save(suite.goContext, transaction, previewPost)
	suite.NotNil(err)
}

func (suite *PreviewPostsRepositoryIntegrationTest) createPost() (uuid.UUID, uuid.UUID, uuid.UUID) {
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
