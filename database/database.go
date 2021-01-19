package database

import (
	"log"

	"userstyles.world/config"
	"userstyles.world/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Connect() {
	db, err := gorm.Open(sqlite.Open(config.DB), &gorm.Config{})

	if err != nil {
		log.Println("Failed to connect database.")
		panic(err)
	}

	DB = db
	log.Println("Database successfully connected.")
}

func Migrate(tables ...interface{}) error {
	log.Println("Migrating database tables.")
	return DB.AutoMigrate(tables...)
}

func Initialize() {
	type user models.User

	Connect()
	Migrate(&user{})
}
