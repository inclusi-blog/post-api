package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"post-api/dbhelper"
	"post-api/story/models"
	"post-api/story/models/db"
	"post-api/story/models/request"
	"post-api/story/service/test_helper"
	"post-api/test_helper/helper"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
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
	userRepository  helper.UserRepository
	dbHelper        helper.DbHelper
}

func (suite *DraftRepositoryIntegrationTest) SetupTest() {
	err := godotenv.Load("../../docker-compose-test.env")
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
	suite.userRepository = helper.NewUserRepository(suite.db)
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

func (suite *DraftRepositoryIntegrationTest) TestSavePostDraft_WhenUpdate() {
	draftUUID := uuid.New()
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draft := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
		UserID: userUUID,
	}
	draftUUID, err = suite.draftRepository.CreateDraft(suite.goContext, draft)
	suite.Nil(err)
	updatedDraft := models.UpsertDraft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
	}
	err = suite.draftRepository.SavePostDraft(updatedDraft, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSavePostDraft_WhenNoSuchPost() {
	newDraft := models.UpsertDraft{
		DraftID: uuid.New(),
		UserID:  uuid.New(),
		Data: models.JSONString{
			JSONText: types.JSONText(`{"title": "hello"}`),
		},
	}

	err := suite.draftRepository.SavePostDraft(newDraft, suite.goContext)
	suite.NotNil(err)
	expectedErr := errors.New("draft not found")
	suite.Equal(expectedErr, err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveDraftTagline_WhenDbReturnsSuccess() {
	draftUUID := uuid.New()
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draft := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
		UserID: userUUID,
	}
	draftUUID, err = suite.draftRepository.CreateDraft(suite.goContext, draft)
	suite.Nil(err)
	saveRequest := request.TaglineSaveRequest{
		UserID:  userUUID,
		DraftID: draftUUID,
		Tagline: "this is some tagline that will be stored",
	}

	err = suite.draftRepository.SaveTaglineToDraft(saveRequest, suite.goContext)

	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveDraftTagline_WhenNoDraftPresent() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  uuid.New(),
		DraftID: uuid.New(),
		Tagline: "this is some tagline that will be stored",
	}

	err := suite.draftRepository.SaveTaglineToDraft(saveRequest, suite.goContext)

	suite.NotNil(err)
	suite.Equal(errors.New("draft not found"), err)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveInterestsToDraft_WhenDraftAvailable() {
	draftUUID := uuid.New()
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draft := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
		UserID: userUUID,
	}
	draftUUID, err = suite.draftRepository.CreateDraft(suite.goContext, draft)
	suite.Nil(err)

	interests := []string{"Sports", "Culture"}

	interestRequest := request.InterestsSaveRequest{UserID: userUUID, DraftID: draftUUID, Interests: interests}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestRequest, suite.goContext)
	suite.Nil(interestErr)
}

func (suite *DraftRepositoryIntegrationTest) TestSaveInterestsToDraft_WhenDraftNotAvailable() {

	interestsSaveRequest := request.InterestsSaveRequest{UserID: uuid.New(), DraftID: uuid.New(), Interests: []string{"Sports", "Culture"}}

	interestErr := suite.draftRepository.SaveInterestsToDraft(interestsSaveRequest, suite.goContext)
	suite.NotNil(interestErr)
	suite.Equal(errors.New("draft not found"), interestErr)

}

func (suite *DraftRepositoryIntegrationTest) TestUpsertPreviewImage_WhenUpsertSuccess() {
	draftUUID := uuid.New()
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draft := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
		UserID: userUUID,
	}
	draftUUID, err = suite.draftRepository.CreateDraft(suite.goContext, draft)
	suite.Nil(err)

	saveRequest := request.PreviewImageSaveRequest{
		UserID:          userUUID,
		DraftID:         draftUUID,
		PreviewImageUrl: "https://some-url",
	}

	err = suite.draftRepository.UpsertPreviewImage(suite.goContext, saveRequest)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestUpsertPreviewImage_WhenNoDraft() {
	previewImageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          uuid.New(),
		DraftID:         uuid.New(),
		PreviewImageUrl: "https://some-url",
	}

	err := suite.draftRepository.UpsertPreviewImage(suite.goContext, previewImageSaveRequest)
	suite.NotNil(err)
	suite.Equal(errors.New("draft not found"), err)
}

func (suite *DraftRepositoryIntegrationTest) TestGetDraft_WhenDbReturnsDraft() {
	draftUUID := uuid.New()
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draft := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
		UserID: userUUID,
	}
	draftUUID, err = suite.draftRepository.CreateDraft(suite.goContext, draft)
	suite.Nil(err)

	taglineSaveRequest := request.TaglineSaveRequest{
		UserID:  userUUID,
		DraftID: draftUUID,
		Tagline: "this is some tagline for draft",
	}

	err = suite.draftRepository.SaveTaglineToDraft(taglineSaveRequest, suite.goContext)
	suite.Nil(err)

	err = suite.draftRepository.SaveInterestsToDraft(request.InterestsSaveRequest{
		UserID:    userUUID,
		DraftID:   draftUUID,
		Interests: []string{"Culture", "Sports"},
	}, suite.goContext)

	suite.Nil(err)

	err = suite.draftRepository.UpsertPreviewImage(suite.goContext, request.PreviewImageSaveRequest{
		UserID:          userUUID,
		DraftID:         draftUUID,
		PreviewImageUrl: "https://www.some-url.com",
	})
	suite.Nil(err)

	previewImage := "https://www.some-url.com"
	tagLine := "this is some tagline for draft"
	interests := "{Culture,Sports}"
	expectedDraft := db.Draft{
		DraftID:      draftUUID,
		UserID:       userUUID,
		Data:         models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
		PreviewImage: &previewImage,
		Tagline:      &tagLine,
		Interests:    &interests,
	}

	savedDraft, err := suite.draftRepository.GetDraft(suite.goContext, draftUUID, userUUID)
	suite.Nil(err)
	suite.Equal(expectedDraft, savedDraft)
}

func (suite *DraftRepositoryIntegrationTest) TestGetDraft_WhenDbReturnsError() {
	draft, err := suite.draftRepository.GetDraft(suite.goContext, uuid.New(), uuid.New())
	suite.NotNil(err)
	suite.Equal(sql.ErrNoRows, err)
	suite.Equal(db.Draft{}, draft)
}

func (suite *DraftRepositoryIntegrationTest) TestGetAllDraft_WhenDBHasNoDrafts() {
	req := models.GetAllDraftRequest{UserID: uuid.New(), StartValue: 1, Limit: 3}

	res, err := suite.draftRepository.GetAllDraft(suite.goContext, req)

	suite.Nil(err)
	suite.Zero(len(res))
}

func (suite *DraftRepositoryIntegrationTest) TestGetAllDraft_WhenDBHasValues() {
	draftUUID := uuid.New()
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draft := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
		UserID: userUUID,
	}
	draftUUID, err = suite.draftRepository.CreateDraft(suite.goContext, draft)
	suite.Nil(err)
	getAllDraftRequest := models.GetAllDraftRequest{UserID: userUUID, StartValue: 0, Limit: 3}

	now := time.Now()
	expectedDraft := []db.Draft{
		{
			DraftID:      draftUUID,
			UserID:       userUUID,
			Data:         models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
			PreviewImage: nil,
			Tagline:      nil,
			Interests:    nil,
			CreatedAt:    &now,
		},
	}

	actualDrafts, err := suite.draftRepository.GetAllDraft(suite.goContext, getAllDraftRequest)

	for _, draft := range actualDrafts {
		*draft.CreatedAt = now
	}

	suite.Nil(err)
	suite.Equal(expectedDraft, actualDrafts)
}

func (suite *DraftRepositoryIntegrationTest) TestGetAllDraft_WhenReturnsMultipleValuesUsingPagination() {
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draftOne := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userUUID,
	}
	draftTwo := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userUUID,
	}
	draftThree := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userUUID,
	}
	draftFour := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userUUID,
	}

	_, err = suite.draftRepository.CreateDraft(suite.goContext, draftOne)
	suite.Nil(err)
	draftTwoUUID, err := suite.draftRepository.CreateDraft(suite.goContext, draftTwo)
	suite.Nil(err)
	draftThreeUUID, err := suite.draftRepository.CreateDraft(suite.goContext, draftThree)
	suite.Nil(err)
	draftFourUUID, err := suite.draftRepository.CreateDraft(suite.goContext, draftFour)
	suite.Nil(err)

	getAllDraftRequest := models.GetAllDraftRequest{UserID: userUUID, StartValue: 0, Limit: 3}
	now := time.Now()
	expectedDraft := []db.Draft{
		{
			DraftID:      draftFourUUID,
			UserID:       userUUID,
			Data:         models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
			PreviewImage: nil,
			Tagline:      nil,
			Interests:    nil,
			CreatedAt:    &now,
		},
		{
			DraftID:      draftThreeUUID,
			UserID:       userUUID,
			Data:         models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
			PreviewImage: nil,
			Tagline:      nil,
			Interests:    nil,
			CreatedAt:    &now,
		},
		{
			DraftID:      draftTwoUUID,
			UserID:       userUUID,
			Data:         models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
			PreviewImage: nil,
			Tagline:      nil,
			Interests:    nil,
			CreatedAt:    &now,
		},
	}

	actualDrafts, err := suite.draftRepository.GetAllDraft(suite.goContext, getAllDraftRequest)

	for _, draft := range actualDrafts {
		*draft.CreatedAt = now
	}

	suite.Nil(err)
	suite.Equal(expectedDraft, actualDrafts)
}

func (suite *DraftRepositoryIntegrationTest) TestGetAllDraft_WhenReturnsMultipleValuesForInBetweenPages() {
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draftOne := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userUUID,
	}
	draftTwo := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userUUID,
	}
	draftThree := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userUUID,
	}
	draftFour := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
		UserID: userUUID,
	}

	draftOneUUID, err := suite.draftRepository.CreateDraft(suite.goContext, draftOne)
	suite.Nil(err)
	draftTwoUUID, err := suite.draftRepository.CreateDraft(suite.goContext, draftTwo)
	suite.Nil(err)
	draftThreeUUID, err := suite.draftRepository.CreateDraft(suite.goContext, draftThree)
	suite.Nil(err)
	_, err = suite.draftRepository.CreateDraft(suite.goContext, draftFour)
	suite.Nil(err)

	now := time.Now()
	expectedDraft := []db.Draft{
		{
			DraftID:      draftThreeUUID,
			UserID:       userUUID,
			Data:         models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
			PreviewImage: nil,
			Tagline:      nil,
			Interests:    nil,
			CreatedAt:    &now,
		},
		{
			DraftID:      draftTwoUUID,
			UserID:       userUUID,
			Data:         models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
			PreviewImage: nil,
			Tagline:      nil,
			Interests:    nil,
			CreatedAt:    &now,
		},
		{
			DraftID:      draftOneUUID,
			UserID:       userUUID,
			Data:         models.JSONString{JSONText: types.JSONText(test_helper.LargeTextData)},
			PreviewImage: nil,
			Tagline:      nil,
			Interests:    nil,
			CreatedAt:    &now,
		},
	}

	allDraftRequest := models.GetAllDraftRequest{UserID: userUUID, StartValue: 1, Limit: 3}
	actualDrafts, err := suite.draftRepository.GetAllDraft(suite.goContext, allDraftRequest)

	for _, draft := range actualDrafts {
		*draft.CreatedAt = now
	}

	suite.Nil(err)
	suite.Equal(expectedDraft, actualDrafts)
}

func (suite *DraftRepositoryIntegrationTest) TestDeleteDraft_WhenDbDeletesDraft() {
	draftUUID := uuid.New()
	userRequest := helper.CreateUserRequest{
		Email:    "dummyUserOne@gmail.com",
		Role:     "User",
		Password: "some-password",
		Username: "some-username",
	}
	userUUID, err := suite.userRepository.CreateUser(suite.goContext, userRequest)
	suite.Nil(err)
	suite.NotEmpty(userUUID)
	draft := models.CreateDraft{
		Data:   models.JSONString{JSONText: types.JSONText(`[{"title": "some text"}]`)},
		UserID: userUUID,
	}
	draftUUID, err = suite.draftRepository.CreateDraft(suite.goContext, draft)
	suite.Nil(err)

	err = suite.draftRepository.DeleteDraft(suite.goContext, draftUUID, userUUID)
	suite.Nil(err)
}

func (suite *DraftRepositoryIntegrationTest) TestDeleteDraft_WhenNoDraftFound() {
	err := suite.draftRepository.DeleteDraft(suite.goContext, uuid.New(), uuid.New())
	suite.NotNil(err)
	suite.Equal(sql.ErrNoRows, err)
}
