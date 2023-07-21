package storage

import (
	"userstyles.world/modules/database"
)

// styleStats contains stats used on style view page.
type styleStats struct {
	TotalViews     int
	TotalInstalls  int
	WeeklyInstalls int
	WeeklyUpdates  int
}

// GetStyleStats returns stats for style view page.
func GetStyleStats(id string) (*styleStats, error) {
	q := `SELECT
(SELECT total_views FROM histories WHERE style_id = s.id ORDER BY id DESC LIMIT 1) total_views,
(SELECT total_installs FROM histories where style_id = s.id ORDER BY id DESC LIMIT 1) total_installs,
(SELECT weekly_installs FROM histories WHERE style_id = s.id ORDER BY id DESC LIMIT 1) weekly_installs,
(SELECT weekly_updates FROM histories WHERE style_id = s.id ORDER BY id DESC LIMIT 1) weekly_updates
FROM styles s
WHERE id = ?`

	var s *styleStats
	if err := database.Conn.Raw(q, id).Scan(&s).Error; err != nil {
		return nil, err
	}

	return s, nil
}

// GetWeeklyViews returns weekly installs for a userstyle.
func GetWeeklyViews(id string) int {
	q := "SELECT count(*) FROM stats WHERE style_id = ? AND created_at > DATETIME('now', '-7 days') AND view > 0"

	var i int
	database.Conn.Raw(q, id).Scan(&i)

	return i
}

// GetWeeklyInstalls returns weekly installs for a userstyle.
func GetWeeklyInstalls(id string) int {
	q := "SELECT count(*) FROM stats WHERE style_id = ? AND created_at > DATETIME('now', '-7 days') AND install > 0"

	var i int
	database.Conn.Raw(q, id).Scan(&i)

	return i
}

// GetTotalViews returns total views for a userstyle.
func GetTotalViews(id string) int {
	q := "SELECT total_views FROM histories WHERE style_id = ? ORDER BY id DESC"

	var i int
	database.Conn.Raw(q, id).Scan(&i)

	return i
}

// GetTotalInstalls returns total installs for a userstyle.
func GetTotalInstalls(id string) int {
	q := "SELECT total_installs FROM histories where style_id = ? ORDER BY id DESC"

	var i int
	database.Conn.Raw(q, id).Scan(&i)

	return i
}
