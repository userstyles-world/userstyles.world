package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	usernameRule = `^[a-zA-Z0-9_]+$`
	v            = validator.New()
)

func usernameValidation(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(usernameRule)
	valid := regex.Match([]byte(fl.Field().String()))

	return valid
}

func registerValidations(v *validator.Validate) *validator.Validate {
	v.RegisterValidation("username", usernameValidation)
	return v
}

func Validate() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("username", usernameValidation)

	return v
}
