package storage

import (
	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/database"
)

var isMirrored = "(styles.mirror_url <> '' or styles.original <> '') and styles.mirror_code = 1"

func CountStylesForUserID(id uint) (int, error) {
	var i int

	stmt := "SELECT COUNT(*) FROM styles WHERE user_id = ? AND " + notDeleted
	tx := database.Conn.Raw(stmt, id).Scan(&i)
	if err := tx.Error; err != nil {
		return 0, err
	}

	return i, nil
}

// FindStyleCode is an optimized query that returns userstyle's source code.
func FindStyleCode(id int) (string, error) {
	var code string

	stmt := "SELECT code FROM styles WHERE id = ? AND " + notDeleted
	tx := database.Conn.Raw(stmt, id).Scan(&code)
	if err := tx.Error; err != nil {
		return "", err
	}

	return code, nil
}

// CountStylesForMirror returns the number of styles with enabled mirroring.
func CountStylesForMirror() (int, error) {
	var i int

	stmt := "SELECT COUNT(*) FROM styles WHERE "
	stmt += isMirrored + " AND " + notDeleted
	tx := database.Conn.Raw(stmt).Scan(&i)
	if err := tx.Error; err != nil {
		return 0, err
	}

	return i, nil
}

// FindStylesForMirror queries for styles with enabled mirroring.
func FindStylesForMirror(action func([]models.Style) error) error {
	var styles []models.Style
	return database.Conn.Where(isMirrored).
		FindInBatches(&styles, 25, func(tx *gorm.DB, size int) error {
			return action(styles)
		}).Error
}
