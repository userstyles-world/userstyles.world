package update

import (
	"time"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

var size = 20

func ImportedStyles() {
	styles, err := models.GetImportedStyles()
	if err != nil {
		log.Warn.Println("Failed to find imported styles:", err.Error())
		return // stop if error occurs
	}

	length := len(styles)
	log.Info.Printf("Updating %d mirrored styles.\n", size)

	for i := 0; i < length; i += size {
		j := i + size
		if j > length {
			j = length
		}

		for _, style := range styles[i:j] {
			time.Sleep(500 * time.Millisecond)
			go Batch(style)
		}

		time.Sleep(time.Second * 15)
	}

	log.Info.Println("Mirrored styles have been updated.")
}
