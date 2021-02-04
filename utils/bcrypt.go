package utils

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func GenerateHashedPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), 14)
	if err != nil {
		log.Println(err)
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
