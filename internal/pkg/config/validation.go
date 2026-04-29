package config

import (
	xtrememdw "github.com/globalxtreme/go-core/v2/middleware"
	"github.com/go-playground/validator/v10"
)

func InitValidation() {
	v := xtrememdw.Validator{}
	v.RegisterValidation(func(validate *validator.Validate) {
		// Enter your custom validation rules
	})
}
