package util

import (
	"golang.org/x/crypto/bcrypt"

	"userstyles.world/modules/config"
)

// HashPassword generates a hash out of a password.
func HashPassword(pw string, secrets *config.SecretsConfig) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), secrets.PasswordCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// VerifyPassword compares the hash to the password.
func VerifyPassword(hash, pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}
