package style

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func ImportGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	return c.Render("import", fiber.Map{
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
		uc, err := usercss.ParseFromURL(r)
		if err != nil {
			log.Println("ParsingFromURL err:", err)
			return c.Render("err", fiber.Map{
				"Title": "Failed to fetch external userstyle",
				"User":  u,
			})
		}
		if errs := usercss.BasicMetadataValidation(uc); errs != nil {
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
		s.Preview = c.FormValue("preview")
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
		log.Println("Style import failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}
