package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"userstyles.world/modules/log"
)

var (
	usernameRule    = `^[a-zA-Z0-9][a-zA-Z0-9-_]+[a-zA-Z0-9]$`
	displayNameRule = `^[a-zA-Z0-9-_ ]+$`
	v               = validator.New()
)

func InitializeValidator() {
	if err := v.RegisterValidation("username", usernameValidation); err != nil {
		log.Warn.Println("Cannot register username validation")
		panic(err)
	}

	if err := v.RegisterValidation("displayName", displayNameValidation); err != nil {
		log.Warn.Println("Cannot register displayName validation")
		panic(err)
	}

	v.RegisterAlias("Username", "username")
	v.RegisterAlias("DisplayName", "displayName")
}

func usernameValidation(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(usernameRule)
	valid := regex.Match([]byte(fl.Field().String()))
	return valid
}

func displayNameValidation(fl validator.FieldLevel) bool {
	regex := regexp.MustCompile(displayNameRule)
	valid := regex.Match([]byte(fl.Field().String()))
	return valid
}

func Validate() *validator.Validate {
	return v
}
