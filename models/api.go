package models

import (
	"time"

	"userstyles.world/modules/database"
	"userstyles.world/modules/errors"
)

type USoFormat struct {
	Username       string `json:"an"`
	Name           string `json:"n"`
	Category       string `json:"c"`
	Screenshot     string `json:"sn"`
	UpdatedAt      int64  `json:"u"`
	TotalInstalls  int64  `json:"t"`
	WeeklyInstalls int64  `json:"w"`
	ID             uint   `json:"i"`
}

type USoStyles []USoFormat

func (s *USoStyles) Query() error {
	lastWeek := time.Now().Add(-time.Hour * 24 * 7)

	stmt := "styles.id, styles.name, styles.user_id, u.username, "
	stmt += "styles.category, strftime('%s', styles.created_at) as created_at, "
	stmt += "strftime('%s', styles.updated_at) as updated_at, "
	stmt += "printf('https://userstyles.world/api/style/preview/%d.webp', styles.id) as screenshot, "
	stmt += "(select count(*) from stats where stats.style_id = styles.id) as total_installs, "
	stmt += "(select count(*) from stats where stats.style_id = styles.id and updated_at > ? and created_at < ?) as weekly_installs"

	err := database.Conn.
		Table("styles").
		Select(stmt, lastWeek, lastWeek).
		Joins("join users u on u.id = styles.user_id").
		Find(&s, "styles.deleted_at is null").
		Error
	if err != nil {
		return errors.ErrStylesNotFound
	}

	return nil
}
