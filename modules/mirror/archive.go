package mirror

import (
	"strconv"
	"sync"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/archive"
	"userstyles.world/modules/log"
	"userstyles.world/modules/search"
	"userstyles.world/modules/util"
)

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

func mirrorWG(wg *sync.WaitGroup, s models.Style) {
	mirror(s)
	wg.Done()
}

func mirror(batch models.Style) {
	// Select which fields to update.
	fields := make(map[string]any)
	fields["id"] = batch.ID

	// Don't update database record if nothing changed.
	var updateReady bool

	// Get new source code.
	uc := new(usercss.UserCSS)
	if err := uc.ParseURL(getSourceCode(batch)); err != nil {
		log.Warn.Printf("Failed to parse style %d from URL: %s\n", batch.ID, err.Error())
		return
	}

	// Exit if source code doesn't pass validation.
	if err := uc.Validate(); err != nil {
		log.Warn.Printf("Failed to validate style %d.\n", batch.ID)
		return
	}

	// Set new source code.
	if batch.MirrorCode {
		old := new(usercss.UserCSS)
		if err := old.Parse(batch.Code); err != nil {
			log.Warn.Printf("Failed to parse style %d.\n", batch.ID)
			return
		}
		if uc.Version != old.Version {
			fields["code"] = util.RemoveUpdateURL(uc.SourceCode)
			updateReady = true
		}
	}

	// Set new style metadata.
	if batch.MirrorMeta {
		url := setStyleURL(batch.MirrorURL, batch.Original)

		if archive.IsFromArchive(url) {
			s, err := archive.ImportFromArchive(url, models.APIUser{})
			if err != nil {
				log.Warn.Printf("Failed to import %s from archive: %s\n", url, err.Error())
				return
			}

			if updateArchiveMeta(&batch, s) {
				fields["name"] = s.Name
				fields["notes"] = s.Notes
				fields["description"] = s.Description
				updateReady = true
			}
		} else if updateMeta(&batch, uc) {
			fields["name"] = uc.Name
			fields["description"] = uc.Description
			fields["homepage"] = uc.HomepageURL
			updateReady = true
		}
	}

	if updateReady {
		// Update database record.
		err := batch.MirrorStyle(fields)
		if err != nil {
			log.Warn.Printf("Failed to mirror style %d: %s\n", batch.ID, err.Error())
		} else {
			log.Info.Printf("Successfully mirrored style %d\n", batch.ID)
		}

		err = models.SaveStyleCode(strconv.Itoa(int(batch.ID)), fields["code"].(string))
		if err != nil {
			log.Warn.Printf("kind=code id=%v err=%q\n", batch.ID, err)
		}

		// Update search index.
		err = search.IndexStyle(batch.ID)
		if err != nil {
			log.Warn.Printf("Failed to re-index style %d: %s\n", batch.ID, err.Error())
		}
	}
}
