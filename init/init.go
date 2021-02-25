package init

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gola-glitch/gola-utils/logging"
	"post-api/configuration"
)

func CreateRouter(data *configuration.ConfigData) *gin.Engine {
	router := gin.Default()
	Validators()
	HttpClient(data)
	Swagger()
	db := Db(data)
	Objects(db, data)
	RegisterRouter(router, data)

	defer func() {
		err := db.Close()
		if err != nil {
			logging.GetLogger(context.TODO()).Errorf("unable to close session %v", err)
		}
	}()
	return router
}
