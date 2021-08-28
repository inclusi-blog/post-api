package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"idp-api/dbhelper"
	"idp-api/models/db"
	"idp-api/repository/helper"
	"testing"
)

type UserDetailsRepositoryIntegrationTest struct { //test suite
	suite.Suite
	goContext             context.Context
	userDetailsRepository UserDetailsRepository
	db                    *sqlx.DB
	dbHelper              helper.DbHelper
}

func (suite *UserDetailsRepositoryIntegrationTest) SetupTest() {
	err := godotenv.Load("../docker-compose-test.env")
	suite.Nil(err)
	connectionString := dbhelper.BuildConnectionString()
	database, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		panic(fmt.Sprintln("Could not connect to test DB", err))
	}
	suite.db = database
	suite.goContext = context.WithValue(context.Background(), "testKey", "testVal")
	suite.userDetailsRepository = NewUserDetailsRepository(database)
	suite.dbHelper = helper.NewDbHelper(database)
}

func (suite *UserDetailsRepositoryIntegrationTest) TearDownTest() {
	suite.ClearData()
	_ = suite.db.Close()
}

func (suite *UserDetailsRepositoryIntegrationTest) ClearData() {
	e := suite.dbHelper.ClearAll()
	if e != nil {
		assert.Error(suite.T(), e)
	}
}

func TestUserDetailsRepositoryIntegrationTest(t *testing.T) {
	suite.Run(t, new(UserDetailsRepositoryIntegrationTest))
}

func (suite *UserDetailsRepositoryIntegrationTest) TestSaveUserDetails_WhenValidData() {
	details := db.SaveUserDetails{
		UUID:     "some-random-string",
		Username: "someuser",
		Email:    "someuser@gmail.com",
		Password: "cksacba@!#$^*(*%$!UT!VBH!@B!B@INcksacba@!#$^*(*%$!UT!VBH!@B!",
		IsActive: false,
	}
	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestSaveUserDetails_WhenDbReturnsError() {
	details := db.SaveUserDetails{
		UUID:     "some-random-string",
		Username: "someuser",
		Email:    "someuser@gmail.com",
		Password: "cksacba@!#$^*(*%$!UT!VBH!@B!B@INcksacba@!#$^*(*%$!UT!VBH!@B!",
		IsActive: false,
	}
	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)
	suite.Nil(err)
	err = suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)
	suite.NotNil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestIsUserAvailable_WhenUserExists() {
	details := db.SaveUserDetails{
		UUID:     "some-random-string",
		Username: "someuser",
		Email:    "someuser@gmail.com",
		Password: "cksacba@!#$^*(*%$!UT!VBH!@B!B@INcksacba@!#$^*(*%$!UT!VBH!@B!",
		IsActive: false,
	}
	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)
	suite.Nil(err)

	available, err := suite.userDetailsRepository.IsEmailAvailable(details.Email, suite.goContext)

	suite.True(available)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestIsUserAvailable_WhenUserNotExists() {
	details := db.SaveUserDetails{
		UUID:     "some-random-string",
		Username: "someuser",
		Email:    "someuser@gmail.com",
		Password: "cksacba@!#$^*(*%$!UT!VBH!@B!B@INcksacba@!#$^*(*%$!UT!VBH!@B!",
		IsActive: false,
	}

	available, err := suite.userDetailsRepository.IsEmailAvailable(details.Email, suite.goContext)

	suite.False(available)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestIsUserNameAvailable_WhenUserNameExists() {
	details := db.SaveUserDetails{
		UUID:     "some-random-string",
		Username: "someuser",
		Email:    "someuser@gmail.com",
		Password: "cksacba@!#$^*(*%$!UT!VBH!@B!B@INcksacba@!#$^*(*%$!UT!VBH!@B!",
		IsActive: false,
	}

	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)
	suite.Nil(err)

	available, err := suite.userDetailsRepository.IsUserNameAvailable(details.Username, suite.goContext)

	suite.True(available)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestIsUserNameAvailable_WhenUserNameNotExists() {
	available, err := suite.userDetailsRepository.IsUserNameAvailable("some-username", suite.goContext)

	suite.False(available)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestIsUserNameAndEmailAvailable_WhenUserNameExists() {
	details := db.SaveUserDetails{
		UUID:     "some-random-string",
		Username: "someuser",
		Email:    "someuser@gmail.com",
		Password: "cksacba@!#$^*(*%$!UT!VBH!@B!B@INcksacba@!#$^*(*%$!UT!VBH!@B!",
		IsActive: false,
	}

	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)
	suite.Nil(err)

	available, err := suite.userDetailsRepository.IsUserNameAndEmailAvailable(details.Username, "random@gmail.com", suite.goContext)

	suite.True(available)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestIsUserNameAndEmailAvailable_WhenEmailExists() {
	details := db.SaveUserDetails{
		UUID:     "some-random-string",
		Username: "someuser",
		Email:    "someuser@gmail.com",
		Password: "cksacba@!#$^*(*%$!UT!VBH!@B!B@INcksacba@!#$^*(*%$!UT!VBH!@B!",
		IsActive: false,
	}

	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)
	suite.Nil(err)

	available, err := suite.userDetailsRepository.IsUserNameAndEmailAvailable("notfoun-user", "someuser@gmail.com", suite.goContext)

	suite.True(available)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestIsUserNameAndEmailAvailable_WhenBothUserNameAndEmailNotExists() {
	available, err := suite.userDetailsRepository.IsUserNameAndEmailAvailable("notfoun-user", "someuser@gmail.com", suite.goContext)

	suite.False(available)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestIsUserNameAndEmailAvailable_WhenBothUserNameAndEmailExists() {
	details := db.SaveUserDetails{
		UUID:     "some-random-string",
		Username: "someuser",
		Email:    "someuser@gmail.com",
		Password: "cksacba@!#$^*(*%$!UT!VBH!@B!B@INcksacba@!#$^*(*%$!UT!VBH!@B!",
		IsActive: false,
	}

	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)
	suite.Nil(err)

	available, err := suite.userDetailsRepository.IsUserNameAndEmailAvailable("someuser", "someuser@gmail.com", suite.goContext)

	suite.True(available)
	suite.Nil(err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestGetUserProfile_WhenDbReturnsUserProfile() {
	profile := db.UserProfile{
		UserID:   "some-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		IsActive: true,
	}

	details := db.SaveUserDetails{
		UUID:     "some-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		Password: "k$q3!@CAF#!@ASsdS!@!@",
		IsActive: true,
	}
	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)

	suite.Nil(err)

	userProfile, err := suite.userDetailsRepository.GetUserProfile("dummy@gmail.com", suite.goContext)

	suite.Nil(err)
	suite.Equal(profile, userProfile)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestGetUserProfile_WhenDbReturnsErrorOnEmptyProfile() {
	profile := db.UserProfile{}

	userProfile, err := suite.userDetailsRepository.GetUserProfile("dummy@gmail.com", suite.goContext)

	suite.NotNil(err)
	suite.Equal(profile, userProfile)
	suite.Equal(sql.ErrNoRows, err)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestGetPassword_WhenDbReturnsPassword() {

	expectedPassword := "k$q3!@CAF#!@ASsdS!@!@"

	details := db.SaveUserDetails{
		UUID:     "some-id",
		Username: "some-user",
		Email:    "dummy@gmail.com",
		Password: "k$q3!@CAF#!@ASsdS!@!@",
		IsActive: true,
	}
	err := suite.userDetailsRepository.SaveUserDetails(details, suite.goContext)

	suite.Nil(err)

	actualPassword, err := suite.userDetailsRepository.GetPassword("dummy@gmail.com", suite.goContext)

	suite.Nil(err)
	suite.Equal(expectedPassword, actualPassword)
}

func (suite *UserDetailsRepositoryIntegrationTest) TestGetPassword_WhenDbReturnsNoRows() {

	actualPassword, err := suite.userDetailsRepository.GetPassword("dummy@gmail.com", suite.goContext)

	suite.NotNil(err)
	suite.Equal(sql.ErrNoRows, err)
	suite.Equal("", actualPassword)
}
