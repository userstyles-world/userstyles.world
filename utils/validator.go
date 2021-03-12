package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	usernameRule = `^[a-zA-Z0-9_]+$`
	v            = validator.New()
)

func InitializeValidator() {
	v.RegisterValidation("username", usernameValidation)
	v.RegisterAlias("Username", "username")
}

func usernameValidation(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(usernameRule)
	valid := regex.Match([]byte(fl.Field().String()))
	return valid
}

func Validate() *validator.Validate {
	return v
}
