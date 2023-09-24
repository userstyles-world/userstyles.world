package init

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

func runMigration(db *gorm.DB) {
	log.Database.Println("Migration started.")
	t := time.Now()

	db.Config.Logger = db.Config.Logger.LogMode(logger.Info)

	// Wrap in a transaction to allow rollbacks.
	db.Transaction(func(tx *gorm.DB) error {
		return models.InitStyleSearch()
	})

	log.Database.Printf("Done in %s.\n", time.Since(t).Round(time.Microsecond))
}
