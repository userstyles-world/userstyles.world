package config

import (
	"fmt"
	"os"
	"strconv"
)

var envNotFound = "__NOT_FOUND__"

func getEnv(name, fallback string) string {
	if env, set := os.LookupEnv(name); set {
		return env
	}

	if fallback != "" {
		return fallback
	}

	panic(fmt.Sprintf(`Env variable not found: %v`, name))
}

func getEnvInt(name string, fallback int) int {
	env := getEnv(name, envNotFound)
	if env == envNotFound {
		return fallback
	}

	envInt, err := strconv.Atoi(env)
	if err != nil {
		return fallback
	}

	return envInt
}

func getEnvBool(name string, fallback bool) bool {
	env := getEnv(name, envNotFound)
	if env == envNotFound {
		return fallback
	}

	envBool, err := strconv.ParseBool(name)
	if err != nil {
		return fallback
	}

	return envBool
}
