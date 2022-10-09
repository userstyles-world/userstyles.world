package update

import (
	"sync"
	"time"

	"userstyles.world/models"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

func ImportedStyles() {
	count, err := storage.CountStylesForMirror()
	if err != nil {
		log.Warn.Println("Failed to find imported styles:", err)
		return
	}

	log.Info.Printf("Updating %d mirrored styles.\n", count)

	var wg sync.WaitGroup
	action := func(styles []models.Style) error {
		wg.Add(len(styles))

		for _, style := range styles {
			go check(&wg, style)
		}

		time.Sleep(time.Second)
		wg.Wait()

		return nil
	}

	if err := storage.FindStylesForMirror(action); err != nil {
		log.Warn.Println("Failed to find imported styles:", err)
		return
	}

	log.Info.Println("Mirrored styles have been updated.")
}
