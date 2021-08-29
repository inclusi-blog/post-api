package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	utils "github.com/gola-glitch/gola-utils/logging"
	"net/http"
	"post-api/idp/constants"
	"post-api/idp/handlers/login"
	"post-api/idp/models/oauth"
)

type TokenController interface {
	ExchangeToken(ctx *gin.Context)
}

type tokenController struct {
	oauthService         login.OauthLoginHandler
	allowInsecureCookies bool
}

func NewTokenController(oauthService login.OauthLoginHandler, allowInsecureCookies bool) TokenController {
	return tokenController{oauthService: oauthService,
		allowInsecureCookies: allowInsecureCookies,
	}
}

// TODO Better way to set access tokens are cookies (need to check possibilites)
func (controller tokenController) ExchangeToken(ctx *gin.Context) {
	logger := utils.GetLogger(ctx)
	logger.Info("initiating exchange token")
	var request = oauth.TokenExchangeRequest{}
	if bindError := ctx.ShouldBindBodyWith(&request, binding.JSON); bindError != nil {
		logger.Errorf("Error in binding exchange token request body. Error: %-v", bindError)
		constants.RespondWithGolaError(ctx, &constants.PayloadValidationError)
		return
	}

	tokenExchangeResponse, _, err := controller.oauthService.ExchangeToken(ctx, request)
	if err != nil {
		logger.Errorf("Error in oauthService. Error: %-v", err)
		constants.RespondWithGolaError(ctx, err)
		return
	}
	logger.Info("setting access token and id token cookiesca.")

	clearCsrfCookies(ctx)
	if !controller.allowInsecureCookies {
		logger.Info("flag to allow insecure cookies is set to false")
		tokenExchangeResponse.IdToken = "dummy.jwt.value"
	}

	logger.Info("Token exchange completed")
	ctx.JSON(http.StatusOK, tokenExchangeResponse)
}

func clearCsrfCookies(ctx *gin.Context) {
	ctx.Writer.Header().Add("set-cookie", "oauth2_authentication_csrf=; Path=/; Max-Age=-1")
	ctx.Writer.Header().Add("set-cookie", "oauth2_authentication_session=; Path=/; Max-Age=-1")
	ctx.Writer.Header().Add("set-cookie", "oauth2_consent_csrf=; Path=/; Max-Age=-1")
}