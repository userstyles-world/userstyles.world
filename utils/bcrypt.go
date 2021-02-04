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
