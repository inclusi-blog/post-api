package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models/db"
	"post-api/service/test_helper"
	"testing"
)

type InterestsServiceTest struct {
	suite.Suite
	mockController          *gomock.Controller
	goContext               context.Context
	mockInterestsRepository *mocks.MockInterestsRepository
	interestsService        InterestsService
}

func TestInterestsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InterestsServiceTest))
}

func (suite *InterestsServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.goContext = context.WithValue(context.Background(), "someKey", "someValue")
	suite.mockInterestsRepository = mocks.NewMockInterestsRepository(suite.mockController)
	suite.interestsService = NewInterestsService(suite.mockInterestsRepository)
}

func (suite *InterestsServiceTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *InterestsServiceTest) TestGetInterests_WhenRepositoryReturnsData() {
	expectedData := []db.Interest{
		{
			Name: "some-interests",
		},
	}
	suite.mockInterestsRepository.EXPECT().GetInterests(suite.goContext, "sports", []string{}).Return(expectedData, nil).Times(1)

	actualInterests, err := suite.interestsService.GetInterests(suite.goContext, "sports", []string{})

	suite.Nil(err)
	suite.Equal(expectedData, actualInterests)
}

func (suite *InterestsServiceTest) TestGetInterests_WhenDbReturnsError() {
	suite.mockInterestsRepository.EXPECT().GetInterests(suite.goContext, "sports", []string{}).Return(nil, errors.New(test_helper.ErrSomethingWentWrong)).Times(1)
	interests, err := suite.interestsService.GetInterests(suite.goContext, "sports", []string{})
	suite.NotNil(err)
	suite.Equal(&constants.PostServiceFailureError, err)
	suite.Len(interests, 0)
}

func (suite *InterestsServiceTest) TestGetInterests_WhenNoDataReturnedWithError() {
	suite.mockInterestsRepository.EXPECT().GetInterests(suite.goContext, "sports", []string{}).Return(nil, errors.New(constants.NoInterestsFoundCode)).Times(1)
	interests, err := suite.interestsService.GetInterests(suite.goContext, "sports", []string{})
	suite.NotNil(err)
	suite.Equal(&constants.NoInterestsFoundError, err)
	suite.Len(interests, 0)
}
