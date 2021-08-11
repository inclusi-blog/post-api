package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func ValidPostUID(fl validator.FieldLevel) bool {
	if info, ok := fl.Field().Interface().(string); ok {
		_, err := uuid.Parse(info)
		if err != nil {
			return false
		}
		return true
	}
	return false
}
