package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gola-glitch/gola-utils/logging"
	"github.com/gola-glitch/gola-utils/mask_util"
	"net/http"
	"post-api/idp/constants"
	"post-api/idp/handlers/login"
	"post-api/idp/models/request"
	"post-api/idp/service"
)

type LoginController struct {
	service     service.LoginService
	authHandler login.OauthLoginHandler
}

func (controller LoginController) LoginByEmailAndPassword(ctx *gin.Context) {
	var loginRequest request.UserLoginRequest

	logger := logging.GetLogger(ctx).WithField("class", "LoginController").WithField("method", "LoginByEmailAndPassword")
	logger.Info("Initiating login with email and password")

	if bindingErr := ctx.ShouldBindBodyWith(&loginRequest, binding.JSON); bindingErr != nil {
		logger.Errorf("Error in binding login request %v", bindingErr)
		ctx.JSON(constants.GetGolaHttpCode(constants.PayloadValidationErrorCode), constants.PayloadValidationError)
		return
	}

	logger.Infof("Successfully bind request body for email %v", mask_util.MaskEmail(ctx, loginRequest.Email))

	response, err := controller.service.LoginWithEmailAndPassword(loginRequest, ctx)

	if err != nil {
		logger.Errorf("Error in login service %v", err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Infof("Login challenge successful")
	ctx.JSON(http.StatusOK, response)
}

func (controller LoginController) AcceptChallenge(ctx *gin.Context) {
	logger := logging.GetLogger(ctx)
	logger.Info("LoginController.AcceptChallenge: success")
	ctx.JSON(http.StatusOK, nil)
}

func (controller LoginController) GrantConsent(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "LoginController").WithField("method", "GrantConsent")
	logger.Info("Grant consent initiated")
	consentChallenge := ctx.Request.URL.Query()["consent_challenge"][0]

	response, err := controller.authHandler.AcceptConsentRequest(ctx, consentChallenge)
	if err != nil {
		logger.Errorf("Error is oauth handler %v", err)
		constants.RespondWithGolaError(ctx, err)
		return
	}

	logger.Info("consent request successful")
	ctx.JSON(http.StatusOK, response)
}

func (controller LoginController) ForgetPassword(ctx *gin.Context) {
	logger := logging.GetLogger(ctx).WithField("class", "LoginController").WithField("method", "ForgetPassword")
	logger.Info("Forget password initiated")
	forgetPassword := new(request.ForgetPassword)
	err := ctx.ShouldBindJSON(forgetPassword)
	if err != nil {
		logger.Errorf("unable to bind request body for forget password %v", err)
		constants.RespondWithGolaError(ctx, constants.PayloadValidationError)
		return
	}

	forgetPasswordErr := controller.service.ForgetPassword(ctx, *forgetPassword)
	if forgetPasswordErr != nil {
		logger.Errorf("unable to do action forget password %v", forgetPasswordErr)
		constants.RespondWithGolaError(ctx, forgetPasswordErr)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func NewLoginController(loginService service.LoginService, handler login.OauthLoginHandler) LoginController {
	return LoginController{
		service:     loginService,
		authHandler: handler,
	}
}
