package init

import (
	"github.com/gin-gonic/gin"
	"post-api/configuration"
)

func CreateRouter(data *configuration.ConfigData) *gin.Engine {
	router := gin.Default()
	Swagger()
	_ = Db(data)
	RegisterRouter(router, data)
	return router
}
