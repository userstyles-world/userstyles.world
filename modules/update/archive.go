package update

import (
	"strings"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/log"
	"userstyles.world/search"
	"userstyles.world/utils"
)

func isArchive(url string) bool {
	return strings.HasPrefix(url, utils.ArchiveURL)
}

func setStyleURL(a, b string) string {
	if a != "" {
		return a
	}

	return b
}

func updateArchiveMeta(a, b *models.Style) bool {
	if a.Name == b.Name &&
		a.Notes == b.Notes &&
		a.Description == b.Description {
		return false
	}

	return true
}

func updateMeta(a *models.Style, b *usercss.UserCSS) bool {
	if a.Name == b.Name &&
		a.Description == b.Description &&
		a.Homepage == b.HomepageURL {
		return false
	}

	return true
}

func getSourceCode(style models.Style) string {
	if style.MirrorURL != "" {
		return style.MirrorURL
	}

	return style.Original
}

func Batch(batch models.Style) {
	// Create new model.
	style := new(models.Style)
	style.ID = batch.ID

	// Don't update database record if nothing changed.
	var updateReady bool

	// Get new source code.
	uc, err := usercss.ParseFromURL(getSourceCode(batch))
	if err != nil {
		log.Warn.Printf("Failed to parse style %d from URL: %s\n", batch.ID, err.Error())
		return
	}

	// Exit if source code doesn't pass validation.
	if errs := usercss.BasicMetadataValidation(uc); errs != nil {
		log.Warn.Printf("Failed to validate style %d.\n", batch.ID)
		return
	}

	// Set new source code.
	if batch.MirrorCode {
		if uc.Version != usercss.ParseFromString(batch.Code).Version {
			style.Code = uc.SourceCode
			updateReady = true
		}
	}

	// Set new style metadata.
	if batch.MirrorMeta {
		url := setStyleURL(batch.MirrorURL, batch.Original)

		if isArchive(url) {
			s, err := utils.ImportFromArchive(url, models.APIUser{})
			if err != nil {
				log.Warn.Printf("Failed to import %s from archive: %s\n", url, err.Error())
			}

			if updateArchiveMeta(&batch, s) {
				style.Name = s.Name
				style.Notes = s.Notes
				style.Preview = s.Preview
				style.Description = s.Description
				updateReady = true
			}
		} else {
			if updateMeta(&batch, uc) {
				style.Name = uc.Name
				style.Description = uc.Description
				style.Homepage = uc.HomepageURL
				updateReady = true
			}
		}
	}

	if updateReady {
		// Update database record.
		if err = models.UpdateStyle(style); err != nil {
			log.Warn.Printf("Failed to mirror style %d: %s\n", batch.ID, err.Error())
		} else {
			log.Info.Printf("Successfully mirrored style %d\n", batch.ID)
		}

		// Update search index.
		if err = search.IndexStyle(style.ID); err != nil {
			log.Warn.Printf("Failed to re-index style %d: %s\n", style.ID, err.Error())
		}
	}
}
