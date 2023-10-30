package init

import (
	"github.com/gin-gonic/gin"
	"post-api/configuration"
	"post-api/utils"
)

func CreateRouter(data *configuration.ConfigData) *gin.Engine {
	router := gin.Default()
	aws := utils.ConnectAws()
	Validators()
	Swagger()
	HttpClient(data)
	db := Db(data, aws)
	Objects(db, data, aws)
	RegisterRouter(router, data)
	return router
}
