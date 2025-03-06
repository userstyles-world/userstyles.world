package init

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

func migration(db *gorm.DB) error {
	f := []string{"slug", "code_size", "code_checksum"}
	for _, f := range []string{"Slug", "CodeSize", "CodeChecksum"} {
		if err := db.Migrator().AddColumn(models.Style{}, f); err != nil {
			return err
		}

		log.Database.Printf("Added '%s' column to 'styles' table.\n", f)
	}

	var sx []models.Style
	err := db.FindInBatches(&sx, 500, func(batch *gorm.DB, i int) error {
		for _, s := range sx {
			s.Prepare()
			if err := db.Model(s).Select(f).UpdateColumns(s).Error; err != nil {
				return err
			}

			fmt.Printf("%8d | %8d | %8s | %.40s\n", s.ID, s.CodeSize, s.CodeChecksum, s.Slug)
		}

		return nil
	}).Error

	return err
}

func runMigration(db *gorm.DB) {
	log.Database.Println("Migration started.")
	t := time.Now()

	db.Config.Logger = db.Config.Logger.LogMode(logger.Error)

	// Wrap in a transaction to allow rollbacks.
	if err := db.Transaction(migration); err != nil {
		log.Database.Fatalf("Migration error: %s\n", err)
	}

	log.Database.Printf("Done in %s.\n", time.Since(t).Round(time.Microsecond))
}
