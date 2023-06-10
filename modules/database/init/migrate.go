package init

import (
	"strconv"
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
		var s models.Style
		if err := tx.Migrator().AddColumn(s, "ImportPrivate"); err != nil {
			log.Database.Fatalf("Failed to add column import_private: %s\n", err)
		}
		if err := tx.Migrator().AddColumn(s, "MirrorPrivate"); err != nil {
			log.Database.Fatalf("Failed to add column mirror_private: %s\n", err)
		}

		var styles []models.Style
		err := tx.Debug().Select("id, code").Find(&styles).Error
		if err != nil {
			log.Database.Fatalf("Failed to find styles: %s\n", err)
		}

		for _, style := range styles {
			err = models.SaveStyleCode(strconv.Itoa(int(style.ID)), style.Code)
			if err != nil {
				log.Warn.Fatalf("Failed to save %d: %s\n", style.ID, err)
			}
		}

		return nil
	})

	log.Database.Printf("Migration completed in %s.\n", time.Since(t))
}
