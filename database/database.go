package database

import (
	"log"
	"os"
	"strconv"
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
	stats models.Stats
)

func Connect() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      utils.DatabaseLogLevel(),
			Colorful:      utils.DatabaseColorful(),
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
	if utils.DatabaseDropTables() && config.DB == "dev.db" {
		log.Println("Dropping database tables.")
		Drop(&user, &style, &stats)
		defer Seed()
	}

	// Split mirror into code and meta. TODO: Remove once it ran once in production.
	DB.Migrator().RenameColumn(models.Style{}, "mirror", "mirror_code")
	DB.Migrator().DropColumn(models.Style{}, "mirror")

	Migrate(&user, &style, &stats)
}

func Drop(dst ...interface{}) error {
	return DB.Migrator().DropTable(dst...)
}

func generateData(amount int) ([]models.Style, []models.User) {
	randomData := utils.B2s(utils.RandStringBytesMaskImprSrcUnsafe(amount * 7 * 4))
	var styleStructs []models.Style
	for i := 0; i < amount; i++ {
		startData := randomData[(i * 7 * 4):]
		styleStructs = append(styleStructs, models.Style{
			UserID:      uint(amount),
			Category:    startData[:4],
			Name:        startData[4:8],
			Description: startData[8:12],
			Notes:       startData[12:16],
			Preview:     startData[16:20],
			Code:        startData[20:24],
			Homepage:    startData[24:28],
			Featured:    true,
		})
	}

	var userStructs []models.User
	randomData = utils.B2s(utils.RandStringBytesMaskImprSrcUnsafe(amount * 4 * 4))
	for i := 0; i < amount; i++ {
		startData := randomData[(i * 4 * 4):]
		userStructs = append(userStructs, models.User{
			Username:  startData[:4],
			Email:     startData[4:8],
			Biography: startData[8:12],
			Password:  startData[12:16],
		})
	}

	return styleStructs, userStructs
}

func Seed() {
	users := []models.User{
		{
			Username:  "vednoc",
			Email:     "vednoc@usw.local",
			Biography: "Something goes here.",
			Password:  utils.GenerateHashedPassword("vednoc123"),
			Role:      models.Admin,
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
			UserID:      1,
			Name:        "Dark-GitHub",
			Description: "Customizable dark theme for GitHub.",
			Notes:       "Some notes go here.",
			Preview:     "https://user-images.githubusercontent.com/18245694/102033688-57232880-3dbc-11eb-8131-2eb21239160d.png",
			Code:        "https://raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl",
			Homepage:    "https://github.com/vednoc/dark-github",
			Category:    "github.com",
			Featured:    true,
		},
		{
			UserID:      1,
			Name:        "Dark-GitLab",
			Description: "Customizable dark theme for GitLab.",
			Notes:       "Some notes go here.",
			Preview:     "https://gitlab.com/vednoc/dark-gitlab/-/raw/master/images/preview.png",
			Code:        "https://gitlab.com/vednoc/dark-gitlab/raw/master/gitlab.user.styl",
			Homepage:    "https://gitlab.com/vednoc/dark-gitlab",
			Category:    "gitlab.com",
			Featured:    true,
		},
		{
			UserID:      1,
			Name:        "Dark-WhatsApp",
			Description: "Customizable dark theme for WhatsApp.",
			Notes:       "Some notes go here.",
			Preview:     "https://raw.githubusercontent.com/vednoc/dark-whatsapp/master/images/preview.png",
			Code:        "/* ==UserStyle==\n@name           Example\n@namespace      example.com\n@version        1.0.0\n@description    A new userstyle\n@author         Me\n==/UserStyle== */\n@-moz-document domain('example.com') {\n/** Your code goes here! */\n}",
			Homepage:    "https://github.com/vednoc/dark-whatsapp",
			Category:    "web.whatsapp.com",
			Featured:    true,
			Original:    "https://raw.githubusercontent.com/vednoc/dark-whatsapp/master/wa.user.styl",
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

	if config.DB_RANDOM_DATA != "false" {
		amount, _ := strconv.Atoi(config.DB_RANDOM_DATA)
		s, u := generateData(amount)
		styles = append(styles, s...)
		users = append(users, u...)
	}

	for _, user := range users {
		DB.Create(&user)
	}
	for _, style := range styles {
		DB.Create(&style)
	}
}
