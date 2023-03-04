package cron

import (
	"time"

	"github.com/go-co-op/gocron"

	// "userstyles.world/models"
	// "userstyles.world/modules/cache"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/database/snapshot"
	"userstyles.world/modules/log"
	"userstyles.world/modules/mirror"
	"userstyles.world/modules/sitemap"
	"userstyles.world/modules/storage"
)

func Initialize() {
	s := gocron.NewScheduler(time.Local)
	s.WaitForScheduleAll()
	s.StartAsync()

	_, err := s.Cron("59 23 * * *").Do(func() { snapshot.StyleStatistics() })
	if err != nil {
		log.Warn.Println("Failed to snapshot style statistics:", err.Error())
	}

	/*
		_, err = s.Every("1h").Do(func() {
			cache.Store.Add("siteStatistics", models.GetHomepageStatistics(), 5*time.Minute)
			styles, err := models.GetAllFeaturedStyles()
			if err != nil {
				return
			}
			cache.Store.Add("featuredStyles", styles, 5*time.Minute)
		})
		if err != nil {
			log.Warn.Println("Failed to cache home page queries:", err.Error())
		}
	*/

	_, err = s.Cron("4 */4 * * *").Do(func() { mirror.MirrorStyles() })
	if err != nil {
		log.Warn.Println("Failed to update imported styles:", err.Error())
	}

	_, err = s.Cron("30 */2 * * *").Do(func() {
		err := sitemap.UpdateSitemapCache()
		if err != nil {
			log.Warn.Println("Failed to update sitemap:", err.Error())
		}
	})
	if err != nil {
		log.Warn.Println("Failed to update sitemap:", err.Error())
	}

	_, err = s.Every("15m").Do(func() {
		index, err := storage.GetStyleCompactIndex()
		if err != nil {
			log.Warn.Printf("Failed to get compact index: %s\n", err)
			return
		}
		cache.Store.Set("index", index, 0)
	})
	if err != nil {
		log.Warn.Println("Failed to set compact index job:", err)
	}
}
