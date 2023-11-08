package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	golaConstants "github.com/inclusi-blog/gola-utils/middleware/introspection/oauth-middleware/constants"
	"github.com/inclusi-blog/gola-utils/model"
)

func GetIDToken(ctx *gin.Context) (model.IdToken, error) {
	token, exists := ctx.Get(golaConstants.ContextDecryptedIdTokenKey)
	if !exists {
		return model.IdToken{}, errors.New("id token not found")
	}
	idToken := token.(model.IdToken)
	return idToken, nil
}
