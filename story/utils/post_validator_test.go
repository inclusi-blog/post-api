package utils

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
	"post-api/configuration"
	"post-api/story/constants"
	"post-api/story/models"
	"post-api/story/models/db"
	"post-api/story/service/test_helper"
	"testing"
)

type PostValidatorTest struct {
	suite.Suite
	mockController *gomock.Controller
	goContext      context.Context
	configData     *configuration.ConfigData
	postValidator  PostValidator
}

func TestPostValidatorTestSuite(t *testing.T) {
	suite.Run(t, new(PostValidatorTest))
}

func (suite *PostValidatorTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.goContext = context.WithValue(context.Background(), "someKey", "someValue")
	suite.configData = &configuration.ConfigData{
		ContentReadTimeConfig: map[string]int{
			"poem": 22,
		},
		MinimumPostReadTime: 21,
	}
	suite.postValidator = NewPostValidator(suite.configData)
}

func (suite *PostValidatorTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *PostValidatorTest) TestValidate_ValidPost() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	tagline := "this is some tagline"
	interests := "{sports,economy,poem}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:      &tagline,
		Interests:    &interests,
		InterestTags: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetReadTime(draft, suite.goContext)
	suite.Nil(err)
	suite.Equal(32, metaData.ReadTime)
	suite.Equal("இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சராச", metaData.Title)
	suite.Equal(&tagline, draft.Tagline)
}

func (suite *PostValidatorTest) TestValidate_InvalidPostData() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	tagline := "this is some tagline"
	interests := "{sports,economy,poem}"
	draft := db.Draft{
		DraftID:      draftUUID,
		UserID:       userUUID,
		Data:         models.JSONString{},
		Tagline:      &tagline,
		Interests:    &interests,
		InterestTags: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetReadTime(draft, suite.goContext)
	suite.NotNil(err)
	suite.Equal("", metaData.Title)
	suite.Equal(&constants.DraftValidationFailedError, err)
	suite.Zero(metaData.ReadTime)
}

func (suite *PostValidatorTest) TestValidate_IfInterestNameEmpty() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	tagline := "this is some tagline"
	interests := ""
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:      &tagline,
		Interests:    &interests,
		InterestTags: nil,
	}
	metaData, err := suite.postValidator.ValidateAndGetReadTime(draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(metaData.Title)
	suite.Equal(&constants.DraftInterestParseError, err)
	suite.Zero(metaData.ReadTime)
}

func (suite *PostValidatorTest) TestValidate_IfReadTimeIsLesserThanConfigTime() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	interests := "{sports,economy,poem}"
	suite.configData.ContentReadTimeConfig = map[string]int{
		"poem": 50,
	}

	tagline := "this is some tagline"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:      &tagline,
		Interests:    &interests,
		InterestTags: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetReadTime(draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(metaData.Title)
	suite.Equal(&constants.InterestReadTimeDoesNotMeetErr, err)
	suite.Zero(metaData.Title)
}

func (suite *PostValidatorTest) TestValidate_IfReadTimeIsLesserThanMinimumConfigTime() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	interests := "{sports,economy,poem}"
	suite.configData.ContentReadTimeConfig = map[string]int{
		"finance": 50,
	}
	suite.configData.MinimumPostReadTime = 50

	tagline := "this is some tagline"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:      &tagline,
		Interests:    &interests,
		InterestTags: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetReadTime(draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(metaData.Title)
	suite.Equal(&constants.ReadTimeNotMeetError, err)
	suite.Zero(metaData.ReadTime)
}

func (suite *PostValidatorTest) TestValidate_ValidPostAndTagLine() {
	draftUUID := uuid.New()
	userUUID := uuid.New()
	interests := "{sports,economy,poem}"
	draft := db.Draft{
		DraftID: draftUUID,
		UserID:  userUUID,
		Data: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Interests:    &interests,
		InterestTags: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetReadTime(draft, suite.goContext)
	suite.Nil(err)
	suite.Equal("இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சராச", metaData.Title)
	suite.Equal(32, metaData.ReadTime)
	suite.Equal("ராஜஸ்தான்:`காங்கிரஸில் வலுக்கும் மோத", metaData.Tagline)
}
