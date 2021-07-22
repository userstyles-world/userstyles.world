package snapshot

import (
	"time"

	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
)

func getViews(id int64) (i int64) {
	day := time.Now().AddDate(0, 0, -1)
	database.Conn.
		Model(models.Stats{}).
		Where("style_id = ? and created_at > ? and view > ?", id, day, day).
		Count(&i)

	return i
}

func getUpdates(id int64) (i int64) {
	day := time.Now().AddDate(0, 0, -1)
	q := "style_id = ? and install > ? and updated_at > ?"
	database.Conn.
		Model(models.Stats{}).
		Where(q, id, day, day).
		Count(&i)

	return i
}

func getInstalls(id int64) (i int64) {
	day := time.Now().AddDate(0, 0, -1)
	q := "style_id = ? and install > ? and created_at > ?"
	database.Conn.
		Model(models.Stats{}).
		Where(q, id, day, day).
		Count(&i)

	return i
}

func getPreviousHistory(id uint) (q *models.History) {
	database.Conn.
		Model(models.History{}).
		Where("style_id = ?", id).
		Order("id DESC").
		Find(&q)

	return q
}

func StyleStatistics() {
again:
	styles, err := models.GetAllStyleIDs()
	if err != nil {
		log.Warn.Println("Failed to get IDs for all styles:", err.Error())
		goto again
	}

	// Store style stats.
	stats := new([]models.History)

	// Iterate over styles and collect their stats.
	for _, v := range styles {
		prev := getPreviousHistory(v.ID)
		views := getViews(int64(v.ID))
		totalViews := prev.TotalViews + views
		installs := getInstalls(int64(v.ID))
		totalInstalls := prev.TotalInstalls + installs
		updates := getUpdates(int64(v.ID))
		totalUpdates := prev.TotalUpdates + updates

		item := models.History{
			StyleID:       v.ID,
			DailyViews:    views,
			DailyInstalls: installs,
			DailyUpdates:  updates,
			TotalViews:    totalViews,
			TotalInstalls: totalInstalls,
			TotalUpdates:  totalUpdates,
		}

		*stats = append(*stats, item)
	}

	log.Info.Println("Collecting stats history.")
	database.Conn.Debug().Create(stats)
	log.Info.Println("Stats history is collected.")
}
