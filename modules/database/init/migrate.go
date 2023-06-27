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
		var l models.Log
		if err := tx.Migrator().AddColumn(l, "Message"); err != nil {
			log.Database.Fatalf("Failed to add column message: %s\n", err)
		}

		return nil
	})

	log.Database.Printf("Done in %s.\n", time.Since(t).Round(time.Microsecond))
}
