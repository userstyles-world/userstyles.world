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
