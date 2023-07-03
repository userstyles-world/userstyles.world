package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"userstyles.world/modules/log"
)

var (
	usernameRule    = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9-_]+[a-zA-Z0-9]$`)
	displayNameRule = regexp.MustCompile(`^[a-zA-Z0-9-_ ]+$`)
	V               = validator.New()
)

// Init initializes validator and registers our validations.
func Init() {
	if err := V.RegisterValidation("username", usernameValidation); err != nil {
		log.Warn.Fatalf("Cannot register username validation: %s\n", err)
	}

	if err := V.RegisterValidation("displayName", displayNameValidation); err != nil {
		log.Warn.Fatalf("Cannot register displayName validation: %s\n", err)
	}

	V.RegisterAlias("Username", "username")
	V.RegisterAlias("DisplayName", "displayName")
}

func usernameValidation(fl validator.FieldLevel) bool {
	return usernameRule.Match([]byte(fl.Field().String()))
}

func displayNameValidation(fl validator.FieldLevel) bool {
	return displayNameRule.Match([]byte(fl.Field().String()))
}
