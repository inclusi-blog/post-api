package service

import (
	"github.com/gola-glitch/gola-utils/middleware/introspection/oauth-middleware/service"
	"post-api/configuration"
	"strings"
)

type protectedUrlService struct {
	configData *configuration.ConfigData
}

func NewProtectedUrlService(configData *configuration.ConfigData) service.ProtectedUrlService {
	return protectedUrlService{configData: configData}
}

func (protectedUrlService protectedUrlService) IsProtected(whiteListUrl string) bool {
	for _, url := range protectedUrlService.configData.TokenValidationIgnoreURLs {
		if strings.Contains(whiteListUrl, url) {
			return false
		}
	}
	return true
}
