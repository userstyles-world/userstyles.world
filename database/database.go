package database

import (
	"fmt"

	"userstyles.world/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Connect() {
	db, err := gorm.Open(sqlite.Open(config.DB), &gorm.Config{})

	if err != nil {
		fmt.Println("Failed to connect database.")
		panic(err)
	}

	DB = db
	fmt.Println("Database connected.")
}

func Migrate(tables ...interface{}) error {
	fmt.Println("Migrated database tables.")
	return DB.AutoMigrate(tables...)
}
