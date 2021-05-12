package update

import (
	"log"
	"time"

	"userstyles.world/database"
	"userstyles.world/models"
)

var (
	size = 5
)

func ImportedStyles() {
	styles, err := models.GetImportedStyles(database.DB)
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
			time.Sleep(time.Second)
			go UpdateBatch(&style)
		}

		time.Sleep(time.Second * 15)
	}
}
