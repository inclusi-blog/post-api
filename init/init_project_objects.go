package init

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"post-api/clients/user_profile"
	"post-api/configuration"
	"post-api/controller"
	"post-api/mapper"
	"post-api/repository"
	"post-api/service"
	"post-api/utils"
)

var (
	draftController     controller.DraftController
	interestsController controller.InterestsController
	postController      controller.PostController
)

func Objects(db neo4j.Session, configData *configuration.ConfigData) {
	userProfileClient := user_profile.NewClient(requestBuilder, configData)
	interestsMapper := mapper.NewInterestsMapper()
	draftRepository := repository.NewDraftRepository(db)
	interestsRepository := repository.NewInterestRepository(db)
	interestsService := service.NewInterestsService(interestsRepository, userProfileClient, interestsMapper)
	interestsController = controller.NewInterestsController(interestsService)
	postRepository := repository.NewPostsRepository(db)
	postValidator := utils.NewPostValidator(configData)
	draftService := service.NewDraftService(draftRepository, postValidator)
	draftController = controller.NewDraftController(draftService)
	postService := service.NewPostService(postRepository, draftRepository, postValidator, db)
	postController = controller.NewPostController(postService)
}
