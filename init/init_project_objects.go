package init

import (
	"github.com/jmoiron/sqlx"
	"post-api/configuration"
	"post-api/controller"
	"post-api/repository"
	"post-api/service"
)

var (
	draftController     controller.DraftController
	interestsController controller.InterestsController
)

func Objects(db *sqlx.DB, configData *configuration.ConfigData) {
	draftRepository := repository.NewDraftRepository(db)
	draftService := service.NewDraftService(draftRepository)
	draftController = controller.NewDraftController(draftService)
	interestsRepository := repository.NewInterestRepository(db)
	interestsService := service.NewInterestsService(interestsRepository)
	interestsController = controller.NewInterestsController(interestsService)
}
