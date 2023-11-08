package service

import (
	"github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"post-api/configuration"
	"testing"
)

type ProtectedUrlServiceTestSuite struct {
	suite.Suite
	mockConfigData      *configuration.ConfigData
	protectedURLService service.ProtectedUrlService
}

func TestProtectedUrlServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ProtectedUrlServiceTestSuite))
}

func (suite *ProtectedUrlServiceTestSuite) SetupTest() {
	suite.mockConfigData = &configuration.ConfigData{
		TokenValidationIgnoreURLs: []string{
			"/api/post/healthz",
			"/api/post/some-url",
		},
	}
	suite.protectedURLService = NewProtectedUrlService(suite.mockConfigData)
}

func (suite *ProtectedUrlServiceTestSuite) TearDownTest() {
}

func (suite *ProtectedUrlServiceTestSuite) TestProtectedUrlService_IsProtected() {
	isProtected := suite.protectedURLService.IsProtected("/api/post/v1/draft/upsertDraft")
	assert.Equal(suite.T(), true, isProtected)
}

func (suite *ProtectedUrlServiceTestSuite) TestProtectedUrlService_ExcludeUnprotectedUrls() {
	isProtected := suite.protectedURLService.IsProtected("/api/post/healthz")
	assert.Equal(suite.T(), false, isProtected)
}

func (suite *ProtectedUrlServiceTestSuite) TestProtectedUrlService_ExcludeUnprotectedUrlsHavingPathParameters() {
	isProtected := suite.protectedURLService.IsProtected("/api/post/some-url/123")
	assert.Equal(suite.T(), false, isProtected)
}
