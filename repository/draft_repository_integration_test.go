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

func (suite *DraftRepositoryIntegrationTest) TestSaveTitleDraft_WhenNewDraftWithTitle() {
	newDraft := models.UpsertDraft{
		DraftID: "abcdef124231",
		UserID:  "1",
		TitleData: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
		Target: "title",
	}

	err := suite.draftRepository.SaveTitleDraft(newDraft, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveTitleDraft_WhenUpsertPostTitle() {
	newDraft := models.UpsertDraft{
		DraftID: "abcdef124231",
		UserID:  "1",
		TitleData: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
		Target: "title",
	}

	err := suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.Nil(err)
	newDraft.TitleData = models.JSONString{
		JSONText: types.JSONText(`{}`),
	}
	err = suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveTitleDraft_ShouldReturnErrorWhenUserIDIsString() {
	newDraft := models.UpsertDraft{
		DraftID: "abcdef124231",
		UserID:  "1hb12kb12",
		TitleData: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
		Target: "title",
	}

	err := suite.draftRepository.SaveTitleDraft(newDraft, suite.goContext)
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
