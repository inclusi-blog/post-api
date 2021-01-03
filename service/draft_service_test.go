package service

import (
	"context"
	"errors"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/models/response"
	"post-api/service/test_helper"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
)

type DraftServiceTest struct {
	suite.Suite
	mockController      *gomock.Controller
	goContext           context.Context
	mockValidator       *mocks.MockPostValidator
	mockDraftRepository *mocks.MockDraftRepository
	draftService        DraftService
}

func TestDraftServiceTestSuite(t *testing.T) {
	suite.Run(t, new(DraftServiceTest))
}

func (suite *DraftServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.mockDraftRepository = mocks.NewMockDraftRepository(suite.mockController)
	suite.mockValidator = mocks.NewMockPostValidator(suite.mockController)
	suite.draftService = NewDraftService(suite.mockDraftRepository, suite.mockValidator)
	suite.goContext = context.WithValue(context.Background(), "someKey", "someValue")
}

func (suite *DraftServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsNilForIsDraftPresent() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().IsDraftPresent(suite.goContext, newDraft.DraftID, newDraft.UserID).Return(nil).Times(1)

	suite.mockDraftRepository.EXPECT().UpdateDraft(newDraft, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsErrorForIsDraftPresent() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().IsDraftPresent(suite.goContext, newDraft.DraftID, newDraft.UserID).Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)
	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)
	suite.NotNil(expectedError)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsDraftNotFoundErrorForIsDraftPresent() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().IsDraftPresent(suite.goContext, newDraft.DraftID, newDraft.UserID).Return(errors.New(constants.NoDraftFoundCode)).Times(1)
	suite.mockDraftRepository.EXPECT().CreateNewPostWithData(newDraft, suite.goContext).Return(nil).Times(1)
	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsErrorForCreateNewDraft() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().IsDraftPresent(suite.goContext, newDraft.DraftID, newDraft.UserID).Return(errors.New(constants.NoDraftFoundCode)).Times(1)
	suite.mockDraftRepository.EXPECT().CreateNewPostWithData(newDraft, suite.goContext).Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)
	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)
	suite.NotNil(expectedError)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsErrorForUpdateDraft() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().IsDraftPresent(suite.goContext, newDraft.DraftID, newDraft.UserID).Return(nil).Times(1)
	suite.mockDraftRepository.EXPECT().UpdateDraft(newDraft, suite.goContext).Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)
	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)
	suite.NotNil(expectedError)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), expectedError)
}

func (suite *DraftServiceTest) TestSaveDraft_WhenDraftRepositoryReturnsNilForUpdateDraft() {
	newDraft := models.UpsertDraft{
		DraftID: "someDraftId1",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(`{ "title": "hello" }`),
		},
	}

	suite.mockDraftRepository.EXPECT().IsDraftPresent(suite.goContext, newDraft.DraftID, newDraft.UserID).Return(nil).Times(1)
	suite.mockDraftRepository.EXPECT().UpdateDraft(newDraft, suite.goContext).Return(nil).Times(1)
	expectedError := suite.draftService.SaveDraft(newDraft, suite.goContext)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertTagline_WhenDraftRepositoryReturnsNoError() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "dummy-id",
		Tagline: "this is some tagline",
	}

	suite.mockDraftRepository.EXPECT().SaveTaglineToDraft(saveRequest, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.UpsertTagline(saveRequest, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertTagline_WhenDraftRepositoryReturnsError() {
	saveRequest := request.TaglineSaveRequest{
		UserID:  "1",
		DraftID: "dummy-id",
		Tagline: "this is some tagline",
	}

	suite.mockDraftRepository.EXPECT().SaveTaglineToDraft(saveRequest, suite.goContext).Return(errors.New("something went wrong")).Times(1)

	expectedError := suite.draftService.UpsertTagline(saveRequest, suite.goContext)

	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestUpsertInterests_WhenDraftRepositoryReturnsNoError() {
	saveRequest := request.InterestsSaveRequest{
		Interest: "sports",
		DraftID:  "121212",
		UserID:   "1",
	}

	suite.mockDraftRepository.EXPECT().SaveInterestsToDraft(saveRequest, suite.goContext).Return(nil).Times(1)

	expectedError := suite.draftService.UpsertInterests(saveRequest, suite.goContext)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestUpsertInterests_WhenDraftRepositoryReturnsError() {
	saveRequest := request.InterestsSaveRequest{
		Interest: "sports",
		DraftID:  "121212",
		UserID:   "1",
	}

	suite.mockDraftRepository.EXPECT().SaveInterestsToDraft(saveRequest, suite.goContext).Return(errors.New("something went wrong")).Times(1)

	expectedError := suite.draftService.UpsertInterests(saveRequest, suite.goContext)

	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestDeleteInterest_WhenDraftRepositoryReturnsNoError() {
	saveRequest := request.InterestsSaveRequest{
		Interest: "sports",
		DraftID:  "121212",
		UserID:   "1",
	}

	suite.mockDraftRepository.EXPECT().DeleteInterest(suite.goContext, saveRequest).Return(nil).Times(1)

	expectedError := suite.draftService.DeleteInterest(suite.goContext, saveRequest)

	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestDeleteInterest_WhenDraftRepositoryReturnsError() {
	saveRequest := request.InterestsSaveRequest{
		Interest: "sports",
		DraftID:  "121212",
		UserID:   "1",
	}

	suite.mockDraftRepository.EXPECT().DeleteInterest(suite.goContext, saveRequest).Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)

	expectedError := suite.draftService.DeleteInterest(suite.goContext, saveRequest)

	suite.NotNil(expectedError)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsError() {
	draftID := "121212"
	expectedDraft := db.DraftDB{}
	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID, "some-user").Return(db.DraftDB{}, errors.New(test_helper.ErrSomethingWentWrong)).Times(1)

	draftData, expectedError := suite.draftService.GetDraft(draftID, "some-user", suite.goContext)

	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsNoDraftFoundError() {
	draftID := "121212"
	expectedDraft := db.DraftDB{}
	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID, "some-user").Return(db.DraftDB{}, errors.New(constants.NoDraftFoundCode)).Times(1)

	draftData, expectedError := suite.draftService.GetDraft(draftID, "some-user", suite.goContext)

	suite.Equal(expectedDraft, draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.NoDraftFoundError, expectedError)
}

func (suite *DraftServiceTest) TestGetDraft_WhenDraftRepositoryReturnsNoError() {
	draftID := "121212"
	draftDB := db.DraftDB{
		DraftID:      "1q2w3e4r5t6y",
		UserID:       "some-user",
		PostData:     models.JSONString{},
		PreviewImage: "some-url",
		Tagline:      "some-tagline",
		Interest:     []string{"sports", "economy"},
	}
	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, draftID, "some-user").Return(draftDB, nil).Times(1)

	actualDraft, expectedError := suite.draftService.GetDraft(draftID, "some-user", suite.goContext)

	suite.Equal(draftDB, actualDraft)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsNoError() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: 1,
		Limit:      5,
	}

	tagline := "My first Data"
	thirdTagline := "படிப்படியாக உயர்ந்த எண்ணிக்கை"
	previewImage := "some preview image"

	draftDBS := []db.DraftDB{
		{
			DraftID: "q2w3e4r5u78i",
			UserID:  "12",
			PostData: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			PreviewImage: previewImage,
			Tagline:      tagline,
			Interest:     []string{"sports", "economy"},
		},
		{
			DraftID: "q2w3e4r5u781",
			UserID:  "12",
			PostData: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			PreviewImage: previewImage,
			Tagline:      tagline,
			Interest:     []string{"sports", "economy"},
		},
		{
			DraftID: "q2w3e4r5u782",
			UserID:  "12",
			PostData: models.JSONString{
				JSONText: types.JSONText(test_helper.LargeTextData),
			},
			PreviewImage: "",
			Tagline:      "",
			Interest:     []string{"sports", "economy"},
		},
	}

	expectedDraft := []db.AllDraft{
		{
			DraftID:   "q2w3e4r5u78i",
			TitleData: "தமிழ்நாட்டில் கொரோனா தொற்று பரவத் தொடங்கியபோது நோயாளிகளின் எண்ணிக்கை படிப்படியாக அதிகரித்து வந்தது. ",
			Tagline:   &tagline,
			Interest:  []string{"sports", "economy"},
		},
		{
			DraftID:   "q2w3e4r5u781",
			TitleData: "தமிழ்நாட்டில் கொரோனா தொற்று பரவத் தொடங்கியபோது நோயாளிகளின் எண்ணிக்கை படிப்படியாக அதிகரித்து வந்தது. ",
			Tagline:   &tagline,
			Interest:  []string{"sports", "economy"},
		},
		{
			DraftID:   "q2w3e4r5u782",
			TitleData: "தமிழ்நாட்டில் கொரோனா தொற்று பரவத் தொடங்கியபோது நோயாளிகளின் எண்ணிக்கை படிப்படியாக அதிகரித்து வந்தது. ",
			Tagline:   &thirdTagline,
			Interest:  []string{"sports", "economy"},
		},
	}

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, allDraftReq).Return(draftDBS, nil).Times(1)

	allDraftActual, expectedError := suite.draftService.GetAllDraft(allDraftReq, suite.goContext)
	suite.Equal(*expectedDraft[2].Tagline, *allDraftActual[2].Tagline)
	suite.Nil(expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsError() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: 1,
		Limit:      5,
	}
	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, allDraftReq).Return([]db.DraftDB{}, errors.New("something went wrong")).Times(1)

	draftData, expectedError := suite.draftService.GetAllDraft(allDraftReq, suite.goContext)
	suite.Nil(draftData)
	suite.NotNil(expectedError)
	suite.Equal(&constants.PostServiceFailureError, expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsInvalidData() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: 1,
		Limit:      5,
	}

	tagline := "My first Data"
	previewImage := "some preview image"

	draftDb := []db.DraftDB{
		{
			DraftID: "q2w3e4r5u78i",
			UserID:  "12",
			PostData: models.JSONString{
				JSONText: types.JSONText(`{`),
			},
			PreviewImage: previewImage,
			Tagline:      tagline,
			Interest:     []string{"sports", "economy"},
		},
	}

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, allDraftReq).Return(draftDb, nil).Times(1)

	allDraftActual, expectedError := suite.draftService.GetAllDraft(allDraftReq, suite.goContext)
	suite.Nil(allDraftActual)
	suite.NotNil(expectedError)
	suite.Equal(constants.StoryInternalServerError("unexpected end of JSON input"), expectedError)
}

func (suite *DraftServiceTest) TestGetAllDraft_WhenDraftRepositoryReturnsNoDraftFoundError() {
	allDraftReq := models.GetAllDraftRequest{
		UserID:     "1",
		StartValue: 1,
		Limit:      5,
	}

	suite.mockDraftRepository.EXPECT().GetAllDraft(suite.goContext, allDraftReq).Return([]db.DraftDB{}, errors.New(constants.NoDraftFoundCode)).Times(1)

	actualDrafts, expectedError := suite.draftService.GetAllDraft(allDraftReq, suite.goContext)
	suite.Nil(actualDrafts)
	suite.NotNil(expectedError)
	suite.Equal(&constants.NoDraftFoundError, expectedError)
}

func (suite *DraftServiceTest) TestSavePreviewImage_WhenSuccess() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "1q2w3e4",
		PreviewImageUrl: "https://some-url",
	}

	suite.mockDraftRepository.EXPECT().UpsertPreviewImage(suite.goContext, imageSaveRequest).Return(nil).Times(1)

	err := suite.draftService.SavePreviewImage(imageSaveRequest, suite.goContext)
	suite.Nil(err)
}

func (suite *DraftServiceTest) TestSavePreviewImage_WhenDbReturnsError() {
	imageSaveRequest := request.PreviewImageSaveRequest{
		UserID:          "1",
		DraftID:         "1q2w3e4",
		PreviewImageUrl: "https://some-url",
	}

	suite.mockDraftRepository.EXPECT().UpsertPreviewImage(suite.goContext, imageSaveRequest).Return(errors.New("something went wrong")).Times(1)

	err := suite.draftService.SavePreviewImage(imageSaveRequest, suite.goContext)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *DraftServiceTest) TestDeleteDraft_WhenDraftRepositoryReturnNoError() {
	suite.mockDraftRepository.EXPECT().DeleteDraft(suite.goContext, "qwertyuiop12", "some-user").Return(nil).Times(1)

	err := suite.draftService.DeleteDraft("qwertyuiop12", "some-user", suite.goContext)
	suite.Nil(err)
}

func (suite *DraftServiceTest) TestDeleteDraft_WhenDraftRepositoryReturnsGenericError() {
	suite.mockDraftRepository.EXPECT().DeleteDraft(suite.goContext, "qwertyuiop12", "some-user").Return(errors.New(test_helper.ErrSomethingWentWrong)).Times(1)

	err := suite.draftService.DeleteDraft("qwertyuiop12", "some-user", suite.goContext)
	suite.NotNil(err)
	suite.Equal(constants.StoryInternalServerError(test_helper.ErrSomethingWentWrong), err)
}

func (suite *DraftServiceTest) TestValidateAndGetDraft_WhenDbReturnsNoError() {
	userId := "some-user"
	draftDb := db.DraftDB{
		DraftID: "1q2w3e4r5t6y",
		UserID:  userId,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.LargeTextData),
		},
		PreviewImage: "this is some preview image",
		Tagline:      "this is some tagline",
		Interest:     []string{"sports", "cricket"},
		IsPublished:  false,
		CreatedAt:    98765434567,
	}

	draft := db.Draft{
		DraftID:      "1q2w3e4r5t6y",
		PostData:     draftDb.PostData,
		PreviewImage: &draftDb.PreviewImage,
		Tagline:      &draftDb.Tagline,
		Interest:     draftDb.Interest,
	}

	data := models.MetaData{
		Title:        "இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சரா",
		Tagline:      "this is new tagline",
		ReadTime:     0,
		PreviewImage: "this is some image",
	}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, "1q2w3e4r5t6y", userId).Return(draftDb, nil).Times(1)
	suite.mockValidator.EXPECT().ValidateAndGetMetaData(draft, suite.goContext).Return(data, nil).Times(1)
	expectedPreviewDraft := response.PreviewDraft{
		DraftID:      "1q2w3e4r5t6y",
		Title:        &data.Title,
		Tagline:      &draftDb.Tagline,
		Interest:     draftDb.Interest,
		PreviewImage: &draftDb.PreviewImage,
		AuthorName:   &userId,
	}

	actualPreviewDraft, err := suite.draftService.ValidateAndGetDraft(suite.goContext, "1q2w3e4r5t6y", userId)
	suite.Nil(err)
	suite.Equal(expectedPreviewDraft, actualPreviewDraft)
}

func (suite *DraftServiceTest) TestValidateAndGetDraft_WhenDbReturnsEmptyTagline() {
	userId := "some-user"
	draftDb := db.DraftDB{
		DraftID: "1q2w3e4r5t6y",
		UserID:  userId,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.LargeTextData),
		},
		PreviewImage: "this is some preview image",
		Tagline:      "",
		Interest:     []string{"sports", "cricket"},
		IsPublished:  false,
		CreatedAt:    98765434567,
	}

	draft := db.Draft{
		DraftID:      "1q2w3e4r5t6y",
		PostData:     draftDb.PostData,
		PreviewImage: &draftDb.PreviewImage,
		Tagline:      &draftDb.Tagline,
		Interest:     draftDb.Interest,
	}

	data := models.MetaData{
		Title:        "இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சரா",
		Tagline:      "this is new tagline",
		ReadTime:     0,
		PreviewImage: "this is some image",
	}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, "1q2w3e4r5t6y", userId).Return(draftDb, nil).Times(1)
	suite.mockValidator.EXPECT().ValidateAndGetMetaData(draft, suite.goContext).Return(data, nil).Times(1)
	expectedPreviewDraft := response.PreviewDraft{
		DraftID:      "1q2w3e4r5t6y",
		Title:        &data.Title,
		Tagline:      &data.Tagline,
		Interest:     draftDb.Interest,
		PreviewImage: &draftDb.PreviewImage,
		AuthorName:   &userId,
	}

	actualPreviewDraft, err := suite.draftService.ValidateAndGetDraft(suite.goContext, "1q2w3e4r5t6y", userId)
	suite.Nil(err)
	suite.Equal(expectedPreviewDraft, actualPreviewDraft)
}

func (suite *DraftServiceTest) TestValidateAndGetDraft_WhenDbReturnsEmptyPreviewImageAndTagline() {
	userId := "some-user"
	draftDb := db.DraftDB{
		DraftID: "1q2w3e4r5t6y",
		UserID:  userId,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		PreviewImage: "",
		Tagline:      "",
		Interest:     []string{"sports", "cricket"},
		IsPublished:  false,
		CreatedAt:    98765434567,
	}

	draft := db.Draft{
		DraftID:      "1q2w3e4r5t6y",
		PostData:     draftDb.PostData,
		PreviewImage: &draftDb.PreviewImage,
		Tagline:      &draftDb.Tagline,
		Interest:     draftDb.Interest,
	}

	data := models.MetaData{
		Title:        "இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சரா",
		Tagline:      "this is new tagline",
		ReadTime:     0,
		PreviewImage: "this is some image",
	}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, "1q2w3e4r5t6y", userId).Return(draftDb, nil).Times(1)
	suite.mockValidator.EXPECT().ValidateAndGetMetaData(draft, suite.goContext).Return(data, nil).Times(1)
	expectedPreviewDraft := response.PreviewDraft{
		DraftID:      "1q2w3e4r5t6y",
		Title:        &data.Title,
		Tagline:      &data.Tagline,
		Interest:     draftDb.Interest,
		PreviewImage: &data.PreviewImage,
		AuthorName:   &userId,
	}

	actualPreviewDraft, err := suite.draftService.ValidateAndGetDraft(suite.goContext, "1q2w3e4r5t6y", userId)
	suite.Nil(err)
	suite.Equal(expectedPreviewDraft, actualPreviewDraft)
}

func (suite *DraftServiceTest) TestValidateAndGetDraft_WhenDbReturnsError() {
	userId := "some-user"

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, "1q2w3e4r5t6y", userId).Return(db.DraftDB{}, errors.New("something went wrong")).Times(1)

	actualPreviewDraft, err := suite.draftService.ValidateAndGetDraft(suite.goContext, "1q2w3e4r5t6y", userId)
	suite.NotNil(err)
	suite.Equal(response.PreviewDraft{}, actualPreviewDraft)
	suite.Equal(constants.StoryInternalServerError("something went wrong"), err)
}

func (suite *DraftServiceTest) TestValidateAndGetDraft_WhenValidationReturnsError() {
	userId := "some-user"
	draftDb := db.DraftDB{
		DraftID: "1q2w3e4r5t6y",
		UserID:  userId,
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.LargeTextData),
		},
		PreviewImage: "this is some preview image",
		Tagline:      "this is some tagline",
		Interest:     []string{"sports", "cricket"},
		IsPublished:  false,
		CreatedAt:    98765434567,
	}

	draft := db.Draft{
		DraftID:      "1q2w3e4r5t6y",
		PostData:     draftDb.PostData,
		PreviewImage: &draftDb.PreviewImage,
		Tagline:      &draftDb.Tagline,
		Interest:     draftDb.Interest,
	}

	suite.mockDraftRepository.EXPECT().GetDraft(suite.goContext, "1q2w3e4r5t6y", userId).Return(draftDb, nil).Times(1)
	suite.mockValidator.EXPECT().ValidateAndGetMetaData(draft, suite.goContext).Return(models.MetaData{}, &constants.DraftValidationFailedError).Times(1)
	actualPreviewDraft, err := suite.draftService.ValidateAndGetDraft(suite.goContext, "1q2w3e4r5t6y", userId)
	suite.NotNil(err)
	suite.Equal(response.PreviewDraft{}, actualPreviewDraft)
	suite.Equal(&constants.DraftValidationFailedError, err)
}
