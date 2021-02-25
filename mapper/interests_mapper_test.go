package mapper

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"post-api/models/db"
	"post-api/models/response"
	"testing"
)

type InterestsMapperTest struct {
	suite.Suite
	mockController  *gomock.Controller
	goContext       context.Context
	interestsMapper InterestsMapper
}

func TestInterestsMapperTestSuite(t *testing.T) {
	suite.Run(t, new(InterestsMapperTest))
}

func (suite *InterestsMapperTest) SetupTest() {
	suite.mockController = gomock.NewController(suite.T())
	suite.goContext = context.WithValue(context.Background(), "someKey", "someValue")
	suite.interestsMapper = NewInterestsMapper()
}

func (suite *InterestsMapperTest) TearDownTest() {
	suite.mockController.Finish()
}

func (suite *InterestsMapperTest) TestMapUserFollowedInterest_WhenValidDataSent() {
	categoryAndInterests := []db.CategoryAndInterest{
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
	userFollowingInterests := []string{"Books", "Series", "Festival"}
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
					IsFollowedByUser: true,
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
					IsFollowedByUser: true,
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
	categoriesAndInterests := suite.interestsMapper.MapUserFollowedInterest(suite.goContext, categoryAndInterests, userFollowingInterests)
	suite.Len(categoriesAndInterests, 3)
	suite.EqualValues(expectedData, categoriesAndInterests)
}

func (suite *InterestsMapperTest) TestMapUserFollowedInterest_WhenEmptyDataSentShouldReturnEmptyData() {
	var categoryAndInterests []db.CategoryAndInterest
	var expectedData []response.CategoryAndInterest
	categoriesAndInterests := suite.interestsMapper.MapUserFollowedInterest(suite.goContext, categoryAndInterests, []string{})
	suite.Len(categoriesAndInterests, 0)
	suite.EqualValues(expectedData, categoriesAndInterests)
}
