package utils

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
	"post-api/configuration"
	"post-api/constants"
	"post-api/models"
	"post-api/models/db"
	"post-api/service/test_helper"
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
	tagline := "this is some tagline"
	draft := db.Draft{
		DraftID: "a1v2b31n",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:  &tagline,
		Interest: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetMetaData(draft, suite.goContext)
	suite.Nil(err)
	suite.Equal(32, metaData.ReadTime)
	suite.Equal("இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சராச", metaData.Title)
	suite.Equal(&tagline, draft.Tagline)
}

func (suite *PostValidatorTest) TestValidate_InvalidPostData() {
	tagline := "this is some tagline"
	draft := db.Draft{
		DraftID:  "a1v2b31n",
		PostData: models.JSONString{},
		Tagline:  &tagline,
		Interest: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetMetaData(draft, suite.goContext)
	suite.NotNil(err)
	suite.Equal("", metaData.Title)
	suite.Equal(&constants.DraftValidationFailedError, err)
	suite.Zero(metaData.ReadTime)
}

func (suite *PostValidatorTest) TestValidate_InvalidInterestData() {
	tagline := "this is some tagline"
	draft := db.Draft{
		DraftID: "a1v2b31n",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:  &tagline,
		Interest: []string{""},
	}
	metaData, err := suite.postValidator.ValidateAndGetMetaData(draft, suite.goContext)
	suite.NotNil(err)
	suite.Equal("", metaData.Title)
	suite.Equal(&constants.MinimumInterestCountNotMatchErr, err)
	suite.Zero(metaData.ReadTime)
}

func (suite *PostValidatorTest) TestValidate_IfInterestNameEmpty() {
	tagline := "this is some tagline"
	draft := db.Draft{
		DraftID: "a1v2b31n",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:  &tagline,
		Interest: []string{""},
	}
	metaData, err := suite.postValidator.ValidateAndGetMetaData(draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(metaData.Title)
	suite.Equal(&constants.MinimumInterestCountNotMatchErr, err)
	suite.Zero(metaData.ReadTime)
}

func (suite *PostValidatorTest) TestValidate_IfReadTimeIsLesserThanConfigTime() {
	suite.configData.ContentReadTimeConfig = map[string]int{
		"poem": 50,
	}

	tagline := "this is some tagline"
	draft := db.Draft{
		DraftID: "a1v2b31n",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:  &tagline,
		Interest: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetMetaData(draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(metaData.Title)
	suite.Equal(&constants.InterestReadTimeDoesNotMeetErr, err)
	suite.Zero(metaData.Title)
}

func (suite *PostValidatorTest) TestValidate_IfReadTimeIsLesserThanMinimumConfigTime() {
	suite.configData.ContentReadTimeConfig = map[string]int{
		"finance": 50,
	}
	suite.configData.MinimumPostReadTime = 50

	tagline := "this is some tagline"
	draft := db.Draft{
		DraftID: "a1v2b31n",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Tagline:  &tagline,
		Interest: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetMetaData(draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(metaData.Title)
	suite.Equal(&constants.ReadTimeNotMeetError, err)
	suite.Zero(metaData.ReadTime)
}

func (suite *PostValidatorTest) TestValidate_ValidPostAndTagLine() {
	draft := db.Draft{
		DraftID: "a1v2b31n",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		Interest: []string{"sports", "economy", "poem"},
	}
	metaData, err := suite.postValidator.ValidateAndGetMetaData(draft, suite.goContext)
	suite.Nil(err)
	suite.Equal("இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சராச", metaData.Title)
	suite.Equal(32, metaData.ReadTime)
	suite.Equal("ராஜஸ்தான்:`காங்கிரஸில் வலுக்கும் மோத", metaData.Tagline)
}
