package storage

import (
	"strconv"
	"strings"

	"gorm.io/gorm"

	"userstyles.world/modules/config"
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

// TotalSearchStyles returns total amount of userstyles for search page.
func TotalSearchStyles(query, sort string) (int, error) {
	q := "SELECT COUNT(*) FROM fts_styles WHERE fts_styles MATCH ?"
	if strings.HasPrefix(sort, "rating") {
		q += " AND (SELECT COUNT(*) FROM reviews WHERE style_id = fts_styles.id AND rating > 0)"
	}

	var total int
	err := database.Conn.
		Raw(q, query).
		Scan(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

// FindSearchStyles returns userstyles for search page.
func FindSearchStyles(query, sort string, page int) ([]*StyleCard, error) {
	var b strings.Builder
	b.WriteString(`SELECT styles.id, styles.name, styles.updated_at, styles.preview,
(SELECT username FROM users WHERE users.id = styles.user_id AND deleted_at IS NULL) AS username,
(SELECT total_views FROM histories WHERE histories.style_id = fts.id ORDER BY id DESC LIMIT 1) AS views,
(SELECT total_installs FROM histories WHERE histories.style_id = fts.id ORDER BY id DESC LIMIT 1) AS installs,
(SELECT ROUND(AVG(rating), 1) FROM reviews r WHERE r.style_id = fts.id AND r.rating > 0 AND r.deleted_at IS NULL) AS rating,
(SELECT COUNT(rating) FROM reviews r WHERE r.style_id = fts.id AND r.rating > 0 AND r.deleted_at IS NULL) AS ReviewCount
FROM fts_styles AS fts
JOIN styles ON styles.id = fts.id
WHERE fts_styles
MATCH ?`)

	if sort == "styles.id ASC" {
		b.WriteString(" ORDER BY bm25(fts_styles, 1, 2, 1.5, 1)")
	} else if sort != "" {
		if strings.HasPrefix(sort, "rating") {
			b.WriteString(" AND rating > 0")
		}
		b.WriteString(" ORDER BY ")
		b.WriteString(sort)
	}
	b.WriteString(" LIMIT ")
	b.WriteString(strconv.Itoa(config.App.PageMaxItems))
	if page > 1 {
		b.WriteString(" OFFSET ")
		b.WriteString(strconv.Itoa((page - 1) * config.App.PageMaxItems))
	}

	var s []*StyleCard
	err := database.Conn.Raw(b.String(), query).Scan(&s).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}

// DeleteSearchData removes a userstyle from FTS table.
func DeleteSearchData(db *gorm.DB, id int) error {
	return db.Exec("DELETE FROM fts_styles WHERE id = ?", id).Error
}
