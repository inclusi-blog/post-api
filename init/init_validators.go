package init

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"post-api/validators"
)

func Validators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("validPostUID", validators.ValidPostUID)
	}
}
