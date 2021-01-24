package repository

import (
	"context"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"post-api/constants"
	"post-api/dbhelper"
	"post-api/models/db"
	"post-api/repository/helper"
	"testing"
)

type InterestsRepositoryIntegrationTest struct {
	suite.Suite
	db                  neo4j.Session
	driver              neo4j.Driver
	adminDb             neo4j.Session
	adminDriver         neo4j.Driver
	goContext           context.Context
	interestsRepository InterestsRepository
	dbHelper            helper.DbHelper
}

func (suite *InterestsRepositoryIntegrationTest) SetupTest() {
	err := godotenv.Load("../docker-compose-test.env")
	suite.Nil(err)
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	logger := logging.GetLogger(context.Background())
	configForNeo4j40 := func(conf *neo4j.Config) { conf.Encrypted = false }
	suite.driver, err = neo4j.NewDriver(dbhelper.BuildConnectionString(), neo4j.BasicAuth(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), ""), configForNeo4j40)
	suite.Nil(err)
	suite.adminDriver, err = neo4j.NewDriver(dbhelper.BuildConnectionString(), neo4j.BasicAuth(os.Getenv("ADMIN_USER"), os.Getenv("ADMIN_PASSWORD"), ""), configForNeo4j40)
	suite.Nil(err)
	suite.NotNil(suite.adminDriver)
	suite.NotNil(suite.driver)

	logger.Info("logging")
	sessionConfig := neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: os.Getenv("DB_SERVICE_NAME")}
	suite.db, err = suite.driver.NewSession(sessionConfig)
	suite.Nil(err)
	suite.adminDb, err = suite.adminDriver.NewSession(sessionConfig)
	suite.Nil(err)

	suite.interestsRepository = NewInterestRepository(suite.db)
	suite.dbHelper = helper.NewDbHelper(suite.adminDb)
	err = suite.dbHelper.ClearAll()
	suite.Nil(err)
	suite.insertInterestEntries()
}

func (suite *InterestsRepositoryIntegrationTest) TearDownTest() {
	suite.ClearInterestsData()
	err := suite.driver.Close()
	suite.Nil(err)
	err = suite.adminDriver.Close()
	suite.Nil(err)
	err = suite.db.Close()
	suite.Nil(err)
	err = suite.adminDb.Close()
	suite.Nil(err)
}

func (suite *InterestsRepositoryIntegrationTest) ClearInterestsData() {
	e := suite.dbHelper.ClearAll()
	if e != nil {
		assert.Error(suite.T(), e)
	}
}

func TestInterestsRepositoryIntegrationTest(t *testing.T) {
	suite.Run(t, new(InterestsRepositoryIntegrationTest))
}

func (suite *InterestsRepositoryIntegrationTest) TestGetInterests_WhenDbReturnsData() {
	interests, err := suite.interestsRepository.GetInterests(suite.goContext, "", []string{"Sports"})
	suite.Nil(err)
	suite.Equal(20, len(interests))
}

func (suite *InterestsRepositoryIntegrationTest) TestGetInterests_WhenSearchKeywordPassedDbReturnsData() {
	interests, err := suite.interestsRepository.GetInterests(suite.goContext, "Festival", []string{})
	suite.Nil(err)
	suite.Equal(1, len(interests))
}

func (suite *InterestsRepositoryIntegrationTest) TestGetInterests_WhenNoInterestsAvailable() {
	suite.ClearInterestsData()
	interests, err := suite.interestsRepository.GetInterests(suite.goContext, "Festival", []string{})
	suite.NotNil(err)
	suite.Equal(0, len(interests))
	suite.Equal(constants.NoInterestsFoundCode, err.Error())
}

func (suite *InterestsRepositoryIntegrationTest) TestFetchCategoriesAndInterests_WhenDbReturnsData() {
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
	interests, err := suite.interestsRepository.FetchCategoriesAndInterests(suite.goContext)
	suite.Nil(err)
	suite.Equal(3, len(interests))
	suite.EqualValues(categoryAndInterests, interests)
}

func (suite *InterestsRepositoryIntegrationTest) TestFetchCategoriesAndInterests_WhenNoInterestsAvailable() {
	suite.ClearInterestsData()
	interests, err := suite.interestsRepository.FetchCategoriesAndInterests(suite.goContext)
	suite.NotNil(err)
	suite.Nil(interests)
	suite.Equal(constants.NoInterestsAndCategoriesCode, err.Error())
}

func (suite *InterestsRepositoryIntegrationTest) insertInterestEntries() {
	interests := []string{
		"CREATE (interest:Category:Interest{name: 'Art'})",
		"CREATE (interest:Category:Interest{name: 'Entertainment'})",
		"CREATE (interest:Category:Interest{name: 'Culture'})",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Poem', image: 'https://upload.wikimedia.org/wikipedia/commons/f/ff/Subramanya_Bharathi_1960_stamp_of_India.jpg'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Short stories'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Books', image: 'https://upload.wikimedia.org/wikipedia/commons/a/a0/Book_fair-Tamil_Nadu-35th-Chennai-january-2012-part_30.JPG'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Literature', image: 'https://upload.wikimedia.org/wikipedia/commons/6/69/Ancient_Tamil_Script.jpg'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Grammar'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Comics', image: 'https://upload.wikimedia.org/wikipedia/ta/9/9d/Sirithiran_cartoon_2.JPG'})-[:BELONGS_TO]->(art)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Movies', image: 'https://upload.wikimedia.org/wikipedia/commons/c/c1/Rajinikanth_and_Vijay_at_the_Nadigar_Sangam_Protest.jpg'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Series', image: 'https://upload.wikimedia.org/wikipedia/commons/1/10/Meta-image-netflix-symbol-black.png'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Anime', image: 'https://www.thenerddaily.com/wp-content/uploads/2018/08/Reasons-To-Watch-Anime.jpg'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Cartoon'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Animation'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Cooking', image: 'https://images.unsplash.com/photo-1604740795024-c06eeca4bf89'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Food'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Agriculture', image: 'https://upload.wikimedia.org/wikipedia/commons/f/f1/%281%29_Agriculture_and_rural_farms_of_India.jpg'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Festival', image: 'https://images.unsplash.com/photo-1576394435759-02a2674ff6e0'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Language', image: 'https://upload.wikimedia.org/wikipedia/commons/3/35/Word_Tamil.svg'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Philosophy', image: 'https://en.wikipedia.org/wiki/Swami_Vivekananda#/media/File:Swami_Vivekananda-1893-09-signed.jpg'})-[:BELONGS_TO]->(culture)",
	}
	for _, query := range interests {
		_, err := suite.adminDb.Run(query, nil)
		suite.Nil(err)
	}
}
