package style

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/archive"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
	"userstyles.world/modules/validator"
)

func ImportGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	return c.Render("style/import", fiber.Map{
		"Title": "Add userstyle",
		"User":  u,
	})
}

func ImportPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Check if someone tries submitting local userstyle.
	origin := strings.TrimSpace(c.FormValue("import"))
	if strings.Contains(origin, "file:///") {
		return c.Render("err", fiber.Map{
			"Title": "Can't import local userstyles",
			"User":  u,
		})
	}

	s := new(models.Style)
	switch {
	// Import from USo-archive.
	case archive.IsFromArchive(origin):
		var err error
		origin, err = archive.RewriteURL(origin)
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": err,
				"User":  u,
			})
		}

		s, err = archive.ImportFromArchive(origin, *u)
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": err,
				"User":  u,
			})
		}

	// Import from anywhere.
	case origin != "":
		// Get userstyle.
		uc := new(usercss.UserCSS)
		if err := uc.ParseURL(origin); err != nil {
			log.Warn.Println("Failed to parse userstyle from URL:", err.Error())
			return c.Render("err", fiber.Map{
				"Title": "Failed to fetch external userstyle",
				"User":  u,
			})
		}

		// Set fields.
		s.UserID = u.ID
		s.Name = uc.Name
		s.Code = uc.SourceCode
		s.License = uc.License
		s.Description = uc.Description
		s.Homepage = uc.HomepageURL
		s.Category = strings.TrimSpace(c.FormValue("category", "unset"))
		s.Original = origin

	// Validation stage.
	default:
		s = &models.Style{
			Name:        strings.TrimSpace(c.FormValue("name")),
			Description: strings.TrimSpace(c.FormValue("description")),
			Notes:       strings.TrimSpace(c.FormValue("notes")),
			Homepage:    strings.TrimSpace(c.FormValue("homepage")),
			License:     strings.TrimSpace(c.FormValue("license", "No License")),
			Code:        strings.TrimSpace(util.RemoveUpdateURL(c.FormValue("code"))),
			Category:    strings.TrimSpace(c.FormValue("category")),
			Original:    c.FormValue("original"),
			UserID:      u.ID,
		}
	}

	// Enable code/meta mirroring.
	s.ImportPrivate = c.FormValue("importPrivate") == "on"
	s.MirrorPrivate = c.FormValue("mirrorPrivate") == "on"
	s.MirrorCode = c.FormValue("mirrorCode") == "on"
	s.MirrorMeta = c.FormValue("mirrorMeta") == "on"

	// Get previewURL
	preview := c.FormValue("previewURL", s.Preview)

	m, err := s.Validate(validator.V, true)
	if err != nil {
		return c.Render("style/import", fiber.Map{
			"Title":      "Import userstyle",
			"User":       u,
			"Style":      s,
			"PreviewURL": preview,
			"err":        m,
			"Error":      "Incorrect userstyle data was entered. Please review the fields bellow.",
		})
	}

	// Prevent importing multiples of the same style.
	err = models.CheckDuplicateStyle(s)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": err,
			"User":  u,
		})
	}

	s.Code = util.RemoveUpdateURL(s.Code)

	s, err = models.CreateStyle(s)
	if err != nil {
		log.Warn.Println("Failed to import style from URL:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	err = models.SaveStyleCode(strconv.Itoa(int(s.ID)), s.Code)
	if err != nil {
		log.Warn.Printf("kind=code id=%v err=%q\n", s.ID, err)
	}

	// Check preview image.
	file, _ := c.FormFile("preview")
	styleID := strconv.FormatUint(uint64(s.ID), 10)
	if file != nil || preview != "" {
		if err := images.Generate(file, styleID, "0", "", preview); err != nil {
			log.Warn.Println("Error:", err)
			s.Preview = ""
		} else {
			s.SetPreview()
			if err = s.UpdateColumn("preview", s.Preview); err != nil {
				log.Warn.Printf("Failed to update preview for %s: %s\n", styleID, err)
			}
		}
	}

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}
