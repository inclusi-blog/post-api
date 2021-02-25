package user_profile

import (
	"context"
	"encoding/json"
	"github.com/gola-glitch/gola-utils/golaerror"
	"github.com/gola-glitch/gola-utils/http/request"
	"github.com/gola-glitch/gola-utils/http/request/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"post-api/configuration"
	"testing"
)

type ClientTestSuite struct {
	suite.Suite
	goContext              context.Context
	configData             *configuration.ConfigData
	mockController         *gomock.Controller
	mockGolaHttpRequest    *mocks.MockHttpRequest
	mockGolaRequestBuilder *mocks.MockHttpRequestBuilder
	client                 Client
}

func (suite *ClientTestSuite) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.goContext = context.WithValue(context.TODO(), "testKey", "testVal")
	suite.configData = &configuration.ConfigData{
		UserProfileBaseUrl:         "http://localhost:8084",
		FetchUserFollowedInterests: "/api/user-profile/user/following/interests",
	}
	suite.mockGolaRequestBuilder = mocks.NewMockHttpRequestBuilder(suite.mockController)
	suite.mockGolaHttpRequest = mocks.NewMockHttpRequest(suite.mockController)
	suite.client = NewClient(suite.mockGolaRequestBuilder, suite.configData)
}

func (suite *ClientTestSuite) TearDownTest() {
	suite.mockController.Finish()
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (suite *ClientTestSuite) TestFetchUserFollowingInterests_WhenUserProfileCallSucceeded() {
	expectedUserFollowingInterests := []string{"Poem", "Art"}
	suite.mockGolaRequestBuilder.EXPECT().NewRequestWithContext(suite.goContext).Return(suite.mockGolaHttpRequest).Times(1)
	suite.mockGolaHttpRequest.EXPECT().Get(suite.configData.UserProfileBaseUrl + suite.configData.FetchUserFollowedInterests).Times(1)
	suite.mockGolaHttpRequest.EXPECT().ResponseAs(gomock.Any()).DoAndReturn(func(response interface{}) request.HttpRequest {
		tempResponsePointer := response.(*[]string)
		*tempResponsePointer = expectedUserFollowingInterests
		return suite.mockGolaHttpRequest
	})

	interests, err := suite.client.FetchUserFollowingInterests(suite.goContext)
	suite.Nil(err)
	suite.Equal(expectedUserFollowingInterests, interests)
}

func (suite *ClientTestSuite) TestFetchUserFollowingInterests_WhenUserProfileCallFails() {
	userProfileErr := golaerror.Error{
		ErrorCode:    "ERR_USER_PROFILE_NO_INTERESTS_FOLLOWED",
		ErrorMessage: "no interests followed by user",
	}

	jsonBytes, err := json.Marshal(userProfileErr)
	httpError := golaerror.HttpError{
		StatusCode:   http.StatusNotFound,
		ResponseBody: jsonBytes,
	}
	expectedUserFollowingInterests := []string{"Poem", "Art"}
	suite.mockGolaRequestBuilder.EXPECT().NewRequestWithContext(suite.goContext).Return(suite.mockGolaHttpRequest).Times(1)
	suite.mockGolaHttpRequest.EXPECT().Get(suite.configData.UserProfileBaseUrl + suite.configData.FetchUserFollowedInterests).Return(httpError).Times(1)
	suite.mockGolaHttpRequest.EXPECT().ResponseAs(gomock.Any()).DoAndReturn(func(response interface{}) request.HttpRequest {
		tempResponsePointer := response.(*[]string)
		*tempResponsePointer = expectedUserFollowingInterests
		return suite.mockGolaHttpRequest
	})

	interests, err := suite.client.FetchUserFollowingInterests(suite.goContext)
	suite.Nil(err)
	suite.Nil(interests)
}

func (suite *ClientTestSuite) TestFetchUserFollowingInterests_WhenUserProfileCallFailsWithGenericError() {
	userProfileErr := golaerror.Error{
		ErrorCode:    "ERR_INTERNAL_SERVER_ERROR",
		ErrorMessage: "no interests followed by user",
	}

	jsonBytes, err := json.Marshal(userProfileErr)
	httpError := golaerror.HttpError{
		StatusCode:   http.StatusInternalServerError,
		ResponseBody: jsonBytes,
	}
	expectedUserFollowingInterests := []string{"Poem", "Art"}
	suite.mockGolaRequestBuilder.EXPECT().NewRequestWithContext(suite.goContext).Return(suite.mockGolaHttpRequest).Times(1)
	suite.mockGolaHttpRequest.EXPECT().Get(suite.configData.UserProfileBaseUrl + suite.configData.FetchUserFollowedInterests).Return(httpError).Times(1)
	suite.mockGolaHttpRequest.EXPECT().ResponseAs(gomock.Any()).DoAndReturn(func(response interface{}) request.HttpRequest {
		tempResponsePointer := response.(*[]string)
		*tempResponsePointer = expectedUserFollowingInterests
		return suite.mockGolaHttpRequest
	})

	_, err = suite.client.FetchUserFollowingInterests(suite.goContext)
	suite.NotNil(err)
	suite.Equal(&userProfileErr, err)
}
