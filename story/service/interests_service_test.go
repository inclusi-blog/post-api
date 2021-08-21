package service

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"post-api/story/constants"
	"post-api/story/mocks"
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
	expectedData := []string{"Sports","Culture"}
	suite.mockInterestsRepository.EXPECT().GetInterests(suite.goContext).Return(expectedData, nil).Times(1)

	actualInterests, err := suite.interestsService.GetInterests(suite.goContext)

	suite.Nil(err)
	suite.Equal(expectedData, actualInterests)
}

func (suite *InterestsServiceTest) TestGetInterests_WhenDbReturnsError() {
	suite.mockInterestsRepository.EXPECT().GetInterests(suite.goContext).Return(nil, sql.ErrNoRows).Times(1)
	interests, err := suite.interestsService.GetInterests(suite.goContext)
	suite.NotNil(err)
	suite.Equal(&constants.PostServiceFailureError, err)
	suite.Len(interests, 0)
}

func (suite *InterestsServiceTest) TestGetInterests_WhenNoDataReturnedWithNoError() {
	suite.mockInterestsRepository.EXPECT().GetInterests(suite.goContext).Return(nil, nil).Times(1)
	interests, err := suite.interestsService.GetInterests(suite.goContext)
	suite.NotNil(err)
	suite.Equal(&constants.NoInterestsFoundError, err)
	suite.Len(interests, 0)
}
