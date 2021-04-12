package config

import (
	"fmt"
	"os"
)

var (
	PORT                   = getEnv("PORT", ":3000")
	DB                     = getEnv("DB", "dev.db")
	DB_DEBUG               = getEnv("DB_DEBUG", "silent")
	DB_COLOR               = getEnv("DB_COLOR", "false")
	DB_DROP                = getEnv("DB_DROP", "false")
	DB_RANDOM_DATA         = getEnv("DB_RANDOM_DATA", "false")
	SALT                   = getEnv("SALT", "10")
	JWT_SIGNING_KEY        = getEnv("JWT_SIGNING_KEY", "ABigSecretPassword")
	VERIFY_JWT_SIGNING_KEY = getEnv("VERIFY_JWT_SIGNING_KEY", "OhNoWeCantUseTheSameAsJWTBeCaUseSeCuRiTy1337")
	CRYPTO_KEY             = getEnv("CRYPTO_KEY", "ABigSecretPasswordWhichIsExact32")
	STATS_KEY              = getEnv("STATS_KEY", "KeyUsedForHashingStatistics")
	EMAIL_ADDRESS          = getEnv("EMAIL_ADDRESS", "test@userstyles.world")
	EMAIL_PWD              = getEnv("EMAIL_PWD", "hahah_not_your_password")
	GIT_COMMIT             = getEnv("GIT_COMMIT", "unset")
)

func getEnv(name, fallback string) string {
	if val, set := os.LookupEnv(name); set {
		return val
	}

	if fallback != "" {
		return fallback
	}

	panic(fmt.Sprintf(`Env variable not found: %v`, name))
}
