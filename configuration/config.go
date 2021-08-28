package configuration

import (
	"github.com/gola-glitch/gola-utils/model"
	"github.com/gola-glitch/gola-utils/redis_util"
)

type ConfigData struct {
	TracingServiceName        string                       `json:"tracing_service_name" binding:"required"`
	TracingOCAgentHost        string                       `json:"tracing_oc_agent_host" binding:"required"`
	DBConnectionPool          model.DBConnectionPoolConfig `json:"dbConnectionPool" binding:"required"`
	LogLevel                  string                       `json:"log_level" binding:"required"`
	ContentReadTimeConfig     map[string]int               `json:"content_read_time_config" binding:"required"`
	MinimumPostReadTime       int                          `json:"minimum_post_read_time" binding:"required"`
	Environment               string                       `json:"environment" binding:"required"`
	AllowedOrigins            []string                     `json:"allowed_origins" binding:"required"`
	RedisStoreConfig          redis_util.RedisStoreConfig  `json:"redis" binding:"required"`
	CryptoServiceURL          string                       `json:"crypto_service_url" binding:"required"`
	Email                     Email                        `json:"email" binding:"required"`
	Oauth                     OAuth                        `json:"oauth" binding:"required"`
	RequestTimeOut            int                          `json:"request_time_out" binding:"required"`
	AllowInsecureCookies      bool                         `json:"allow_insecure_cookies"`
	ActivationCallback        string                       `json:"activationCallback" binding:"required"`
	TokenValidationIgnoreURLs []string                     `json:"token_validation_ignore_urls"`
}

type Email struct {
	GatewayURL    string         `json:"gateway_url"`
	DefaultSender string         `json:"default_sender"`
	TemplatePaths TemplatesPaths `json:"template_paths"`
}

type OAuth struct {
	AdminUrl                string `json:"adminBaseUrl"`
	PublicUrl               string `json:"public_url"`
	AcceptLoginRequestUrl   string `json:"accept_login_request_url"`
	GetConsentRequestUrl    string `json:"get_consent_request_url"`
	AcceptConsentRequestUrl string `json:"accept_consent_request_url"`
	GetTokenUrl             string `json:"get_token_url"`
}
type TemplatesPaths struct {
	NewUserActivation string `json:"new_user_activation"`
}

func (configData *ConfigData) GetDBConnectionPoolConfig() model.DBConnectionPoolConfig {
	return configData.DBConnectionPool
}
