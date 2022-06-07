package config

import (
	"log"
	"os"
	"strconv"
)

func getEnv(name, fallback string) string {
	env, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	return env
}

func getEnvInt(name string, fallback int) int {
	env, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	res, err := strconv.Atoi(env)
	if err != nil {
		log.Fatalf("Failed to convert %q to an int for %q.\n", env, name)
	}

	return res
}

func getEnvBool(name string, fallback bool) bool {
	env, ok := os.LookupEnv(name)
	if !ok {
		return fallback
	}

	res, err := strconv.ParseBool(env)
	if err != nil {
		log.Fatalf("Failed to convert %q to a bool for %q.\n", env, name)
	}

	return res
}
