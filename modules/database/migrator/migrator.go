// Package migrator provides functionality for migrating database schema.
package migrator

import (
	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
)

// migration is a helper struct.
type migration struct {
	Version int
	Execute func(db *gorm.DB) error
}

// Migrate is the migration engine.
func Migrate() error {
	last, err := models.GetLastMigration(database.Conn)
	if err != nil {
		log.Database.Printf("Failed to find last migration: %v\n", err)
	}

	// TODO: Handle all cases.
	if last.Version == 0 {
		if err := database.Conn.Transaction(func(tx *gorm.DB) error {
			log.Database.Printf("Adding migrations table.\n")
			return tx.AutoMigrate(models.Migration{})
		}); err != nil {
			return err
		}
	}

	for _, m := range migrations() {
		if m.Version > last.Version {
			if err := database.Conn.Transaction(func(tx *gorm.DB) error {
				if err := m.Execute(tx); err != nil {
					return nil
				}

				return models.CreateMigration(tx, models.Migration{
					Version: m.Version,
				})
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

// migrations contains all schema migrations.
func migrations() []migration {
	return []migration{
		{Version: 1, Execute: func(db *gorm.DB) error {
			return nil
		}},
	}
}
