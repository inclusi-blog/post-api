package configuration

import "github.com/gola-glitch/gola-utils/model"

type ConfigData struct {
	TracingServiceName    string                       `json:"tracing_service_name" binding:"required"`
	TracingOCAgentHost    string                       `json:"tracing_oc_agent_host" binding:"required"`
	DBConnectionPool      model.DBConnectionPoolConfig `json:"dbConnectionPool" binding:"required"`
	LogLevel              string                       `json:"log_level" binding:"required"`
	ContentReadTimeConfig map[string]int               `json:"content_read_time_config" binding:"required"`
	MinimumPostReadTime   int                          `json:"minimum_post_read_time" binding:"required"`
	Environment           string                       `json:"environment" binding:"required"`
	AllowedOrigins        []string                     `json:"allowed_origins" binding:"required"`
}

func (configData *ConfigData) GetDBConnectionPoolConfig() model.DBConnectionPoolConfig {
	return configData.DBConnectionPool
}
