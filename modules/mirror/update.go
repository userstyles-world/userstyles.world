// Package mirror provides functionality for mirroring userstyles.
package mirror

import (
	"sync"
	"time"

	"userstyles.world/models"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

// MirrorStyles runs an update check for mirrored userstyles in the background.
func MirrorStyles() {
	count, err := storage.CountStylesForMirror()
	if err != nil {
		log.Warn.Println("Failed to count mirrored styles:", err)
		return
	}

	log.Info.Printf("Checking %d mirrored styles.\n", count)

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
		log.Warn.Println("Failed to find mirrored styles:", err)
		return
	}

	log.Info.Printf("Done checking mirrored styles.\n")
}
