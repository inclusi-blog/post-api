package repository

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/jmoiron/sqlx/types"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"os"
	"post-api/constants"
	"post-api/dbhelper"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/repository/helper"
	"post-api/service/test_helper"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DraftRepositoryIntegrationTest struct {
	suite.Suite
	db              neo4j.Session
	driver          neo4j.Driver
	adminDb         neo4j.Session
	adminDriver     neo4j.Driver
	goContext       context.Context
	draftRepository DraftRepository
	dbHelper        helper.DbHelper
}

func (suite *DraftRepositoryIntegrationTest) SetupTest() {
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

	suite.draftRepository = NewDraftRepository(suite.db)
	suite.dbHelper = helper.NewDbHelper(suite.adminDb)
	suite.createSampleUser()
}

func (suite *DraftRepositoryIntegrationTest) TearDownTest() {
	suite.ClearDraftData()
	err := suite.driver.Close()
	suite.Nil(err)
	err = suite.adminDriver.Close()
	suite.Nil(err)
	err = suite.db.Close()
	suite.Nil(err)
	err = suite.adminDb.Close()
	suite.Nil(err)
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

func (suite *DraftRepositoryIntegrationTest) TestIsDraftPresentWhenThereIsNoDraft() {
	err := suite.draftRepository.IsDraftPresent(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.NotNil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestIsDraftPresentWhenThereIsADraft() {
	draft := models.UpsertDraft{
		DraftID: "1q2w3e4r5t6y",
		UserID:  "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.LargeTextData),
		},
	}
	err := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(err)
	err = suite.draftRepository.IsDraftPresent(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestCreateNewPostWithData_WhenDbReturnsNoError() {
	draft := models.UpsertDraft{
		DraftID: "1q2w3e4r5t6y",
		UserID:  "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.LargeTextData),
		},
	}
	err := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestCreateNewPostWithData_WhenAlreadyADraftPresent() {
	draft := models.UpsertDraft{
		DraftID: "1q2w3e4r5t6y",
		UserID:  "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.LargeTextData),
		},
	}
	err := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(err)
	err = suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.NotNil(err)
	suite.True(neo4j.IsClientError(err))
}

func (suite *DraftRepositoryIntegrationTest) TestUpdateDraft_WhenThereIsADraft() {
	draft := models.UpsertDraft{
		DraftID: "1q2w3e4r5t6y",
		UserID:  "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText("[{children: { text: \"some text\"}]"),
		},
	}

	err := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(err)
	draft.PostData = models.JSONString{
		JSONText: types.JSONText(test_helper.LargeTextData),
	}
	err = suite.draftRepository.UpdateDraft(draft, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestUpdateDraft_WhenThereIsNoADraft() {
	draft := models.UpsertDraft{
		DraftID: "1q2w3e4r5t6y",
		UserID:  "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText("[{children: { text: \"some text\"}]"),
		},
	}

	err := suite.draftRepository.UpdateDraft(draft, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveTaglineToDraft_WhenThereIsADraft() {
	draft := models.UpsertDraft{
		DraftID: "1q2w3e4r5t6y",
		UserID:  "some-user",
		PostData: models.JSONString{
			JSONText: types.JSONText("[{children: { text: \"some text\"}]"),
		},
	}
	saveRequest := request.TaglineSaveRequest{
		UserID:  "some-user",
		DraftID: "1q2w3e4r5t6y",
		Tagline: "this is some tagline",
	}
	err := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(err)
	draft.PostData = models.JSONString{
		JSONText: types.JSONText(test_helper.LargeTextData),
	}
	err = suite.draftRepository.SaveTaglineToDraft(saveRequest, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveTaglineToDraft_WhenThereIsNoADraft() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "some-user",
		DraftID: "1q2w3e4r5t6y",
		Tagline: "this is some tagline",
	}

	err := suite.draftRepository.SaveTaglineToDraft(saveRequest, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveInterestsToDraft_WhenDraftAvailable() {
	suite.insertInterestEntries()
	draft := models.UpsertDraft{DraftID: "1q2w3e4r5t6y", UserID: "some-user", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}

	draftErr := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(draftErr)

	interestRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t6y", Interest: "Art"}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveInterestsToDraft_WhenDraftNotAvailable() {
	suite.insertInterestEntries()
	interestRequest := request.InterestsSaveRequest{UserID: "some-user-one", DraftID: "13", Interest: ""}
	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)
}

func (suite *DraftRepositoryIntegrationTest) TestDeleteInterest_WhenDraftAvailable() {
	suite.insertInterestEntries()
	draft := models.UpsertDraft{DraftID: "1q2w3e4r5t6y", UserID: "some-user", PostData: models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}}

	draftErr := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(draftErr)

	interestRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t6y", Interest: "Art"}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)
	interestErr = suite.draftRepository.DeleteInterest(suite.goContext, interestRequest)
	suite.Nil(interestErr)
}

func (suite *DraftRepositoryIntegrationTest) TestDeleteInterest_WhenDraftNotAvailable() {
	suite.insertInterestEntries()
	interestRequest := request.InterestsSaveRequest{UserID: "some-user-one", DraftID: "13", Interest: ""}
	interestErr := suite.draftRepository.DeleteInterest(suite.goContext, interestRequest)
	suite.Nil(interestErr)
}

func (suite *DraftRepositoryIntegrationTest) TestGetDraft_WhenThereIsADraftAvailable() {
	suite.insertInterestEntries()
	postData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	draft := models.UpsertDraft{DraftID: "1q2w3e4r5t6y", UserID: "some-user", PostData: postData}

	draftErr := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(draftErr)

	interestRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t6y", Interest: "Art"}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)

	draftDB, err := suite.draftRepository.GetDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	expectedDraft := db.DraftDB{
		DraftID:      "1q2w3e4r5t6y",
		UserID:       "some-user",
		PostData:     postData,
		PreviewImage: "",
		Tagline:      "",
		Interest:     []string{"Art"},
		IsPublished:  false,
	}
	suite.Nil(err)
	suite.Equal(expectedDraft.PreviewImage, draftDB.PreviewImage)
	suite.Equal(expectedDraft.DraftID, draftDB.DraftID)
	suite.Equal(expectedDraft.Interest, draftDB.Interest)
	suite.Equal(expectedDraft.UserID, draftDB.UserID)
	suite.Equal(expectedDraft.IsPublished, draftDB.IsPublished)
	suite.NotEmpty(draftDB.CreatedAt)
}

func (suite *DraftRepositoryIntegrationTest) TestGetDraft_WhenThereIsNoDraftAvailable() {
	suite.insertInterestEntries()
	draftDB, err := suite.draftRepository.GetDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.NotNil(err)
	suite.Equal(db.DraftDB{}, draftDB)
	suite.Equal(constants.NoDraftFoundCode, err.Error())
}

func (suite *DraftRepositoryIntegrationTest) TestUpsertPreviewImage_WhenThereIsADraft() {
	suite.insertInterestEntries()
	postData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	draft := models.UpsertDraft{DraftID: "1q2w3e4r5t6y", UserID: "some-user", PostData: postData}

	draftErr := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(draftErr)

	interestRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t6y", Interest: "Art"}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)

	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t6y",
		PreviewImageUrl: "this is some url",
	}

	draftErr = suite.draftRepository.UpsertPreviewImage(suite.goContext, imageSaveRequest)
	suite.Nil(draftErr)
	draftDB, err := suite.draftRepository.GetDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	expectedDraft := db.DraftDB{
		DraftID:      "1q2w3e4r5t6y",
		UserID:       "some-user",
		PostData:     postData,
		PreviewImage: "this is some url",
		Tagline:      "",
		Interest:     []string{"Art"},
	}
	suite.Nil(err)
	suite.Equal(expectedDraft.PreviewImage, draftDB.PreviewImage)
	suite.Equal(expectedDraft.DraftID, draftDB.DraftID)
	suite.Equal(expectedDraft.Interest, draftDB.Interest)
	suite.Equal(expectedDraft.UserID, draftDB.UserID)
}

func (suite *DraftRepositoryIntegrationTest) TestUpsertPreviewImage_WhenThereIsNoDraft() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t6y",
		PreviewImageUrl: "this is some url",
	}

	draftErr := suite.draftRepository.UpsertPreviewImage(suite.goContext, imageSaveRequest)
	suite.Nil(draftErr)
	draftDB, err := suite.draftRepository.GetDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.NotNil(err)
	suite.Equal(db.DraftDB{}, draftDB)
	suite.Equal(constants.NoDraftFoundCode, err.Error())
}

func (suite *DraftRepositoryIntegrationTest) TestDeleteDraft_WhenThereIsADraft() {
	suite.insertInterestEntries()
	postData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	draft := models.UpsertDraft{DraftID: "1q2w3e4r5t6y", UserID: "some-user", PostData: postData}

	draftErr := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(draftErr)

	interestRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t6y", Interest: "Art"}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)

	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t6y",
		PreviewImageUrl: "this is some url",
	}

	draftErr = suite.draftRepository.UpsertPreviewImage(suite.goContext, imageSaveRequest)
	suite.Nil(draftErr)
	draftDB, err := suite.draftRepository.GetDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	expectedDraft := db.DraftDB{
		DraftID:      "1q2w3e4r5t6y",
		UserID:       "some-user",
		PostData:     postData,
		PreviewImage: "this is some url",
		Tagline:      "",
		Interest:     []string{"Art"},
	}
	suite.Nil(err)
	suite.Equal(expectedDraft.PreviewImage, draftDB.PreviewImage)
	suite.Equal(expectedDraft.DraftID, draftDB.DraftID)
	suite.Equal(expectedDraft.Interest, draftDB.Interest)
	suite.Equal(expectedDraft.UserID, draftDB.UserID)

	draftErr = suite.draftRepository.DeleteDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.Nil(draftErr)

	draftAfterDeleted, err := suite.draftRepository.GetDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.NotNil(err)
	suite.Equal(constants.NoDraftFoundCode, err.Error())
	suite.Equal(db.DraftDB{}, draftAfterDeleted)
}

func (suite *DraftRepositoryIntegrationTest) TestDeleteDraft_WhenThereIsNoDraft() {
	draftBeforeDelete, err := suite.draftRepository.GetDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.NotNil(err)
	suite.Equal(db.DraftDB{}, draftBeforeDelete)
	suite.Equal(constants.NoDraftFoundCode, err.Error())

	draftErr := suite.draftRepository.DeleteDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.Nil(draftErr)

	draftAfterDeleted, err := suite.draftRepository.GetDraft(suite.goContext, "1q2w3e4r5t6y", "some-user")
	suite.NotNil(err)
	suite.Equal(db.DraftDB{}, draftAfterDeleted)
	suite.Equal(constants.NoDraftFoundCode, err.Error())
}

func (suite *DraftRepositoryIntegrationTest) TestFetchAllDraft_WhenLimitSentShouldReturnDraftInThatRange() {
	suite.insertInterestEntries()
	postData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	draft := models.UpsertDraft{DraftID: "1q2w3e4r5t6y", UserID: "some-user", PostData: postData}

	draftErr := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(draftErr)

	interestRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t6y", Interest: "Art"}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)

	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t6y",
		PreviewImageUrl: "this is some url",
	}

	draftErr = suite.draftRepository.UpsertPreviewImage(suite.goContext, imageSaveRequest)
	suite.Nil(draftErr)

	secondPostData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	secondDraft := models.UpsertDraft{DraftID: "1q2w3e4r5t61", UserID: "some-user", PostData: secondPostData}
	secondDraftInterests := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t61", Interest: "Art"}
	secondDraftPreviewImageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t61",
		PreviewImageUrl: "this is some url",
	}

	_ = suite.draftRepository.CreateNewPostWithData(secondDraft, suite.goContext)
	_ = suite.draftRepository.SaveInterestsToDraft(secondDraftInterests, suite.goContext)
	_ = suite.draftRepository.UpsertPreviewImage(suite.goContext, secondDraftPreviewImageSaveRequest)

	thirdPostData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	thirdDraft := models.UpsertDraft{DraftID: "1q2w3e4r5t62", UserID: "some-user", PostData: thirdPostData}
	thirdDraftInterestSaveRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t62", Interest: "Art"}
	thirdDraftPreviewImageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t62",
		PreviewImageUrl: "this is some url",
	}

	_ = suite.draftRepository.CreateNewPostWithData(thirdDraft, suite.goContext)
	_ = suite.draftRepository.SaveInterestsToDraft(thirdDraftInterestSaveRequest, suite.goContext)
	_ = suite.draftRepository.UpsertPreviewImage(suite.goContext, thirdDraftPreviewImageSaveRequest)
	allDraftRequest := models.GetAllDraftRequest{
		UserID:     "some-user",
		StartValue: 1,
		Limit:      2,
	}
	allDraft, draftErr := suite.draftRepository.GetAllDraft(suite.goContext, allDraftRequest)
	suite.Nil(draftErr)
	suite.Len(allDraft, 2)
	suite.Equal("1q2w3e4r5t61", allDraft[0].DraftID)
	suite.Equal("1q2w3e4r5t6y", allDraft[1].DraftID)
}

func (suite *DraftRepositoryIntegrationTest) TestFetchAllDraft_WhenLimitSentShouldReturnDraftsWithRangeZeroToThree() {
	suite.insertInterestEntries()
	postData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	draft := models.UpsertDraft{DraftID: "1q2w3e4r5t6y", UserID: "some-user", PostData: postData}

	draftErr := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(draftErr)

	interestRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t6y", Interest: "Art"}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)

	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t6y",
		PreviewImageUrl: "this is some url",
	}

	draftErr = suite.draftRepository.UpsertPreviewImage(suite.goContext, imageSaveRequest)
	suite.Nil(draftErr)

	secondPostData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	secondDraft := models.UpsertDraft{DraftID: "1q2w3e4r5t61", UserID: "some-user", PostData: secondPostData}
	secondDraftInterests := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t61", Interest: "Art"}
	secondDraftPreviewImageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t61",
		PreviewImageUrl: "this is some url",
	}

	_ = suite.draftRepository.CreateNewPostWithData(secondDraft, suite.goContext)
	_ = suite.draftRepository.SaveInterestsToDraft(secondDraftInterests, suite.goContext)
	_ = suite.draftRepository.UpsertPreviewImage(suite.goContext, secondDraftPreviewImageSaveRequest)

	thirdPostData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	thirdDraft := models.UpsertDraft{DraftID: "1q2w3e4r5t62", UserID: "some-user", PostData: thirdPostData}
	thirdDraftInterestSaveRequest := request.InterestsSaveRequest{UserID: "some-user", DraftID: "1q2w3e4r5t62", Interest: "Art"}
	thirdDraftPreviewImageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "some-user",
		DraftID:         "1q2w3e4r5t62",
		PreviewImageUrl: "this is some url",
	}

	_ = suite.draftRepository.CreateNewPostWithData(thirdDraft, suite.goContext)
	_ = suite.draftRepository.SaveInterestsToDraft(thirdDraftInterestSaveRequest, suite.goContext)
	_ = suite.draftRepository.UpsertPreviewImage(suite.goContext, thirdDraftPreviewImageSaveRequest)
	allDraftRequest := models.GetAllDraftRequest{
		UserID:     "some-user",
		StartValue: 0,
		Limit:      3,
	}
	allDraft, draftErr := suite.draftRepository.GetAllDraft(suite.goContext, allDraftRequest)
	suite.Nil(draftErr)
	suite.Len(allDraft, 3)
}

func (suite *DraftRepositoryIntegrationTest) TestFetchAllDraft_WhenThereIsNoDraft() {
	suite.insertInterestEntries()
	allDraftRequest := models.GetAllDraftRequest{
		UserID:     "some-user",
		StartValue: 0,
		Limit:      3,
	}
	allDraft, draftErr := suite.draftRepository.GetAllDraft(suite.goContext, allDraftRequest)
	suite.NotNil(draftErr)
	suite.Len(allDraft, 0)
	suite.Equal(constants.NoDraftFoundCode, draftErr.Error())
}

func (suite *DraftRepositoryIntegrationTest) TestUpdatePublishedStatus_WhenThereIsADraft() {
	suite.insertInterestEntries()
	postData := models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)}
	draft := models.UpsertDraft{DraftID: "1q2w3e4r5t6y", UserID: "some-user", PostData: postData}

	draftErr := suite.draftRepository.CreateNewPostWithData(draft, suite.goContext)
	suite.Nil(draftErr)

	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.draftRepository.UpdatePublishedStatus(suite.goContext, draft.DraftID, draft.UserID, transaction)
		suite.Nil(err)
		return nil, err
	})
	suite.Nil(err)
	suite.Nil(result)
}

func (suite *DraftRepositoryIntegrationTest) TestUpdatePublishedStatus_WhenThereIsNoDraft() {
	suite.insertInterestEntries()
	result, err := suite.db.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		err := suite.draftRepository.UpdatePublishedStatus(suite.goContext, "1q2w3e4r5t6y", "some-user", transaction)
		suite.NotNil(err)
		err = transaction.Close()
		suite.Nil(err)
		return nil, err
	})
	suite.NotNil(err)
	suite.Nil(result)
}

func (suite *DraftRepositoryIntegrationTest) insertInterestEntries() {
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

func (suite *DraftRepositoryIntegrationTest) createSampleUser() {
	_, err := suite.adminDb.Run("CREATE (person:Person{ userId: $userId})", map[string]interface{}{
		"userId": "some-user",
	})
	suite.Nil(err)
}
