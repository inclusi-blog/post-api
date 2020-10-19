package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models/db"
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
	expectedInterests := []db.Interest{{
		Name: "some-interest",
	}}
	suite.mockInterestService.EXPECT().GetInterests(suite.context).Return(expectedInterests, nil).Times(1)
	marshal, err := json.Marshal(expectedInterests)

	suite.context.Request, err = http.NewRequest(http.MethodGet, "/api/v1/post/get-interests", nil)
	suite.Nil(err)

	suite.interestsController.GetInterests(suite.context)

	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
	suite.Equal(200, suite.recorder.Code)
}

func (suite *InterestsControllerTest) TestGetInterests_WhenServiceReturnsError() {
	suite.mockInterestService.EXPECT().GetInterests(suite.context).Return(nil, &constants.PostServiceFailureError).Times(1)

	marshal, err := json.Marshal(&constants.PostServiceFailureError)
	suite.context.Request, err = http.NewRequest(http.MethodGet, "/api/v1/post/get-interests", nil)
	suite.Nil(err)
	suite.interestsController.GetInterests(suite.context)

	suite.Equal(500, suite.recorder.Code)
	suite.Equal(string(marshal), string(suite.recorder.Body.Bytes()))
}
