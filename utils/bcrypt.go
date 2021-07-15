package utils

import (
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"

	"userstyles.world/modules/config"
)

var salt = getSalt()

func getSalt() int {
	salt, err := strconv.Atoi(config.Salt)
	if err != nil {
		log.Fatalln("Failed to convert SALT env variable, err:", err)
	}

	return salt
}

func GenerateHashedPassword(pw string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), salt)
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
