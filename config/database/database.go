package database

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

var (
	// File holds the database name.
	File = getEnv("DB", "dev.db")

	// Conn holds the pointer to database connections.
	Conn *gorm.DB
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
