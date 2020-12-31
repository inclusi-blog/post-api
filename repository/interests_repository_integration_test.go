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

func (suite *InterestsRepositoryIntegrationTest) insertInterestEntries() {
	interests := []string{
		"CREATE (interest:Category:Interest{name: 'Art'})",
		"CREATE (interest:Category:Interest{name: 'Entertainment'})",
		"CREATE (interest:Category:Interest{name: 'Culture'})",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Poem'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Short stories'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Books'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Literature'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Grammar'})-[:BELONGS_TO]->(art)",
		"MATCH (art:Category{name: 'Art'}) CREATE (interest:Interest{name: 'Comics'})-[:BELONGS_TO]->(art)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Movies'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Series'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Anime'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Cartoon'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (entertainment:Category{name: 'Entertainment'}) CREATE (interest:Interest{name: 'Animation'})-[:BELONGS_TO]->(entertainment)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Cooking'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Food'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Agriculture'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Festival'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Language'})-[:BELONGS_TO]->(culture)",
		"MATCH (culture:Category{name: 'Culture'}) CREATE (interest:Interest{name: 'Philosophy'})-[:BELONGS_TO]->(culture)",
	}
	for _, query := range interests {
		_, err := suite.adminDb.Run(query, nil)
		suite.Nil(err)
	}
}
