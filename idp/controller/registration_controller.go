package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"net/http"
	"post-api/idp/constants"
	"post-api/idp/models/request"
	"post-api/idp/service"
)

type RegistrationController struct {
	registrationCacheService service.RegistrationCacheService
	userRegistrationService  service.UserRegistrationService
}

// NewRegistration godoc
// @Tags registration
// @Summary NewRegistration
// @Description initiate new user registration request
// @Accept json
// @Param request body request.NewRegistrationRequest true "Request Body"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/idp/v1/user/register [post]
func (controller RegistrationController) NewRegistration(ctx *gin.Context) {
	log := logging.GetLogger(ctx).WithField("class", "RegistrationController").WithField("method", "NewRegistration")

	var registrationRequest request.NewRegistrationRequest

	err := ctx.ShouldBindBodyWith(&registrationRequest, binding.JSON)

	if err != nil {
		log.Errorf("Unable to bind request body for new registration request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Successfully bind request body %v", registrationRequest)

	cacheSaveErr := controller.registrationCacheService.SaveUserDetailsInCache(registrationRequest, ctx)

	if cacheSaveErr != nil {
		log.Errorf("Unable to save user details in cache for request payload . %v", cacheSaveErr)
		constants.RespondWithGolaError(ctx, cacheSaveErr)
		return
	}

	ctx.Status(http.StatusOK)
}

// ActivateUser godoc
// @Tags registration
// @Summary ActivateUser
// @Description activates the account by clicking on uuid link sent to email
// @Accept json
// @Param activation_hash path string true "activation_hash"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /api/idp/v1/user/activate/{activation_hash} [get]
func (controller RegistrationController) ActivateUser(ctx *gin.Context) {
	log := logging.GetLogger(ctx).WithField("class", "RegistrationController").WithField("method", "ActivateUser")

	var initiateRegistrationRequest request.InitiateRegistrationRequest

	err := ctx.ShouldBindBodyWith(&initiateRegistrationRequest, binding.JSON)

	if err != nil {
		log.Errorf("Unable to bind request body for new registration request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Successfully bind request body %v", initiateRegistrationRequest)

	userSaveErr := controller.userRegistrationService.InitiateRegistration(initiateRegistrationRequest, ctx)

	if userSaveErr != nil {
		log.Errorf("Unable to save user details in cache for request payload . %v", userSaveErr)
		constants.RespondWithGolaError(ctx, userSaveErr)
		return
	}

	ctx.Status(http.StatusOK)
}

func (controller RegistrationController) IsEmailAvailable(ctx *gin.Context) {
	log := logging.GetLogger(ctx).WithField("class", "RegistrationController").WithField("method", "IsEmailAvailable")

	var emailAvailabilityRequest request.EmailAvailabilityRequest

	if err := ctx.ShouldBindBodyWith(&emailAvailabilityRequest, binding.JSON); err != nil {
		log.Errorf("Unable to bind request body for email availability request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Successfully bind request body %v", emailAvailabilityRequest)
	emailAvailabilityResponse, err := controller.userRegistrationService.IsEmailRegistered(emailAvailabilityRequest, ctx)

	if err != nil {
		log.Errorf("Error occurred while fetching user existence from service for emaik %v. %v", emailAvailabilityRequest.Email, err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, emailAvailabilityResponse)
}

func (controller RegistrationController) IsUsernameAvailable(ctx *gin.Context) {
	log := logging.GetLogger(ctx).WithField("class", "RegistrationController").WithField("method", "IsEmailAvailable")

	var availabilityRequest request.UsernameAvailabilityRequest

	if err := ctx.ShouldBindBodyWith(&availabilityRequest, binding.JSON); err != nil {
		log.Errorf("Unable to bind request body for username availability request %v", err)
		ctx.JSON(http.StatusBadRequest, constants.PayloadValidationError)
		return
	}

	log.Infof("Successfully bind request body %v", availabilityRequest)
	usernameAvailabilityResponse, err := controller.userRegistrationService.IsUsernameRegistered(availabilityRequest, ctx)

	if err != nil {
		log.Errorf("Error occurred while fetching username existence from service for username %v. %v", availabilityRequest.Username, err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, usernameAvailabilityResponse)
}

func NewRegistrationController(cacheService service.RegistrationCacheService, registrationService service.UserRegistrationService) RegistrationController {
	return RegistrationController{
		registrationCacheService: cacheService,
		userRegistrationService:  registrationService,
	}
}
