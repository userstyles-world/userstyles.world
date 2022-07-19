package update

import (
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

	action := func(styles []models.Style) error {
		for _, style := range styles {
			time.Sleep(200 * time.Millisecond)
			go Batch(style)
		}

		return nil
	}

	if err := storage.FindStylesForMirror(action); err != nil {
		log.Warn.Println("Failed to find imported styles:", err)
		return
	}

	log.Info.Println("Mirrored styles have been updated.")
}
