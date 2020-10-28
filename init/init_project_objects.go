package init

import (
	"github.com/jmoiron/sqlx"
	"post-api/configuration"
	"post-api/controller"
	"post-api/repository"
	"post-api/service"
	"post-api/utils"
)

var (
	draftController     controller.DraftController
	interestsController controller.InterestsController
	postController      controller.PostController
)

func Objects(db *sqlx.DB, configData *configuration.ConfigData) {
	draftRepository := repository.NewDraftRepository(db)
	draftService := service.NewDraftService(draftRepository)
	draftController = controller.NewDraftController(draftService)
	interestsRepository := repository.NewInterestRepository(db)
	interestsService := service.NewInterestsService(interestsRepository)
	interestsController = controller.NewInterestsController(interestsService)
	postRepository := repository.NewPostsRepository(db)
	postValidator := utils.NewPostValidator(configData)
	previewPostRepository := repository.NewPreviewPostsRepository(db)
	postService := service.NewPostService(postRepository, draftRepository, postValidator, previewPostRepository)
	postController = controller.NewPostController(postService)
}
