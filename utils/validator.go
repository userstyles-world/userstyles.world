package utils

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	usernameRule    = `^[a-zA-Z0-9_]+$`
	displayNameRule = `^[a-zA-Z0-9-_ ]+$`
	v               = validator.New()
)

func InitializeValidator() {
	if err := v.RegisterValidation("username", usernameValidation); err != nil {
		log.Println("Cannot register username validation")
		panic(err)
	}

	if err := v.RegisterValidation("displayName", usernameValidation); err != nil {
		log.Println("Cannot register displayName validation")
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
