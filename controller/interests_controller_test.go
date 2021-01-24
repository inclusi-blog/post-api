package controller

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models/db"
	"post-api/models/request"
	"post-api/models/response"
	"testing"
)

type InterestsControllerTest struct {
	suite.Suite
	mockCtrl            *gomock.Controller
	recorder            *httptest.ResponseRecorder
	context             *gin.Context
	mockInterestService *mocks.MockInterestsService
	interestsController InterestsController
}

func (suite *InterestsControllerTest) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockInterestService = mocks.NewMockInterestsService(suite.mockCtrl)
	suite.recorder = httptest.NewRecorder()
	suite.context, _ = gin.CreateTestContext(suite.recorder)
	suite.interestsController = NewInterestsController(suite.mockInterestService)
}

func (suite *InterestsControllerTest) TearDownTest() {
	suite.mockCtrl.Finish()
}

func TestInterestsControllerTestSuite(t *testing.T) {
	suite.Run(t, new(InterestsControllerTest))
}

func (suite *InterestsControllerTest) TestGetInterests_WhenSuccess() {
	interestSearchRequest := request.SearchInterests{
		SearchKeyword: "sport",
		SelectedTags:  []string{"he"},
	}

	bytesJson, err := json.Marshal(interestSearchRequest)
	suite.Nil(err)

	expectedInterests := []db.Interest{{
		Name: "some-interest",
	}}
	suite.mockInterestService.EXPECT().GetInterests(suite.context, "sport", []string{"he"}).Return(expectedInterests, nil).Times(1)
	marshal, err := json.Marshal(expectedInterests)

	suite.context.Request, err = http.NewRequest(http.MethodGet, "/api/v1/post/get-interests", bytes.NewBufferString(string(bytesJson)))
	suite.Nil(err)

	suite.interestsController.GetInterests(suite.context)

	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
	suite.Equal(200, suite.recorder.Code)
}

func (suite *InterestsControllerTest) TestGetInterests_WhenServiceReturnsError() {
	interestSearchRequest := request.SearchInterests{
		SearchKeyword: "sport",
		SelectedTags:  []string{"he"},
	}

	bytesJson, err := json.Marshal(interestSearchRequest)
	suite.Nil(err)

	suite.mockInterestService.EXPECT().GetInterests(suite.context, "sport", []string{"he"}).Return(nil, &constants.PostServiceFailureError).Times(1)

	marshal, err := json.Marshal(&constants.PostServiceFailureError)
	suite.context.Request, err = http.NewRequest(http.MethodGet, "/api/v1/post/get-interests", bytes.NewBufferString(string(bytesJson)))
	suite.Nil(err)
	suite.interestsController.GetInterests(suite.context)

	suite.Equal(500, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *InterestsControllerTest) TestGetInterests_WhenInvalidRequest() {
	suite.mockInterestService.EXPECT().GetInterests(suite.context, "sport", []string{}).Return([]db.Interest{}, nil).Times(0)

	marshal, err := json.Marshal(&constants.PayloadValidationError)
	suite.context.Request, err = http.NewRequest(http.MethodGet, "/api/v1/post/get-interests", bytes.NewBufferString(`{}`))
	suite.Nil(err)
	suite.interestsController.GetInterests(suite.context)

	suite.Equal(400, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *InterestsControllerTest) TestGetExploreInterests_WhenSuccess() {
	expectedData := []response.CategoryAndInterest{
		{
			Category: "Art",
			Interests: []response.InterestWithIcon{
				{
					Name:             "Comics",
					Image:            "https://upload.wikimedia.org/wikipedia/ta/9/9d/Sirithiran_cartoon_2.JPG",
					IsFollowedByUser: false,
				},
				{
					Name:             "Literature",
					Image:            "https://upload.wikimedia.org/wikipedia/commons/6/69/Ancient_Tamil_Script.jpg",
					IsFollowedByUser: false,
				},
				{
					Name:             "Books",
					Image:            "https://upload.wikimedia.org/wikipedia/commons/a/a0/Book_fair-Tamil_Nadu-35th-Chennai-january-2012-part_30.JPG",
					IsFollowedByUser: false,
				},
				{
					Name:             "Poem",
					Image:            "https://upload.wikimedia.org/wikipedia/commons/f/ff/Subramanya_Bharathi_1960_stamp_of_India.jpg",
					IsFollowedByUser: false,
				},
			},
		},
		{
			Category: "Entertainment",
			Interests: []response.InterestWithIcon{
				{
					Name:             "Anime",
					Image:            "https://www.thenerddaily.com/wp-content/uploads/2018/08/Reasons-To-Watch-Anime.jpg",
					IsFollowedByUser: false,
				},
				{
					Name:             "Series",
					Image:            "https://upload.wikimedia.org/wikipedia/commons/1/10/Meta-image-netflix-symbol-black.png",
					IsFollowedByUser: false,
				},
				{
					Name:             "Movies",
					Image:            "https://upload.wikimedia.org/wikipedia/commons/c/c1/Rajinikanth_and_Vijay_at_the_Nadigar_Sangam_Protest.jpg",
					IsFollowedByUser: false,
				},
			},
		},
		{
			Category: "Culture",
			Interests: []response.InterestWithIcon{
				{
					Name:             "Philosophy",
					Image:            "https://en.wikipedia.org/wiki/Swami_Vivekananda#/media/File:Swami_Vivekananda-1893-09-signed.jpg",
					IsFollowedByUser: false,
				},
				{
					Name:             "Language",
					Image:            "https://upload.wikimedia.org/wikipedia/commons/3/35/Word_Tamil.svg",
					IsFollowedByUser: false,
				},
				{
					Name:             "Festival",
					Image:            "https://images.unsplash.com/photo-1576394435759-02a2674ff6e0",
					IsFollowedByUser: false,
				},
				{
					Name:             "Agriculture",
					Image:            "https://upload.wikimedia.org/wikipedia/commons/f/f1/%281%29_Agriculture_and_rural_farms_of_India.jpg",
					IsFollowedByUser: false,
				},
				{
					Name:             "Cooking",
					Image:            "https://images.unsplash.com/photo-1604740795024-c06eeca4bf89",
					IsFollowedByUser: false,
				},
			},
		},
	}
	suite.mockInterestService.EXPECT().GetExploreCategoriesAndInterests(suite.context).Return(expectedData, nil).Times(1)
	marshal, err := json.Marshal(expectedData)

	suite.context.Request, err = http.NewRequest(http.MethodGet, "/api/v1/post/interest/topics-and-interests", nil)
	suite.Nil(err)

	suite.interestsController.GetExploreInterests(suite.context)

	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
	suite.Equal(200, suite.recorder.Code)
}

func (suite *InterestsControllerTest) TestGetExploreInterests_WhenServiceReturnsError() {
	suite.mockInterestService.EXPECT().GetExploreCategoriesAndInterests(suite.context).Return(nil, &constants.PostServiceFailureError).Times(1)

	marshal, err := json.Marshal(&constants.PostServiceFailureError)
	suite.context.Request, err = http.NewRequest(http.MethodGet, "/api/v1/post/interest/topics-and-interests", nil)
	suite.Nil(err)
	suite.interestsController.GetExploreInterests(suite.context)

	suite.Equal(500, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}

func (suite *InterestsControllerTest) TestGetExploreInterests_WhenServiceReturnsNoCategoriesAndInterestFoundError() {
	suite.mockInterestService.EXPECT().GetExploreCategoriesAndInterests(suite.context).Return(nil, &constants.NoInterestsAndCategoriesErr).Times(1)

	marshal, err := json.Marshal(&constants.NoInterestsAndCategoriesErr)
	suite.context.Request, err = http.NewRequest(http.MethodGet, "/api/v1/post/interest/topics-and-interests", nil)
	suite.Nil(err)
	suite.interestsController.GetExploreInterests(suite.context)

	suite.Equal(404, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}
