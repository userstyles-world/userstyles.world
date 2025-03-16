package migrator

import (
	"strings"

	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

// deindent removes unneeded control characters from string literals.
func deindent(s string) string {
	s = strings.TrimPrefix(s, "\n")
	s = strings.ReplaceAll(s, "\t\t", "    ")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

func m1(db *gorm.DB) error {
	const q = `
	CREATE TABLE migrations(
		version INTEGER,
		name TEXT,
		applied_at DATETIME
	);`
	return db.Exec(deindent(q)).Error
}

func m2(db *gorm.DB) error {
	cx := []string{"slug", "code_size", "code_checksum"}

	var s models.Style
	for _, c := range cx {
		if !db.Migrator().HasColumn(s, c) {
			if err := db.Migrator().AddColumn(s, c); err != nil {
				return err
			}
			log.Database.Printf("Added '%s' column to 'styles' table.\n", c)
		}
	}

	var sx []models.Style
	return db.FindInBatches(&sx, 500, func(tx *gorm.DB, size int) error {
		for i, s := range sx {
			s.Prepare()
			if err := tx.Select(cx).UpdateColumns(s).Error; err != nil {
				return err
			}

			// Print everything in debug mode or every 100th message.
			if config.App.Debug || i%100 == 0 {
				const f = "id=%-8d cs=%-8d cc=%-8s s=%.40s\n"
				log.Database.Printf(f, s.ID, s.CodeSize, s.CodeChecksum, s.Slug)
			}
		}

		return nil
	}).Error
}
