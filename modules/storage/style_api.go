package storage

import (
	"bytes"
	"encoding/json"
	"strings"

	"gorm.io/gorm"

	"userstyles.world/modules/errors"
)

const (
	selectTotalInstalls  = "(SELECT total_installs FROM histories h WHERE h.style_id = styles.id ORDER BY id DESC LIMIT 1) AS TotalInstalls"
	selectWeeklyInstalls = "(SELECT weekly_installs FROM histories h WHERE h.style_id = styles.id ORDER BY id DESC LIMIT 1) AS WeeklyInstalls"
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
	var size int
	err := db.Raw("SELECT count(*) FROM styles").Scan(&size).Error
	if err != nil || size < 1 {
		return nil, errors.ErrStylesNotFound
	}

	var buf bytes.Buffer
	buf.WriteString(`{"data":`)

	styles := make([]StyleCompact, size)
	err = db.
		Select(selectCompactIndex).
		Where(notDeleted).
		Find(&styles).Error
	if err != nil {
		return nil, errors.ErrStylesNotFound
	}

	b, err := json.Marshal(&styles)
	if err != nil {
		return nil, err
	}
	buf.Write(b)

	buf.WriteRune('}')

	return buf.Bytes(), nil
}
