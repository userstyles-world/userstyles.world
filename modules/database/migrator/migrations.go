package migrator

import (
	"fmt"

	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

func m1(db *gorm.DB) error {
	return db.AutoMigrate(models.Migration{})
}

func m2(db *gorm.DB) error {
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
