package database

import (
	"log"

	"userstyles.world/config"
	"userstyles.world/models"
	"userstyles.world/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

type (
	user  models.User
	style models.Style
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
	Connect()

	// Generate data for development.
	if config.DB == "dev.db" {
		DropTables()
		defer Seed()
	}

	Migrate(&user{}, &style{})
}

func DropTables() {
	DB.Migrator().DropTable(&user{})
	DB.Migrator().DropTable(&style{})
}

func Seed() {
	users := []user{
		{
			Username: "vednoc",
			Email:    "vednoc@usw.local",
			Password: utils.GenerateHashedPassword("vednoc123"),
		},
		{
			Username: "john",
			Email:    "john@usw.local",
			Password: utils.GenerateHashedPassword("johnjohn"),
		},
		{
			Username: "jane",
			Email:    "jane@usw.local",
			Password: utils.GenerateHashedPassword("janejane"),
		},
	}

	styles := []style{
		{
			UserID:   1,
			Name:     "Dark-GitHub",
			Preview:  "https://user-images.githubusercontent.com/18245694/102144636-23044200-3e66-11eb-8d4b-e104de055f07.png",
			Featured: true,
		},
		{
			UserID:   1,
			Name:     "Dark-GitLab",
			Preview:  "https://user-images.githubusercontent.com/18245694/99352060-3d1c2600-28a2-11eb-8ab8-b22aea43f330.png",
			Featured: true,
		},
		{
			UserID:   1,
			Name:     "Dark-WhatsApp",
			Preview:  "https://user-images.githubusercontent.com/18245694/99352056-3c838f80-28a2-11eb-8bda-06a8c807dcb2.png",
			Featured: true,
		},
		{
			UserID:   2,
			Name:     "Archived userstyle",
			Archived: true,
		},
		{
			UserID:   3,
			Name:     "Featured userstyle",
			Featured: true,
		},
		{
			UserID: 3,
			Name:   "Temporary userstyle",
		},
	}

	for _, user := range users {
		DB.Create(&user)
	}
	for _, style := range styles {
		DB.Create(&style)
	}
}
