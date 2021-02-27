package config

import (
	"fmt"
	"os"
)

var (
	PORT     = getEnv("PORT", ":3000")
	DB       = getEnv("DB", "dev.db")
	DB_DEBUG = getEnv("DB_DEBUG", "silent")
	DB_COLOR = getEnv("DB_COLOR", "false")
	DB_DROP  = getEnv("DB_DROP", "false")
	SALT     = getEnv("SALT", "10")
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
