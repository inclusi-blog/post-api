package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"post-api/dbhelper"
	"post-api/repository/helper"
	"testing"
)

type InterestsRepositoryIntegrationTest struct {
	suite.Suite
	db                  *sqlx.DB
	goContext           context.Context
	interestsRepository InterestsRepository
	dbHelper            helper.DbHelper
}

func (suite *InterestsRepositoryIntegrationTest) SetupTest() {
	connectionString := dbhelper.BuildConnectionString()
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintln("Could not connect to test DB", err))
	}
	fmt.Print(db)
	suite.db = db
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	suite.interestsRepository = NewInterestRepository(db)
	suite.dbHelper = helper.NewDbHelper(db)
}

func (suite *InterestsRepositoryIntegrationTest) TearDownTest() {
	suite.ClearDraftData()
	_ = suite.db.Close()
}

func (suite *InterestsRepositoryIntegrationTest) ClearDraftData() {
	e := suite.dbHelper.ClearAll()
	if e != nil {
		assert.Error(suite.T(), e)
	}
}

func TestInterestsRepositoryIntegrationTest(t *testing.T) {
	suite.Run(t, new(InterestsRepositoryIntegrationTest))
}

func (suite *InterestsRepositoryIntegrationTest) TestGetInterests_WhenDbReturnsData() {
	interests, err := suite.interestsRepository.GetInterests(suite.goContext)
	suite.Nil(err)
	suite.Equal(14, len(interests))
}
