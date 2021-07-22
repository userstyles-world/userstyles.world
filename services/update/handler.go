package update

import (
	"time"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

var size = 20

func ImportedStyles() {
	styles, err := models.GetImportedStyles()
	log.Info.Printf("Updating %d mirrored styles.\n", len(styles))

	if err != nil {
		log.Warn.Println("Failed to find imported styles:", err.Error())
		return // stop if error occurs
	}

	for i := 0; i < len(styles); i += size {
		j := i + size
		if j > len(styles) {
			j = len(styles)
		}

		for _, style := range styles[i:j] {
			time.Sleep(500 * time.Millisecond)
			go Batch(style)
		}

		time.Sleep(time.Second * 15)
	}

	log.Info.Println("Mirrored styles have been updated.")
}
