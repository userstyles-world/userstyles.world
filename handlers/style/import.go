package style

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/search"
	"userstyles.world/utils"
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
	r := c.FormValue("import")
	s := new(models.Style)

	// Check if someone tries submitting local userstyle.
	if strings.Contains(r, "file:///") {
		return c.Render("err", fiber.Map{
			"Title": "Can't import local userstyles.",
			"User":  u,
		})
	}

	// Check if userstyle is imported from USo-archive.
	if strings.HasPrefix(r, utils.ArchiveURL) {
		style, err := utils.ImportFromArchive(r, *u)
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": err,
				"User":  u,
			})
		}

		// Move style content to outer scope.
		s = style
	} else {
		// Get userstyle.
		uc := new(usercss.UserCSS)
		if err := uc.ParseURL(r); err != nil {
			log.Warn.Println("Failed to parse userstyle from URL:", err.Error())
			return c.Render("err", fiber.Map{
				"Title": "Failed to fetch external userstyle",
				"User":  u,
			})
		}
		if errs := uc.Validate(); errs != nil {
			return c.Render("err", fiber.Map{
				"Title": "Failed to validate external userstyle",
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
		s.Original = r
	}

	// Enable code/meta mirroring.
	s.MirrorCode = c.FormValue("mirrorCode") == "on"
	s.MirrorMeta = c.FormValue("mirrorMeta") == "on"

	// Prevent importing multiples of the same style.
	err := models.CheckDuplicateStyle(s)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": err,
			"User":  u,
		})
	}

	s, err = models.CreateStyle(s)
	if err != nil {
		log.Warn.Println("Failed to import style from URL:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	// Check preview image.
	file, _ := c.FormFile("preview")
	preview := c.FormValue("previewURL")
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

	if err = search.IndexStyle(s.ID); err != nil {
		log.Warn.Printf("Failed to index style %d: %s\n", s.ID, err)
	}

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}
