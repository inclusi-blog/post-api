package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"post-api/dbhelper"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/repository/helper"
	"testing"
)

type DraftRepositoryIntegrationTest struct {
	suite.Suite
	db              *sqlx.DB
	goContext       context.Context
	draftRepository DraftRepository
	dbHelper        helper.DbHelper
}

func (suite *DraftRepositoryIntegrationTest) SetupTest() {
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
	suite.draftRepository = NewDraftRepository(db)
	suite.dbHelper = helper.NewDbHelper(db)
}

func (suite *DraftRepositoryIntegrationTest) TearDownTest() {
	suite.ClearDraftData()
	_ = suite.db.Close()
}

func (suite *DraftRepositoryIntegrationTest) ClearDraftData() {
	e := suite.dbHelper.ClearAll()
	if e != nil {
		assert.Error(suite.T(), e)
	}
}

func TestDraftRepositoryIntegrationTest(t *testing.T) {
	suite.Run(t, new(DraftRepositoryIntegrationTest))
}

func (suite *DraftRepositoryIntegrationTest) TestSavePostDraft_WhenNewDraft() {
	newDraft := models.UpsertDraft{
		DraftID: "abcdef124231",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
	}

	err := suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSavePostDraft_WhenUpsertPost() {
	newDraft := models.UpsertDraft{
		DraftID: "abcdef124231",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
	}

	err := suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.Nil(err)
	newDraft.PostData = models.JSONString{
		JSONText: types.JSONText(`{}`),
	}
	err = suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSavePostDraft_ShouldReturnErrorWhenUserIDIsString() {
	newDraft := models.UpsertDraft{
		DraftID: "abcdef124231",
		UserID:  "1hb12kb12",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
	}

	err := suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.NotNil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveDraftTagline_WhenDbReturnsSuccess() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "1b2b2b23h",
		Tagline: "this is some tagline that will be stored",
	}

	err := suite.draftRepository.SaveTaglineToDraft(saveRequest, suite.goContext)

	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveDraftTagline_WhenUpsertDbReturnsSuccess() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "1b2b2b23h",
		Tagline: "this is some tagline that will be stored",
	}

	err := suite.draftRepository.SaveTaglineToDraft(saveRequest, suite.goContext)

	suite.Nil(err)

	upsertRequest := request.TaglineSaveRequest{
		UserID:  saveRequest.UserID,
		DraftID: saveRequest.DraftID,
		Tagline: "revereted",
	}

	err = suite.draftRepository.SaveTaglineToDraft(upsertRequest, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveDraftTagline_WhenDbReturnsReturns() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "some-user-id",
		DraftID: "1b2b2b23h",
		Tagline: "this is some tagline that will be stored",
	}

	err := suite.draftRepository.SaveTaglineToDraft(saveRequest, suite.goContext)

	suite.NotNil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestGetDraft_WhenDbReturnsDraft() {
	newDraft := models.UpsertDraft{
		DraftID: "abcdef124231",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title": "some post data"}`),
		},
	}

	taglineSaveRequest := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "abcdef124231",
		Tagline: "this is some tagline for draft",
	}

	err := suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.Nil(err)

	err = suite.draftRepository.SaveTaglineToDraft(taglineSaveRequest, suite.goContext)
	suite.Nil(err)

	err = suite.draftRepository.SaveInterestsToDraft(request.InterestsSaveRequest{
		UserID:  "1",
		DraftID: "abcdef124231",
		Interests: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}, suite.goContext)

	err = suite.draftRepository.UpsertPreviewImage(suite.goContext, request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "abcdef124231",
		PreviewImageUrl: "https://www.some-url.com",
	})

	suite.Nil(err)

	expectedDraft := db.Draft{
		DraftID: "abcdef124231",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title": "some post data"}`),
		},
		PreviewImage: "https://www.some-url.com",
		Tagline:      "this is some tagline for draft",
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{"name":"sports","id":"1"},{"name":"economy","id":"2"}]`),
		},
	}

	draft, err := suite.draftRepository.GetDraft(suite.goContext, "abcdef124231")
	suite.Nil(err)
	suite.Equal(expectedDraft, draft)
}

func (suite *DraftRepositoryIntegrationTest) TestGetDraft_WhenDbReturnsError() {
	draft, err := suite.draftRepository.GetDraft(suite.goContext, "abcdef124231")
	suite.NotNil(err)
	suite.Equal(sql.ErrNoRows, err)
	suite.Equal(db.Draft{}, draft)
}

func (suite *DraftRepositoryIntegrationTest) TestUpsertPreviewImage_WhenUpsertSuccess() {
	newDraft := models.UpsertDraft{
		DraftID: "abcdef124231",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
	}

	previewImageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "abcdef124231",
		PreviewImageUrl: "https://some-url",
	}

	err := suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.Nil(err)

	err = suite.draftRepository.UpsertPreviewImage(suite.goContext, previewImageSaveRequest)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestUpsertPreviewImage_WhenNewDraftSuccess() {
	previewImageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "abcdef124231",
		PreviewImageUrl: "https://some-url",
	}

	err := suite.draftRepository.UpsertPreviewImage(suite.goContext, previewImageSaveRequest)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestUpsertPreviewImage_WhenNewDraftEmpty() {
	previewImageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "hello",
		DraftID:         "abcedsc",
		PreviewImageUrl: "https://some-url",
	}

	err := suite.draftRepository.UpsertPreviewImage(suite.goContext, previewImageSaveRequest)
	suite.NotNil(err)
}
