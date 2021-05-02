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

func getViews(id int64) (i int64) {
	day := time.Now().AddDate(0, 0, -1)
	database.DB.
		Model(models.Stats{}).
		Where("style_id = ? and updated_at > ?", id, day).
		Count(&i)

	return i
}

func getUpdates(id int64) (i int64) {
	day := time.Now().AddDate(0, 0, -1)
	q := "style_id = ? and install = ? and updated_at > ?"
	database.DB.
		Debug().
		Model(models.Stats{}).
		Where(q, id, true, day).
		Count(&i)

	return i
}

func getInstalls(id int64) (i int64) {
	day := time.Now().AddDate(0, 0, -1)
	q := "style_id = ? and install = ? and created_at > ?"
	database.DB.
		Debug().
		Model(models.Stats{}).
		Where(q, id, true, day).
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
			DailyUpdates:  getUpdates(int64(v.ID)),
		}

		*items = append(*items, item)
	}

	log.Println("Stats history.")
	database.DB.Debug().Create(items)
}
