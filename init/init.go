package init

import (
	"github.com/gin-gonic/gin"
	"post-api/configuration"
	"post-api/utils"
)

func CreateRouter(data *configuration.ConfigData) *gin.Engine {
	router := gin.Default()
	Validators()
	Swagger()
	db := Db(data)
	HttpClient(data)
	aws := utils.ConnectAws(data)
	Objects(db, data, aws)
	RegisterRouter(router, data)
	return router
}
