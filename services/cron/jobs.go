package cron

import (
	"time"

	"github.com/go-co-op/gocron"

	"userstyles.world/modules/log"
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

	_, err = s.Cron("*/30 * * * *").Do(func() { update.ImportedStyles() })
	if err != nil {
		log.Warn.Println("Failed to update imported styles:", err.Error())
	}
}
