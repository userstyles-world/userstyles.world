package models

import (
	"userstyles.world/modules/errors"
)

type USoFormat struct {
	Username       string  `json:"an"`
	Name           string  `json:"n"`
	Category       string  `json:"c"`
	Screenshot     string  `json:"sn"`
	UpdatedAt      int64   `json:"u"`
	TotalInstalls  int64   `json:"t"`
	WeeklyInstalls int64   `json:"w"`
	Rating         float64 `json:"r"`
	ID             uint    `json:"i"`
}

type USoStyles []USoFormat

func (s *USoStyles) Query() error {
	stmt := "styles.id, styles.name, styles.user_id, styles.category, "
	stmt += "STRFTIME('%s', styles.updated_at) AS updated_at, "
	stmt += "PRINTF('https://userstyles.world/api/style/preview/%d.webp', styles.id) AS screenshot, "
	stmt += "(SELECT username FROM users u WHERE u.id = styles.user_id) AS username, "
	stmt += "(SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = styles.id AND r.deleted_at IS NULL) AS rating, "
	stmt += "(SELECT count(*) FROM stats s WHERE s.style_id = styles.id AND s.install > DATETIME('now', '-7 days')) AS TotalInstalls, "
	stmt += "(SELECT count(*) FROM stats s WHERE s.style_id = styles.id AND s.install > DATETIME('now', '-7 days') AND s.created_at > DATETIME('now', '-7 days')) AS WeeklyInstalls"

	err := db().Model(modelStyle).Select(stmt).Find(&s).Error
	if err != nil {
		return errors.ErrStylesNotFound
	}

	return nil
}
