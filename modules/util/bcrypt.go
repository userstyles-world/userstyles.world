package util

import (
	"golang.org/x/crypto/bcrypt"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

var salt = config.Salt

func GenerateHashedPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), salt)
	if err != nil {
		log.Warn.Println(err)
	}

	return string(hash)
}

func CompareHashedPassword(user, form string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user), []byte(form))
	if err != nil {
		return err
	}

	return nil
}
