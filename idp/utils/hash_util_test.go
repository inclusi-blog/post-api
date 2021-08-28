package util

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type HashUtilTestSuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	hashUtil HashUtil
}

func TestHashUtilTestSuite(t *testing.T) {
	suite.Run(t, new(HashUtilTestSuite))

}

func (suite *HashUtilTestSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.hashUtil = NewHashUtil()
}

func (suite *HashUtilTestSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite HashUtilTestSuite) TestGenerateBcryptHash_ShouldReturnBcryptHash() {
	actual, err := suite.hashUtil.GenerateBcryptHash("some-data")
	suite.NotNil(actual)
	suite.Nil(err)
}

func (suite HashUtilTestSuite) TestMatchBcryptHash_ShouldReturnNilWhenMatchFound() {
	err := suite.hashUtil.MatchBcryptHash("$2a$10$KcRz/b/Frfi4aSHP.qNUj.xOEoWKFn2XvV9Qf39WYO4Ip1naMGzWW", "some-data")
	suite.Nil(err)
}

func (suite HashUtilTestSuite) TestMatchBcryptHash_ShouldReturnErrorWhenNoMatchFound() {
	err := suite.hashUtil.MatchBcryptHash("some-data", "some-data")
	suite.NotNil(err)
}
