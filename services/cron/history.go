package cron

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron"

	"userstyles.world/database"
	"userstyles.world/models"
)

func Initialize() {
	s := gocron.NewScheduler(time.Local)
	job, err := s.Cron("20 00 * * *").Do(func() { snapshotStats() })
	s.StartAsync()
	fmt.Printf("job: %v, err: %v\n", job, err)
}

// TODO: Refactor and be use proper data after 2021-05-01.
// NOTE: As the first entry, we'll use _all_ the available data.
func getViews(id int64) (i int64) {
	database.DB.
		Model(models.Stats{}).
		Where("style_id = ?", id).
		Count(&i)

	return i
}

func getInstalls(id int64) (i int64) {
	database.DB.
		Model(models.Stats{}).
		Where("style_id = ? and install = ?", id, true).
		Count(&i)

	return i
}

func snapshotStats() {
	// TODO: Write a query to only get style IDs.
	styles, err := models.GetAllStyles(database.DB)
	if err != nil {
		log.Fatalln(err)
	}
	items := new([]models.History)

	for _, v := range *styles {
		item := models.History{
			StyleID:       v.ID,
			DailyViews:    getViews(int64(v.ID)),
			DailyInstalls: getInstalls(int64(v.ID)),
			DailyUpdates:  getInstalls(int64(v.ID)),
		}

		*items = append(*items, item)
	}

	log.Println("Stats history.")
	database.DB.Debug().Create(items)
}
