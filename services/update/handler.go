package update

import (
	"log"
	"time"

	"userstyles.world/models"
)

var size = 20

func ImportedStyles() {
	styles, err := models.GetImportedStyles()
	log.Printf("Updating %d mirrored styles.\n", len(styles))

	if err != nil {
		log.Printf("Failed to find imported styles, err: %v", err)
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

	log.Println("Mirrored styles have been updated.")
}
