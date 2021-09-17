package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"post-api/dbhelper"
	"post-api/test_helper/helper"
	"testing"
)

type UserInterestsRepositoryTest struct {
	suite.Suite
	db                      *sqlx.DB
	ginContext              *gin.Context
	recorder                *httptest.ResponseRecorder
	dbHelper                helper.DbHelper
	userHelper              helper.UserRepository
	userInterestsRepository UserInterestsRepository
}

func (suite *UserInterestsRepositoryTest) SetupTest() {
	err := godotenv.Load("../../docker-compose-test.env")
	suite.Nil(err)
	connectionString := dbhelper.BuildConnectionString()
	database, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintln("Could not connect to test DB", err))
	}
	suite.db = database
	suite.recorder = httptest.NewRecorder()
	suite.ginContext, _ = gin.CreateTestContext(suite.recorder)
	suite.ginContext.Request = httptest.NewRequest(http.MethodGet, "/some-url", nil)
	suite.userInterestsRepository = NewUserInterestsRepository(database)
	suite.dbHelper = helper.NewDbHelper(database)
}

func (suite *UserInterestsRepositoryTest) TearDownTest() {
	suite.ClearData()
	_ = suite.db.Close()
}

func (suite *UserInterestsRepositoryTest) ClearData() {
	e := suite.dbHelper.ClearAll()
	if e != nil {
		assert.Error(suite.T(), e)
	}
}

func TestUserInterestsRepositoryTest(t *testing.T) {
	suite.Run(t, new(UserInterestsRepositoryTest))
}

func (suite *UserInterestsRepositoryTest) TestGetFollowedInterests_WhenNotFollowedAnyInterests() {
	_, err := suite.userInterestsRepository.GetFollowedInterest(suite.ginContext, uuid.New())
	suite.Nil(err)
}
