package init

import (
	"github.com/jmoiron/sqlx"
	"post-api/configuration"
	"post-api/controller"
	"post-api/repository"
	"post-api/service"
)

var (
	draftController controller.DraftController
)

func Objects(db *sqlx.DB, configData *configuration.ConfigData) {
	draftRepository := repository.NewDraftRepository(db)
	draftService := service.NewDraftService(draftRepository)
	draftController = controller.NewDraftController(draftService)
}
