package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"post-api/dbhelper"
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
	previewPostsRepository PreviewPostsRepository
	dbHelper               helper.DbHelper
}

func (suite *PreviewPostsRepositoryIntegrationTest) SetupTest() {
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
	suite.previewPostsRepository = NewPreviewPostsRepository(db)
	suite.postsRepository = NewPostsRepository(db)
	suite.dbHelper = helper.NewDbHelper(db)
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
	post := db.PublishPost{
		PUID:   "q12we34r",
		UserID: "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		ReadTime:  180,
		ViewCount: 0,
	}
	postID, err := suite.postsRepository.CreatePost(suite.goContext, post)
	suite.Nil(err)

	previewPost := db.PreviewPost{
		PostID:       postID,
		Title:        "Some title",
		Tagline:      "some tagline",
		PreviewImage: "some image url",
		LikeCount:    0,
		CommentCount: 0,
		ViewTime:     0,
	}

	_, err = suite.previewPostsRepository.SavePreview(suite.goContext, previewPost)
	suite.Nil(err)
}

func (suite *PreviewPostsRepositoryIntegrationTest) TestSavePreview_WhenNoPostSavedInPostsTable() {
	previewPost := db.PreviewPost{
		PostID:       4,
		Title:        "Some title",
		Tagline:      "some tagline",
		PreviewImage: "some image url",
		LikeCount:    0,
		CommentCount: 0,
		ViewTime:     0,
	}

	_, err := suite.previewPostsRepository.SavePreview(suite.goContext, previewPost)
	suite.NotNil(err)
}
