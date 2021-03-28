package updater

import (
	"fmt"
	"log"

	"github.com/vednoc/go-usercss-parser"
	"userstyles.world/database"
	"userstyles.world/models"
)

func UpdateBatch(batch models.APIStyle) {
	if batch.Archived {
		return
	}
	t := new(models.Style)
	style, err := usercss.ParseFromURL(batch.Original)
	if err != nil {
		fmt.Printf("Updater: Cannot fetch style %d.\n", batch.ID)
		return
	}
	if valid, _ := usercss.BasicMetadataValidation(style); !valid {
		fmt.Printf("Updater: Cannot validate style %d.\n", batch.ID)
		return
	}
	if style.SourceCode != batch.Code {
		fmt.Printf("Updater: Style %d was changed.\n", batch.ID)
		s := models.Style{
			UserID:      batch.UserID,
			Name:        batch.Name,
			Code:        style.SourceCode,
			Description: batch.Description,
			Homepage:    batch.Homepage,
			Preview:     batch.Preview,
			Category:    batch.Category,
			Original:    batch.Original,
		}
		s.Code = style.SourceCode

		err := database.DB.
			Model(t).
			Where("id", batch.ID).
			Updates(s).
			Error

		if err != nil {
			log.Printf("Updater: Updating style %d failed, caught err %s", batch.ID, err)

		}

	}
}
