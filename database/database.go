package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"userstyles.world/config"
	"userstyles.world/models"
	"userstyles.world/utils"
)

var (
	DB    *gorm.DB
	user  models.User
	style models.Style
)

func Connect() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      utils.DatabaseLogLevel(DB),
			Colorful:      utils.DatabaseColorful(DB),
		},
	)

	db, err := gorm.Open(sqlite.Open(config.DB), &gorm.Config{
		Logger: newLogger,
	})

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
	if utils.DatabaseDropTables(DB) && config.DB == "dev.db" {
		log.Println("Dropping database tables.")
		Drop(&user, &style)
		defer Seed()
	}

	Migrate(&user, &style)
}

func Drop(dst ...interface{}) error {
	return DB.Migrator().DropTable(dst...)
}

func Seed() {
	users := []models.User{
		{
			Username:  "vednoc",
			Email:     "vednoc@usw.local",
			Biography: "Something goes here.",
			Password:  utils.GenerateHashedPassword("vednoc123"),
		},
		{
			Username:  "john",
			Email:     "john@usw.local",
			Biography: "Something.",
			Password:  utils.GenerateHashedPassword("johnjohn"),
		},
		{
			Username: "jane",
			Email:    "jane@usw.local",
			Password: utils.GenerateHashedPassword("janejane"),
		},
	}

	styles := []models.Style{
		{
			UserID:   1,
			Name:     "Dark-GitHub",
			Summary:  "Customizable dark theme for GitHub.",
			Preview:  "https://user-images.githubusercontent.com/18245694/102033688-57232880-3dbc-11eb-8131-2eb21239160d.png",
			Code:     "https://raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl",
			Homepage: "https://github.com/vednoc/dark-github",
			Category: "github.com",
			Featured: true,
		},
		{
			UserID:   1,
			Name:     "Dark-GitLab",
			Summary:  "Customizable dark theme for GitLab.",
			Preview:  "https://gitlab.com/vednoc/dark-gitlab/-/raw/master/images/preview.png",
			Code:     "https://gitlab.com/vednoc/dark-gitlab/raw/master/gitlab.user.styl",
			Homepage: "https://gitlab.com/vednoc/dark-gitlab",
			Category: "gitlab.com",
			Featured: true,
		},
		{
			UserID:   1,
			Name:     "Dark-WhatsApp",
			Summary:  "Customizable dark theme for WhatsApp.",
			Preview:  "https://raw.githubusercontent.com/vednoc/dark-whatsapp/master/images/preview.png",
			Code:     "https://raw.githubusercontent.com/vednoc/dark-whatsapp/master/wa.user.styl",
			Homepage: "https://github.com/vednoc/dark-whatsapp",
			Category: "web.whatsapp.com",
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
