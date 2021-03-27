package updater

import (
	"log"
	"time"

	"userstyles.world/database"
	"userstyles.world/models"
)

var UpdaterTick time.Duration

func Initialize() {
	UpdaterTick = time.Minute * 30
	go Tick()
}

func Tick() {
	for {
		importedStyles, err := models.GetAllImportedStyles(database.DB)
		if err == nil {
			for i, len := 0, len(*importedStyles); i < len; i++ {
				go UpdateBatch((*importedStyles)[i])
			}
		} else {
			log.Printf("Updater: Wannted to update imported styles, but caught error: %s", err)
		}
		time.Sleep(UpdaterTick)
	}
}
