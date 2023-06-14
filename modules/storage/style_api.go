package storage

import (
	"bytes"
	"encoding/json"
	"strings"
	"time"

	"gorm.io/gorm"

	"userstyles.world/modules/errors"
)

const (
	selectWeeklyInstalls = "(SELECT COUNT(*) FROM stats s WHERE s.style_id = styles.id AND s.install > DATETIME('now', '-7 days') AND s.created_at > DATETIME('now', '-7 days')) AS WeeklyInstalls"
	selectTotalInstalls  = "(SELECT COUNT(*) FROM stats s WHERE s.style_id = styles.id AND s.install > 0) AS TotalInstalls"
	selectUSoRatings     = "(SELECT ROUND(AVG(rating)*0.6, 1) FROM reviews r WHERE r.style_id = styles.id AND rating > 0 AND r.deleted_at IS NULL) AS Rating"
)

var (
	selectCompactIndex = strings.Join([]string{
		"styles.id", "styles.name", "styles.preview", "styles.category",
		"STRFTIME('%s', styles.updated_at) AS updated_at",
		selectAuthor, selectTotalInstalls, selectWeeklyInstalls, selectUSoRatings,
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

// GetStyleCompactIndex returns a compact index for our integration with Stylus.
func GetStyleCompactIndex(db *gorm.DB) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(`{"data":`)

	var styles []StyleCompact
	action := func(tx *gorm.DB, batch int) error {
		b, err := json.Marshal(&styles)
		if err != nil {
			return err
		}

		if batch == 1 {
			buf.Write(b[:len(b)-1])
		} else {
			buf.WriteRune(',')
			buf.Write(b[1 : len(b)-1])
		}

		time.Sleep(100 * time.Millisecond)
		return nil
	}

	err := db.
		Select(selectCompactIndex).
		Where(notDeleted).
		FindInBatches(&styles, 100, action).Error
	if err != nil {
		return nil, errors.ErrStylesNotFound
	}

	buf.WriteString("]}")

	return buf.Bytes(), nil
}
