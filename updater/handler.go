package updater

import (
	"log"
	"time"

	"userstyles.world/database"
	"userstyles.world/models"
)

var (
	startupDelay time.Duration
	updaterTick  time.Duration
	batchSize    int
)

func Initialize() {
	startupDelay = time.Second * 15
	updaterTick = time.Minute * 30
	batchSize = 5
	go Tick()
}

func Tick() {
	time.Sleep(startupDelay)
	for {
		importedStyles, err := models.GetAllImportedStyles(database.DB)
		if err == nil {
			stylesLen := len(*importedStyles)
			n, i := 0, 0
		out:
			for {
				for n < batchSize {
					if i >= stylesLen {
						break out
					}
					go UpdateBatch((*importedStyles)[i])
					i++
					n++
				}
				time.Sleep(time.Millisecond * 50)
				n = 0
			}
		} else {
			log.Printf("Updater: Wannted to update imported styles, but caught error: %s", err)
		}
		time.Sleep(updaterTick)
	}
}
