package snapshot

import (
	"time"

	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
)

const q = `INSERT INTO histories(style_id, created_at, updated_at, daily_views, daily_installs, daily_updates, total_views, total_installs, total_updates)
SELECT
	s.id, DATETIME('now'), DATETIME('now'),
	(SELECT COUNT(*) FROM stats WHERE style_id = s.id AND view > DATE('now', '-1 day') AND created_at > DATE('now', '-1 day')) AS daily_views,
	(SELECT COUNT(*) FROM stats WHERE style_id = s.id AND install > DATE('now', '-1 day') AND created_at > DATE('now', '-1 day')) AS daily_installs,
	(SELECT COUNT(*) FROM stats WHERE style_id = s.id AND install > DATE('now', '-1 day') AND created_at != install) AS daily_updates,
	(SELECT COUNT(*) FROM stats WHERE style_id = s.id AND view > 0) AS total_views,
	(SELECT COUNT(*) FROM stats WHERE style_id = s.id AND install > 0) AS total_installs,
	(SELECT COUNT(*) FROM stats WHERE style_id = s.id AND install != created_at) AS total_updates
FROM styles s
WHERE deleted_at IS NULL
`

func StyleStatistics() {
	log.Info.Println("Collecting stats history.")

	for i := 0; i < 10; i++ {
		// NOTE: Might need some tweaks; it looks a bit too easy.
		err := database.Conn.Exec(q).Error
		if err == nil {
			// Exit if query was successful, otherwise try again.
			break
		}

		log.Database.Printf("Failed to take a snapshot on try %d: %s\n", i, err)
		time.Sleep(500 * time.Millisecond)
	}

	log.Info.Println("Stats history is collected.")
}
