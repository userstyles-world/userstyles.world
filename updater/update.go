package updater

import (
	"log"
	"strings"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

const archiveURL = "https://raw.githubusercontent.com/33kk/uso-archive/"

func isUSo(url string) bool {
	if strings.HasPrefix(url, archiveURL) {
		return true
	}

	return false
}

func updateMeta(a, b *models.Style) bool {
	if a.Name == b.Name &&
		a.Notes == b.Notes &&
		a.Preview == b.Preview &&
		a.Description == b.Description {
		return false
	}

	return true
}

func UpdateBatch(batch *models.Style) {
	if batch.Archived {
		return
	}

	// Update style metadata if style comes from USo-archive.
	if isUSo(batch.Original) && batch.MirrorMeta {
		new, err := utils.ImportFromArchive(batch.Original, models.APIUser{})
		if err != nil {
			log.Printf("%v\n", err)
		}

		// Run update if fields differ.
		if updateMeta(batch, new) {
			s := models.Style{
				Name:        new.Name,
				Notes:       new.Notes,
				Preview:     new.Preview,
				Description: new.Description,
			}

			err := database.DB.
				Debug().
				Model(models.Style{}).
				Where("id", batch.ID).
				Updates(s).
				Error

			if err != nil {
				log.Printf("Updater: Updating style %d failed, caught err %s", batch.ID, err)
			}
		}
	}

	// Update style source code.
	style, err := usercss.ParseFromURL(batch.Original)
	if err != nil {
		log.Printf("Updater: Cannot fetch style %d.\n", batch.ID)
		return
	}
	if errs := usercss.BasicMetadataValidation(style); errs != nil {
		log.Printf("Updater: Cannot validate style %d.\n", batch.ID)
		return
	}

	// Mirror source code if versions don't match.
	if style.Version != usercss.ParseFromString(batch.Code).Version {
		log.Printf("Updater: Style %d was changed.\n", batch.ID)
		s := models.Style{
			Code: style.SourceCode,
		}

		err := database.DB.
			Model(models.Style{}).
			Where("id", batch.ID).
			Updates(s).
			Error

		if err != nil {
			log.Printf("Updater: Updating style %d failed, caught err %s", batch.ID, err)
		}
	}
}
