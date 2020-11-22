package init

import (
	"github.com/gin-gonic/gin"
	"post-api/configuration"
)

func CreateRouter(data *configuration.ConfigData) *gin.Engine {
	router := gin.Default()
	Validators()
	Swagger()
	db := Db(data)
	Objects(db, data)
	RegisterRouter(router, data)
	return router
}
