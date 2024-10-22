package mirror

import (
	"strconv"
	"sync"

	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/archive"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/log"
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
	style := &models.Style{ID: batch.ID}

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

	// Select which database fields should be updated.
	var f []string

	// Set new source code.
	if batch.MirrorCode {
		old := new(usercss.UserCSS)
		if err := old.Parse(batch.Code); err != nil {
			log.Warn.Printf("Failed to parse style %d.\n", batch.ID)
			return
		}
		if uc.Version != old.Version {
			style.Code = util.RemoveUpdateURL(uc.SourceCode)
			f = append(f, "code")
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
				f = append(f, "name", "description", "notes")
				style.Name = s.Name
				style.Description = s.Description
				style.Notes = s.Notes
				updateReady = true
			}
		} else if updateMeta(&batch, uc) {
			f = append(f, "name", "description", "homepage")
			style.Name = uc.Name
			style.Description = uc.Description
			style.Homepage = uc.HomepageURL
			updateReady = true
		}
	}

	if updateReady {
		// Update database record.
		err := models.MirrorStyle(style, f)
		if err != nil {
			log.Database.Printf("Failed to mirror %d: %s\n", style.ID, err)
			return
		}

		if style.Code != "" {
			i := int(style.ID)
			err = models.SaveStyleCode(strconv.Itoa(i), style.Code)
			if err != nil {
				log.Warn.Printf("kind=code id=%v err=%q\n", style.ID, err)
				return
			}

			cache.Code.Update(i, []byte(style.Code))
		}

		log.Info.Printf("Successfully mirrored style %d\n", style.ID)
	}
}
