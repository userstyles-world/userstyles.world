package init

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
	"userstyles.world/modules/database"
	"userstyles.world/utils"
)

var tables = []struct {
	name  string
	model interface{}
}{
	{"users", &models.User{}},
	{"styles", &models.Style{}},
	{"stats", &models.Stats{}},
	{"oauths", &models.OAuth{}},
	{"histories", &models.History{}},
}

func connect() (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel(),
			Colorful:      colorful(),
		},
	)

	conn, err := gorm.Open(sqlite.Open(config.DB), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Initialize the database connection.
func Initialize() {
	conn, err := connect()
	if err != nil {
		log.Fatalf("Failed to connect database, err: %s\n", err.Error())
	}

	database.Conn = conn
	log.Println("Database successfully connected.")

	// Generate data for development.
	if dropTables() && !config.Production {
		for _, table := range tables {
			if err := drop(table.model); err != nil {
				log.Fatalf("Failed to drop %s, err: %s", table.name, err.Error())
			}
			log.Printf("Dropped database table %s.\n", table.name)
		}
		defer seed()
	}

	// Migrate tables.
	for _, table := range tables {
		if err := migrate(table.model); err != nil {
			log.Fatalf("Failed to migrate %s, err: %s", table.name, err.Error())
		}
		log.Printf("Migrated database table %s.\n", table.name)
	}
}

func generateData(amount int) ([]models.Style, []models.User) {
	randomData := utils.UnsafeString(utils.RandStringBytesMaskImprSrcUnsafe(amount * 7 * 4))
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
	randomData = utils.UnsafeString(utils.RandStringBytesMaskImprSrcUnsafe(amount * 4 * 4))
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

func seed() {
	log.Println("Seeding database mock data.")
	defer log.Println("Finished seeding mock data.")

	users := []models.User{
		{
			Username:  "admin",
			Email:     "admin@usw.local",
			Biography: "Admin of USw.",
			Password:  utils.GenerateHashedPassword("admin123"),
			Role:      models.Admin,
		},
		{
			Username:  "moderator",
			Email:     "moderator@usw.local",
			Biography: "I'm a moderator.",
			Password:  utils.GenerateHashedPassword("moderator"),
		},
		{
			Username: "regular",
			Email:    "regular@usw.local",
			Password: utils.GenerateHashedPassword("regular"),
		},
	}

	styles := []models.Style{
		{
			UserID:      1,
			Name:        "Dark-GitHub",
			Description: "Customizable dark theme for GitHub.",
			Notes:       "Some notes go here.",
			Preview:     "https://userstyles.world/api/style/preview/2.webp",
			Original:    "https://raw.githubusercontent.com/vednoc/dark-github/main/github.user.styl",
			Homepage:    "https://github.com/vednoc/dark-github",
			Category:    "github.com",
			MirrorCode:  true,
			Featured:    true,
		},
		{
			UserID:      1,
			Name:        "Dark-GitLab",
			Description: "Customizable dark theme for GitLab.",
			Notes:       "Some notes go here.",
			Preview:     "https://userstyles.world/api/style/preview/3.webp",
			Original:    "https://gitlab.com/vednoc/dark-gitlab/raw/master/gitlab.user.styl",
			Homepage:    "https://gitlab.com/vednoc/dark-gitlab",
			Category:    "gitlab.com",
			MirrorCode:  true,
			Featured:    true,
		},
		{
			UserID:      1,
			Name:        "Dark-WhatsApp",
			Description: "Customizable dark theme for WhatsApp.",
			Notes:       "Some notes go here.",
			Preview:     "https://userstyles.world/api/style/preview/4.webp",
			Original:    "https://raw.githubusercontent.com/vednoc/dark-whatsapp/master/wa.user.styl",
			Homepage:    "https://github.com/vednoc/dark-whatsapp",
			Category:    "web.whatsapp.com",
			MirrorCode:  true,
			Featured:    true,
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

	OAuths := []models.OAuth{
		{
			UserID:       1,
			Name:         "USw integration",
			Description:  "Just some integration",
			Scopes:       []string{"1", "user"},
			ClientID:     "publicccc_client",
			ClientSecret: "secreettUwU",
			RedirectURI:  "https://gusted.xyz/callback_helper",
		},
	}

	if config.DB_RANDOM_DATA != "false" {
		amount, _ := strconv.Atoi(config.DB_RANDOM_DATA)
		s, u := generateData(amount)
		styles = append(styles, s...)
		users = append(users, u...)
	}

	for i := range users {
		database.Conn.Create(&users[i])
	}
	for i := range styles {
		database.Conn.Create(&styles[i])
	}
	for i := range OAuths {
		database.Conn.Create(&OAuths[i])
	}
}
