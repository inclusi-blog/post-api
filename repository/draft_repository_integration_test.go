package repository

import (
	"context"
	"database/sql"
	"fmt"
	"post-api/dbhelper"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/repository/helper"
	"post-api/service/test_helper"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
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
	database, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintln("Could not connect to test DB", err))
	}
	fmt.Print(database)
	suite.db = database
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	suite.draftRepository = NewDraftRepository(database)
	suite.dbHelper = helper.NewDbHelper(database)
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

	expectedDraft := db.DraftDB{
		DraftID: "abcdef124231",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{"title": "some post data"}`),
		},
		PreviewImage: sql.NullString{String: "https://www.some-url.com", Valid: true},
		Tagline: sql.NullString{
			String: "this is some tagline for draft",
			Valid:  true,
		},
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
	suite.Equal(db.DraftDB{}, draft)
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

func (suite *DraftRepositoryIntegrationTest) TestGetAllDraft_WhenDBHasNoValues() {

	req := models.GetAllDraftRequest{UserID: "3", StartValue: 1, Limit: 3}

	res, err := suite.draftRepository.GetAllDraft(suite.goContext, req)

	suite.Nil(err)

	suite.Zero(len(res))

}

func (suite *DraftRepositoryIntegrationTest) TestGetAllDraft_WhenDBHasValues() {

	getAllDraftRequest := models.GetAllDraftRequest{UserID: "3", StartValue: 0, Limit: 3}

	draft := models.UpsertDraft{DraftID: "11", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}

	expectedDraft := []db.Draft{{DraftID: draft.DraftID, UserID: draft.UserID, PostData: draft.PostData, PreviewImage: nil, Tagline: nil, Interest: models.JSONString{JSONText: types.JSONText(``)}}}

	err := suite.draftRepository.SavePostDraft(draft, suite.goContext)
	suite.Nil(err)

	getAllDraftResponse, err := suite.draftRepository.GetAllDraft(suite.goContext, getAllDraftRequest)

	suite.Nil(err)

	suite.Equal(expectedDraft, getAllDraftResponse)

}

func (suite *DraftRepositoryIntegrationTest) TestGetAllDraft_WhenReturnsMultipleValuesUsingPagination() {

	getAllDraftRequest := models.GetAllDraftRequest{UserID: "3", StartValue: 0, Limit: 3}

	draft1 := models.UpsertDraft{DraftID: "11", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}
	draft2 := models.UpsertDraft{DraftID: "12", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}
	draft3 := models.UpsertDraft{DraftID: "13", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}
	draft4 := models.UpsertDraft{DraftID: "14", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}

	expectedDraft := []db.Draft{{DraftID: draft4.DraftID, UserID: draft1.UserID, PostData: draft1.PostData, PreviewImage: nil, Tagline: nil, Interest: models.JSONString{JSONText: types.JSONText(``)}},
		{DraftID: draft3.DraftID, UserID: draft3.UserID, PostData: draft3.PostData, PreviewImage: nil, Tagline: nil, Interest: models.JSONString{JSONText: types.JSONText(``)}},
		{DraftID: draft2.DraftID, UserID: draft2.UserID, PostData: draft2.PostData, PreviewImage: nil, Tagline: nil, Interest: models.JSONString{JSONText: types.JSONText(``)}}}

	draft1Err := suite.draftRepository.SavePostDraft(draft1, suite.goContext)
	suite.Nil(draft1Err)

	draft2Err := suite.draftRepository.SavePostDraft(draft2, suite.goContext)
	suite.Nil(draft2Err)

	draft3Err := suite.draftRepository.SavePostDraft(draft3, suite.goContext)
	suite.Nil(draft3Err)

	draft4Err := suite.draftRepository.SavePostDraft(draft3, suite.goContext)
	suite.Nil(draft4Err)

	err := suite.draftRepository.SavePostDraft(draft4, suite.goContext)
	suite.Nil(err)

	getAllDraftResponse, err := suite.draftRepository.GetAllDraft(suite.goContext, getAllDraftRequest)

	suite.Nil(err)

	suite.Equal(expectedDraft, getAllDraftResponse)

}

func (suite *DraftRepositoryIntegrationTest) TestGetAllDraft_WhenReturnsMultipleValuesForInBetweenPages() {

	getAllDraftRequest := models.GetAllDraftRequest{UserID: "3", StartValue: 1, Limit: 3}

	draft1 := models.UpsertDraft{DraftID: "11", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}
	draft2 := models.UpsertDraft{DraftID: "12", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}
	draft3 := models.UpsertDraft{DraftID: "13", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}
	draft4 := models.UpsertDraft{DraftID: "14", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}

	expectedDraft := []db.Draft{{DraftID: draft3.DraftID, UserID: draft3.UserID, PostData: draft3.PostData, PreviewImage: nil, Tagline: nil, Interest: models.JSONString{JSONText: types.JSONText(``)}},
		{DraftID: draft2.DraftID, UserID: draft2.UserID, PostData: draft2.PostData, PreviewImage: nil, Tagline: nil, Interest: models.JSONString{JSONText: types.JSONText(``)}},
		{DraftID: draft1.DraftID, UserID: draft1.UserID, PostData: draft1.PostData, PreviewImage: nil, Tagline: nil, Interest: models.JSONString{JSONText: types.JSONText(``)}}}

	draft1Err := suite.draftRepository.SavePostDraft(draft1, suite.goContext)
	suite.Nil(draft1Err)

	draft2Err := suite.draftRepository.SavePostDraft(draft2, suite.goContext)
	suite.Nil(draft2Err)

	draft3Err := suite.draftRepository.SavePostDraft(draft3, suite.goContext)
	suite.Nil(draft3Err)

	draft4Err := suite.draftRepository.SavePostDraft(draft3, suite.goContext)
	suite.Nil(draft4Err)

	err := suite.draftRepository.SavePostDraft(draft4, suite.goContext)
	suite.Nil(err)

	getAllDraftResponse, err := suite.draftRepository.GetAllDraft(suite.goContext, getAllDraftRequest)

	suite.Nil(err)

	suite.Equal(expectedDraft, getAllDraftResponse)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveInterestsToDraft_WhenDraftAvailable() {

	draft := models.UpsertDraft{DraftID: "11", UserID: "3", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}

	draftErr := suite.draftRepository.SavePostDraft(draft, suite.goContext)
	suite.Nil(draftErr)

	interstRequest := request.InterestsSaveRequest{UserID: "3", DraftID: "11", Interests: models.JSONString{JSONText: types.JSONText(``)}}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interstRequest, suite.goContext)
	suite.Nil(interestErr)

}

func (suite *DraftRepositoryIntegrationTest) TestSaveInterestsToDraft_WhenDraftNotAvailable() {

	interstRequest := request.InterestsSaveRequest{UserID: "3", DraftID: "13", Interests: models.JSONString{JSONText: types.JSONText(``)}}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interstRequest, suite.goContext)
	suite.Nil(interestErr)

}

func (suite *DraftRepositoryIntegrationTest) TestDeleteDraft_WhenDbDeletesDraft() {
	draft := models.UpsertDraft{
		DraftID: "q2w3e4r5t6y1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestDataMoreThan100Len),
		},
	}
	err := suite.draftRepository.SavePostDraft(draft, suite.goContext)
	suite.Nil(err)

	err = suite.draftRepository.DeleteDraft(suite.goContext, "q2w3e4r5t6y1")
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestDeleteDraft_WhenNoDraftFound() {
	err := suite.draftRepository.DeleteDraft(suite.goContext, "q2w3e4r5t6y1")
	suite.NotNil(err)
	suite.Equal(sql.ErrNoRows, err)
}
