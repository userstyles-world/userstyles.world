package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type deletedAt = gorm.DeletedAt

// Changelog represents schema for changes table.
type Changelog struct {
	ID          int       `gorm:"column:id; primaryKey"`
	UserID      int       `gorm:"column:user_id"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	DeletedAt   deletedAt `gorm:"column:deleted_at; index"`
	Title       string    `gorm:"column:title" form:"title"`
	Description string    `gorm:"column:description" form:"description"`
	User        *User
}

// GetChangelogs tries to return changelogs sorted from most to least recent.
func GetChangelogs(db *gorm.DB) (clx []Changelog, err error) {
	err = db.Order("id DESC").Find(&clx).Error
	return clx, err
}

// CreateChangelog tries to insert a new changelog and return inserted data.
func CreateChangelog(db *gorm.DB, cl Changelog) error {
	return db.Clauses(clause.Returning{}).Create(&cl).Error
}
