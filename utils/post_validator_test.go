package utils

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/suite"
	"post-api/configuration"
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
	draft := db.Draft{
		DraftID: "a1v2b31n",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline: "this is some tagline",
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{ "id": "1", "name": "sports"}, {"id": "2", "name": "economy"}, {"id": "3", "name": "poem"}]`),
		},
	}
	titleText, overAllReadTime, err := suite.postValidator.ValidateAndGetReadTime(&draft, suite.goContext)
	suite.Nil(err)
	suite.Equal(34, overAllReadTime)
	suite.Equal("Install apps via helm in kubernetes",titleText)
	suite.Equal("this is some tagline", draft.Tagline)
}

func (suite *PostValidatorTest) TestValidate_InvalidPostData() {
	draft := db.Draft{
		DraftID:  "a1v2b31n",
		UserID:   "1",
		PostData: models.JSONString{},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline: "this is some tagline",
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{ "id": "1", "name": "sports"}, {"id": "2", "name": "economy"}, {"id": "3", "name": "poem"}]`),
		},
	}
	titleText, overAllReadTime, err := suite.postValidator.ValidateAndGetReadTime(&draft, suite.goContext)
	suite.NotNil(err)
	suite.Equal("", titleText)
	suite.Equal("json: cannot unmarshal object into Go value of type []interface {}", err.Error())
	suite.Zero(overAllReadTime)
}

func (suite *PostValidatorTest) TestValidate_InvalidTitleData() {
	draft := db.Draft{
		DraftID: "a1v2b31n",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{},
		Tagline:   "this is some tagline",
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{ "id": "1", "name": "sports"}, {"id": "2", "name": "economy"}, {"id": "3", "name": "poem"}]`),
		},
	}
	titleText, overAllReadTime, err := suite.postValidator.ValidateAndGetReadTime(&draft, suite.goContext)
	suite.NotNil(err)
	suite.Equal("", titleText)
	suite.Equal("json: cannot unmarshal object into Go value of type []interface {}", err.Error())
	suite.Zero(overAllReadTime)
}

func (suite *PostValidatorTest) TestValidate_InvalidInterestData() {
	draft := db.Draft{
		DraftID: "a1v2b31n",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline:  "this is some tagline",
		Interest: models.JSONString{},
	}
	titleText, overAllReadTime, err := suite.postValidator.ValidateAndGetReadTime(&draft, suite.goContext)
	suite.NotNil(err)
	suite.Equal("", titleText)
	suite.Equal("json: cannot unmarshal object into Go value of type []interface {}", err.Error())
	suite.Zero(overAllReadTime)
}

func (suite *PostValidatorTest) TestValidate_IfInterestNameEmpty() {
	draft := db.Draft{
		DraftID: "a1v2b31n",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline: "this is some tagline",
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{ "id": "1", "name": ""}, {"id": "2", "name": "economy"}, {"id": "3", "name": "poem"}]`),
		},
	}
	titleString, overAllReadTime, err := suite.postValidator.ValidateAndGetReadTime(&draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(titleString)
	suite.Equal("interest is invalid", err.Error())
	suite.Zero(overAllReadTime)
}

func (suite *PostValidatorTest) TestValidate_IfReadTimeIsLesserThanConfigTime() {
	suite.configData.ContentReadTimeConfig = map[string]int{
		"poem": 50,
	}

	draft := db.Draft{
		DraftID: "a1v2b31n",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline: "this is some tagline",
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{ "id": "1", "name": "sports"}, {"id": "2", "name": "economy"}, {"id": "3", "name": "poem"}]`),
		},
	}
	titleText, overAllReadTime, err := suite.postValidator.ValidateAndGetReadTime(&draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(titleText)
	suite.Equal("post interest doesn't meet required read time", err.Error())
	suite.Zero(overAllReadTime)
}

func (suite *PostValidatorTest) TestValidate_IfReadTimeIsLesserThanMinimumConfigTime() {
	suite.configData.ContentReadTimeConfig = map[string]int{
		"finance": 50,
	}
	suite.configData.MinimumPostReadTime = 50

	draft := db.Draft{
		DraftID: "a1v2b31n",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline: "this is some tagline",
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{ "id": "1", "name": "sports"}, {"id": "2", "name": "economy"}, {"id": "3", "name": "poem"}]`),
		},
	}
	titleText, overAllReadTime, err := suite.postValidator.ValidateAndGetReadTime(&draft, suite.goContext)
	suite.NotNil(err)
	suite.Empty(titleText)
	suite.Equal("post doesn't meet minimum read time", err.Error())
	suite.Zero(overAllReadTime)
}

func (suite *PostValidatorTest) TestValidate_ValidPostAndNoTagLine() {
	draft := db.Draft{
		DraftID: "a1v2b31n",
		UserID:  "1",
		PostData: models.JSONString{
			JSONText: types.JSONText(test_helper.ContentTestData),
		},
		TitleData: models.JSONString{
			JSONText: types.JSONText(test_helper.TitleTestData),
		},
		Tagline: "",
		Interest: models.JSONString{
			JSONText: types.JSONText(`[{ "id": "1", "name": "sports"}, {"id": "2", "name": "economy"}, {"id": "3", "name": "poem"}]`),
		},
	}
	titleText, overAllReadTime, err := suite.postValidator.ValidateAndGetReadTime(&draft, suite.goContext)
	suite.Nil(err)
	suite.Equal("Install apps via helm in kubernetes", titleText)
	suite.Equal(34, overAllReadTime)
	suite.Equal("இந்தக் கேள்விதான் சென்னைவாசிகள் உட்பட அனைத்து தமிழக மக்களின் மனதிலும் எழுந்துள்ளது. ஒருநாளைக்கு சராச", draft.Tagline)
}
