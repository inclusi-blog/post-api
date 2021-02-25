package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"post-api/constants"
	"post-api/mocks"
	"post-api/models/db"
	"post-api/models/response"
	"post-api/service/test_helper"
	"testing"
)

type InterestsServiceTest struct {
	suite.Suite
	mockController          *gomock.Controller
	goContext               context.Context
	mockInterestsRepository *mocks.MockInterestsRepository
	mockUserProfileClient   *mocks.MockClient
	mockInterestsMapper     *mocks.MockInterestsMapper
	interestsService        InterestsService
}

func TestInterestsServiceTestSuite(t *testing.T) {
	suite.Run(t, new(InterestsServiceTest))
}

func (suite *InterestsServiceTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.goContext = context.WithValue(context.Background(), "someKey", "someValue")
	suite.mockInterestsRepository = mocks.NewMockInterestsRepository(suite.mockController)
	suite.mockInterestsMapper = mocks.NewMockInterestsMapper(suite.mockController)
	suite.mockUserProfileClient = mocks.NewMockClient(suite.mockController)
	suite.interestsService = NewInterestsService(suite.mockInterestsRepository, suite.mockUserProfileClient, suite.mockInterestsMapper)
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

func (suite *InterestsServiceTest) TestGetExploreCategoriesAndInterests_WhenRepositoryReturnsData() {
	dbReturnedCategoriesAndInterest := []db.CategoryAndInterest{
		{
			Category: "Art",
			Interests: []db.InterestWithIcon{
				{
					Name:  "Comics",
					Image: "https://upload.wikimedia.org/wikipedia/ta/9/9d/Sirithiran_cartoon_2.JPG",
				},
				{
					Name:  "Literature",
					Image: "https://upload.wikimedia.org/wikipedia/commons/6/69/Ancient_Tamil_Script.jpg",
				},
				{
					Name:  "Books",
					Image: "https://upload.wikimedia.org/wikipedia/commons/a/a0/Book_fair-Tamil_Nadu-35th-Chennai-january-2012-part_30.JPG",
				},
				{
					Name:  "Poem",
					Image: "https://upload.wikimedia.org/wikipedia/commons/f/ff/Subramanya_Bharathi_1960_stamp_of_India.jpg",
				},
			},
		},
		{
			Category: "Entertainment",
			Interests: []db.InterestWithIcon{
				{
					Name:  "Anime",
					Image: "https://www.thenerddaily.com/wp-content/uploads/2018/08/Reasons-To-Watch-Anime.jpg",
				},
				{
					Name:  "Series",
					Image: "https://upload.wikimedia.org/wikipedia/commons/1/10/Meta-image-netflix-symbol-black.png",
				},
				{
					Name:  "Movies",
					Image: "https://upload.wikimedia.org/wikipedia/commons/c/c1/Rajinikanth_and_Vijay_at_the_Nadigar_Sangam_Protest.jpg",
				},
			},
		},
		{
			Category: "Culture",
			Interests: []db.InterestWithIcon{
				{
					Name:  "Philosophy",
					Image: "https://en.wikipedia.org/wiki/Swami_Vivekananda#/media/File:Swami_Vivekananda-1893-09-signed.jpg",
				},
				{
					Name:  "Language",
					Image: "https://upload.wikimedia.org/wikipedia/commons/3/35/Word_Tamil.svg",
				},
				{
					Name:  "Festival",
					Image: "https://images.unsplash.com/photo-1576394435759-02a2674ff6e0",
				},
				{
					Name:  "Agriculture",
					Image: "https://upload.wikimedia.org/wikipedia/commons/f/f1/%281%29_Agriculture_and_rural_farms_of_India.jpg",
				},
				{
					Name:  "Cooking",
					Image: "https://images.unsplash.com/photo-1604740795024-c06eeca4bf89",
				},
			},
		},
	}

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
					IsFollowedByUser: true,
				},
				{
					Name:             "Poem",
					Image:            "https://upload.wikimedia.org/wikipedia/commons/f/ff/Subramanya_Bharathi_1960_stamp_of_India.jpg",
					IsFollowedByUser: true,
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
	userFollowingInterests := []string{"Poem", "Art", "Books"}
	suite.mockUserProfileClient.EXPECT().FetchUserFollowingInterests(suite.goContext).Return(userFollowingInterests, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().FetchCategoriesAndInterests(suite.goContext).Return(dbReturnedCategoriesAndInterest, nil).Times(1)
	suite.mockInterestsMapper.EXPECT().MapUserFollowedInterest(suite.goContext, dbReturnedCategoriesAndInterest, userFollowingInterests).Return(expectedData).Times(1)

	actualInterests, err := suite.interestsService.GetExploreCategoriesAndInterests(suite.goContext)

	suite.Nil(err)
	suite.Equal(expectedData, actualInterests)
}

func (suite *InterestsServiceTest) TestGetExploreCategoriesAndInterests_WhenDbReturnsError() {
	userFollowingInterests := []string{"Poem", "Art", "Books"}
	suite.mockUserProfileClient.EXPECT().FetchUserFollowingInterests(suite.goContext).Return(userFollowingInterests, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().FetchCategoriesAndInterests(suite.goContext).Return(nil, errors.New(test_helper.ErrSomethingWentWrong)).Times(1)
	interests, err := suite.interestsService.GetExploreCategoriesAndInterests(suite.goContext)
	suite.NotNil(err)
	suite.Equal(&constants.PostServiceFailureError, err)
	suite.Nil(interests)
}

func (suite *InterestsServiceTest) TestGetExploreCategoriesAndInterests_WhenNoDataReturnedWithError() {
	userFollowingInterests := []string{"Poem", "Art", "Books"}

	suite.mockUserProfileClient.EXPECT().FetchUserFollowingInterests(suite.goContext).Return(userFollowingInterests, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().FetchCategoriesAndInterests(suite.goContext).Return(nil, errors.New(constants.NoInterestsAndCategoriesCode)).Times(1)

	interests, err := suite.interestsService.GetExploreCategoriesAndInterests(suite.goContext)
	suite.NotNil(err)
	suite.Equal(&constants.NoInterestsAndCategoriesErr, err)
	suite.Nil(interests)
}

func (suite *InterestsServiceTest) TestGetExploreCategoriesAndInterests_WhenNoDataReturnedWithNoError() {
	userFollowingInterests := []string{"Poem", "Art", "Books"}
	suite.mockUserProfileClient.EXPECT().FetchUserFollowingInterests(suite.goContext).Return(userFollowingInterests, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().FetchCategoriesAndInterests(suite.goContext).Return([]db.CategoryAndInterest{}, nil).Times(1)
	interests, err := suite.interestsService.GetExploreCategoriesAndInterests(suite.goContext)
	suite.NotNil(err)
	suite.Equal(&constants.NoInterestsAndCategoriesErr, err)
	suite.Nil(interests)
}

func (suite *InterestsServiceTest) TestGetExploreCategoriesAndInterests_WhenNilReturnedWithNoErrorFromB() {
	userFollowingInterests := []string{"Poem", "Art", "Books"}
	suite.mockUserProfileClient.EXPECT().FetchUserFollowingInterests(suite.goContext).Return(userFollowingInterests, nil).Times(1)
	suite.mockInterestsRepository.EXPECT().FetchCategoriesAndInterests(suite.goContext).Return(nil, nil).Times(1)
	interests, err := suite.interestsService.GetExploreCategoriesAndInterests(suite.goContext)
	suite.NotNil(err)
	suite.Equal(&constants.NoInterestsAndCategoriesErr, err)
	suite.Nil(interests)
}

func (suite *InterestsServiceTest) TestGetExploreCategoriesAndInterests_WhenFetchUserFollowInterestsReturnsError() {
	dbReturnedCategoriesAndInterest := []db.CategoryAndInterest{
		{
			Category: "Art",
			Interests: []db.InterestWithIcon{
				{
					Name:  "Comics",
					Image: "https://upload.wikimedia.org/wikipedia/ta/9/9d/Sirithiran_cartoon_2.JPG",
				},
				{
					Name:  "Literature",
					Image: "https://upload.wikimedia.org/wikipedia/commons/6/69/Ancient_Tamil_Script.jpg",
				},
				{
					Name:  "Books",
					Image: "https://upload.wikimedia.org/wikipedia/commons/a/a0/Book_fair-Tamil_Nadu-35th-Chennai-january-2012-part_30.JPG",
				},
				{
					Name:  "Poem",
					Image: "https://upload.wikimedia.org/wikipedia/commons/f/ff/Subramanya_Bharathi_1960_stamp_of_India.jpg",
				},
			},
		},
		{
			Category: "Entertainment",
			Interests: []db.InterestWithIcon{
				{
					Name:  "Anime",
					Image: "https://www.thenerddaily.com/wp-content/uploads/2018/08/Reasons-To-Watch-Anime.jpg",
				},
				{
					Name:  "Series",
					Image: "https://upload.wikimedia.org/wikipedia/commons/1/10/Meta-image-netflix-symbol-black.png",
				},
				{
					Name:  "Movies",
					Image: "https://upload.wikimedia.org/wikipedia/commons/c/c1/Rajinikanth_and_Vijay_at_the_Nadigar_Sangam_Protest.jpg",
				},
			},
		},
		{
			Category: "Culture",
			Interests: []db.InterestWithIcon{
				{
					Name:  "Philosophy",
					Image: "https://en.wikipedia.org/wiki/Swami_Vivekananda#/media/File:Swami_Vivekananda-1893-09-signed.jpg",
				},
				{
					Name:  "Language",
					Image: "https://upload.wikimedia.org/wikipedia/commons/3/35/Word_Tamil.svg",
				},
				{
					Name:  "Festival",
					Image: "https://images.unsplash.com/photo-1576394435759-02a2674ff6e0",
				},
				{
					Name:  "Agriculture",
					Image: "https://upload.wikimedia.org/wikipedia/commons/f/f1/%281%29_Agriculture_and_rural_farms_of_India.jpg",
				},
				{
					Name:  "Cooking",
					Image: "https://images.unsplash.com/photo-1604740795024-c06eeca4bf89",
				},
			},
		},
	}

	suite.mockUserProfileClient.EXPECT().FetchUserFollowingInterests(suite.goContext).Return(nil, &constants.InternalServerError).Times(1)
	suite.mockInterestsRepository.EXPECT().FetchCategoriesAndInterests(suite.goContext).Return(dbReturnedCategoriesAndInterest, nil).Times(0)

	_, err := suite.interestsService.GetExploreCategoriesAndInterests(suite.goContext)

	suite.NotNil(err)
	suite.Equal(&constants.InternalServerError, err)
}
