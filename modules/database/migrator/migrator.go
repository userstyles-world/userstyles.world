// Package migrator provides functionality for migrating database schema.
package migrator

import (
	"time"

	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
)

// migration is a helper struct.
type migration struct {
	Version int
	Execute func(db *gorm.DB) error
	Name    string
}

// Migrate is the migration engine.
func Migrate() error {
	last, err := models.GetLastMigration(database.Conn)
	if err != nil {
		log.Database.Printf("Failed to find last migration: %v\n", err)
	}

	mx := migrations()
	if last.Version == 0 {
		if err := database.Conn.Transaction(func(tx *gorm.DB) error {
			// Check if migrations table already exists, then assume it already
			// contains the latest schema and insert the latest migration if so.
			if tx.Migrator().HasTable(models.Migration{}) {
				m := mx[len(mx)-1]
				defer func() { last.Version = m.Version }()
				return models.CreateMigration(tx, models.Migration{
					Version: m.Version,
					Name:    m.Name,
				})
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

			if err := database.Conn.Transaction(func(tx *gorm.DB) error {
				if err := m.Execute(tx); err != nil {
					return nil
				}

				return models.CreateMigration(tx, models.Migration{
					Version: m.Version,
					Name:    m.Name,
				})
			}); err != nil {
				return err
			}

			last.Version = m.Version // bump version
			log.Database.Printf("Migration done in %s.", time.Since(t))
		}
	}

	return nil
}

// migrations contains all schema migrations.
func migrations() []migration {
	m1 := func(db *gorm.DB) error {
		return database.Conn.AutoMigrate(models.Migration{})
	}

	return []migration{
		{1, m1, "add migrations table"},
	}
}
