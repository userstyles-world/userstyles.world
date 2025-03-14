// Package migrator provides functionality for migrating database schema.
package migrator

import (
	"time"

	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

// Migrate is the migration engine.
func Migrate(db *gorm.DB) error {
	last, err := models.GetLastMigration(db)
	if err != nil {
		log.Database.Printf("Failed to find last migration: %v\n", err)
	}

	mx := migrations()
	if last.Version == 0 {
		if err := db.Transaction(func(tx *gorm.DB) error {
			// Check if migrations table already exists, then assume it already
			// contains the latest schema and insert the latest migration if so.
			if tx.Migrator().HasTable(models.Migration{}) {
				m := mx[len(mx)-1]
				defer func() { last = m }()
				return models.CreateMigration(tx, m)
			}

			return nil
		}); err != nil {
			return err
		}
	}

	for _, m := range mx {
		if m.Version > last.Version {
			log.Database.Printf("Executing %q migration.", m.Name)
			t := time.Now()

			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := m.Execute(tx); err != nil {
					return nil
				}

				return models.CreateMigration(tx, m)
			}); err != nil {
				return err
			}

			last = m // update last migration
			log.Database.Printf("Migration done in %s.", time.Since(t))
		}
	}

	return nil
}

// migrations contains all schema migrations.
func migrations() []models.Migration {
	m1 := func(db *gorm.DB) error {
		return db.AutoMigrate(models.Migration{})
	}

	return []models.Migration{
		{Version: 1, Execute: m1, Name: "add migrations table"},
	}
}
