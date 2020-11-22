package validators

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidPostUID(fl validator.FieldLevel) bool {
	if info, ok := fl.Field().Interface().(string); ok {
		matched, _ := regexp.MatchString("^[a-z0-9]{12}$", info)
		if matched {
			return true
		}
	}
	return false
}
