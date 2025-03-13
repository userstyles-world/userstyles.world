package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Migration represents schema for migrations table.
type Migration struct {
	Version   int       `gorm:"column:version"`
	Name      string    `gorm:"column:name"`
	CreatedAt time.Time `gorm:"column:applied_at"`
}

// GetLastMigration returns the last migration or an error if it's doesn't exit.
func GetLastMigration(db *gorm.DB) (m Migration, err error) {
	err = db.Last(&m).Error
	return m, err
}

// CreateMigration inserts a new migration and returns inserted data.
func CreateMigration(db *gorm.DB, m Migration) error {
	return db.Clauses(clause.Returning{}).Create(&m).Error
}
