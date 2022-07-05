package storage

import (
	"gorm.io/gorm"

	"userstyles.world/modules/database"
)

// StyleSearch is a field-aligned struct optimized for style search.
type StyleSearch struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Notes       string `json:"notes"`
	ID          int    `json:"id"`
}

// FindStylesForSearch queries for styles in batches, and runs a passed action
// that might return an error.
func FindStylesForSearch(action func([]StyleSearch) error) error {
	var res []StyleSearch
	fn := func(tx *gorm.DB, batch int) error {
		return action(res)
	}

	return database.Conn.
		Table("styles").Select("id, name, description, notes", selectAuthor).
		Where("deleted_at IS NULL").FindInBatches(&res, 250, fn).Error
}
