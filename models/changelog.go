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

// GetChangelog tries to return a changelog with a given id.
func GetChangelog(db *gorm.DB, id int) (cl Changelog, err error) {
	err = db.Find(&cl, id).Error
	return cl, err
}

// CreateChangelog tries to insert a new changelog and return inserted data.
func CreateChangelog(db *gorm.DB, cl Changelog) error {
	return db.Clauses(clause.Returning{}).Create(&cl).Error
}

// UpdateChangelog tries to update an existing changelog.
func UpdateChangelog(db *gorm.DB, cl Changelog) error {
	return db.Updates(cl).Error
}

// DeleteChangelog tries to delete an existing changelog.
func DeleteChangelog(db *gorm.DB, cl Changelog) error {
	return db.Delete(&cl).Error
}
