package storage

import "userstyles.world/modules/database"

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
