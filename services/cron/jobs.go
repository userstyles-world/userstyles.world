package cron

import (
	"log"
	"time"

	"github.com/go-co-op/gocron"

	"userstyles.world/services/snapshot"
)

func Initialize() {
	s := gocron.NewScheduler(time.Local)
	s.WaitForScheduleAll()
	s.StartAsync()

	_, err := s.Cron("59 23 * * *").Do(func() { snapshot.StyleStatistics() })
	if err != nil {
		log.Printf("History snapshop failed, err: %v\n", err)
	}
}
