package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"post-api/dbhelper"
	"post-api/models"
	"post-api/models/db"
	"post-api/repository/helper"
	"testing"
)

type PostsRepositoryIntegrationTest struct {
	suite.Suite
	db              *sqlx.DB
	goContext       context.Context
	postsRepository PostsRepository
	dbHelper        helper.DbHelper
}

func (suite *PostsRepositoryIntegrationTest) SetupTest() {
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
		TitleData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"Install apps via helm in kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)

	suite.Nil(err)
}

func (suite *PostsRepositoryIntegrationTest) TestCreatePost_WhenDBReturnsError() {
	post := db.PublishPost{
		PUID:   "1q323e4r4r43",
		UserID: "some-invalid-id",
		PostData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"You can use helm to deploy your apps via kubernetes"}]}]`),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(`[{"children":[{"text":"Install apps via helm in kubernetes"}]}]`),
		},
		ReadTime:  73,
		ViewCount: 0,
	}

	err := suite.postsRepository.CreatePost(suite.goContext, post)

	suite.NotNil(err)
}
