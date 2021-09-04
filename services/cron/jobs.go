package cron

import (
	"time"

	"github.com/go-co-op/gocron"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/log"
	"userstyles.world/modules/sitemap"
	"userstyles.world/modules/update"
	"userstyles.world/services/snapshot"
)

func Initialize() {
	s := gocron.NewScheduler(time.Local)
	s.WaitForScheduleAll()
	s.StartAsync()

	_, err := s.Cron("59 23 * * *").Do(func() { snapshot.StyleStatistics() })
	if err != nil {
		log.Warn.Println("Failed to snapshot style statistics:", err.Error())
	}

	_, err = s.Cron("*/5 * * * *").Do(func() {
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

	_, err = s.Cron("*/30 * * * *").Do(func() { update.ImportedStyles() })
	if err != nil {
		log.Warn.Println("Failed to update imported styles:", err.Error())
	}

	_, err = s.Cron("*/30 * * * *").Do(func() {
		err := sitemap.UpdateSitemapCache()
		if err != nil {
			log.Warn.Println("Failed to update sitemap:", err.Error())
		}
	})
	if err != nil {
		log.Warn.Println("Failed to update sitemap:", err.Error())
	}
}
