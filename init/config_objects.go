package init

import (
	"github.com/inclusi-blog/gola-utils/configuration_loader"
	"post-api/configuration"
	"post-api/constants"
)

func LoadConfig() *configuration.ConfigData {
	var configData configuration.ConfigData
	err := configuration_loader.NewConfigLoader().Load(constants.FILE_NAME, &configData)

	if err != nil {
		panic(err)
	}
	return &configData
}
