package configuration

import "github.com/gola-glitch/gola-utils/model"

type ConfigData struct {
	TracingServiceName         string                       `json:"tracing_service_name" binding:"required"`
	TracingOCAgentHost         string                       `json:"tracing_oc_agent_host" binding:"required"`
	DBConnectionPool           model.DBConnectionPoolConfig `json:"dbConnectionPool" binding:"required"`
	LogLevel                   string                       `json:"log_level" binding:"required"`
	ContentReadTimeConfig      map[string]int               `json:"content_read_time_config" binding:"required"`
	MinimumPostReadTime        int                          `json:"minimum_post_read_time" binding:"required"`
	Environment                string                       `json:"environment" binding:"required"`
	AllowedOrigins             []string                     `json:"allowed_origins" binding:"required"`
	CryptoServiceUrl           string                       `json:"crypto_service_url" binding:"required"`
	OAuthUrl                   string                       `json:"oauth_url" binding:"required"`
	TokenValidationIgnoreURLs  []string                     `json:"token_validation_ignore_urls"`
	UserProfileBaseUrl         string                       `json:"user_profile_base_url" binding:"required"`
	RequestTimeout             int                          `json:"request_timeout" binding:"required"`
	FetchUserFollowedInterests string                       `json:"fetch_user_followed_interests" binding:"required"`
}

func (configData *ConfigData) GetDBConnectionPoolConfig() model.DBConnectionPoolConfig {
	return configData.DBConnectionPool
}
