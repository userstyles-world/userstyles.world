package config

import (
	"fmt"
	"os"
)

var (
	PORT = getEnv("PORT", ":3000")
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
