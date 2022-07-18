package storage

import (
	"strings"

	"gorm.io/gorm"

	"userstyles.world/modules/database"
	"userstyles.world/modules/errors"
)

var (
	selectWeeklyInstalls = "(SELECT COUNT(*) FROM stats s WHERE s.style_id = styles.id AND s.install > DATETIME('now', '-7 days') AND s.created_at > DATETIME('now', '-7 days')) AS WeeklyInstalls"
	selectTotalInstalls  = "(SELECT COUNT(*) FROM stats s WHERE s.style_id = styles.id AND s.install > 0) AS TotalInstalls"

	selectCompactIndex = strings.Join([]string{
		"styles.id", "styles.name", "styles.preview", "styles.category",
		"STRFTIME('%s', styles.updated_at) AS updated_at",
		selectAuthor, selectTotalInstalls, selectWeeklyInstalls, selectRatings,
	}, ", ")
)

// StyleCompact is a field-aligned struct optimized for compact style index.
type StyleCompact struct {
	Name           string  `json:"n"`
	Username       string  `json:"an"`
	Preview        string  `json:"sn"`
	Category       string  `json:"c"`
	ID             int     `json:"i"`
	UpdatedAt      int     `json:"u"`
	TotalInstalls  int     `json:"t"`
	WeeklyInstalls int     `json:"w"`
	Rating         float64 `json:"r"`
}

// TableName returns which table in database to use with GORM.
func (StyleCompact) TableName() string { return "styles" }

// GetStyleCompactIndex returns a compact JSON for our integration with Stylus.
func GetStyleCompactIndex() ([]StyleCompact, error) {
	var styles, batch []StyleCompact
	action := func(*gorm.DB, int) error {
		styles = append(styles, batch...)
		return nil
	}

	err := database.Conn.
		Select(selectCompactIndex).
		Where(notDeleted).
		FindInBatches(&batch, 250, action).Error
	if err != nil {
		return nil, errors.ErrStylesNotFound
	}

	return styles, nil
}
