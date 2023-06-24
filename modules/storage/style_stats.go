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
	// TODO: Add weekly stats collection to speed up this query.
	q := `SELECT
(SELECT total_views FROM histories WHERE style_id = s.id ORDER BY id DESC) total_views,
(SELECT total_installs FROM histories where style_id = s.id ORDER BY id DESC) total_installs,
(SELECT count(*) FROM stats WHERE style_id = s.id AND created_at > DATETIME('now', '-7 days') AND install > 0) weekly_installs,
(SELECT count(*) FROM stats WHERE style_id = s.id AND created_at < DATETIME('now', '-7 days') AND install > DATETIME('now', '-7 days')) weekly_updates
FROM styles s
WHERE id = ?`

	var s *styleStats
	if err := database.Conn.Raw(q, id).Scan(&s).Error; err != nil {
		return nil, err
	}

	return s, nil
}
