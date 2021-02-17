package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"userstyles.world/config"
)

var (
	salt       = getSalt()
	MonitorURL = RandomString(64)
)

func getSalt() int {
	salt, err := strconv.Atoi(config.SALT)
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

func RandomString(size int) string {
	b := make([]byte, size)

	if _, err := rand.Read(b); err != nil {
		log.Fatalln("Failed to generate RandomString, err:", err)
	}

	return fmt.Sprintf("%X", b[0:size])
}
