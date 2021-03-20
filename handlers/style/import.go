package style

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func StyleImportGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	return c.Render("import", fiber.Map{
		"Title": "Add userstyle",
		"User":  u,
	})
}

func StyleImportPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	r := c.FormValue("import")

	// Check if someone tries submitting local userstyle.
	if strings.Contains(r, "file:///") {
		return c.Render("err", fiber.Map{
			"Title": "Can't import local userstyles.",
			"User":  u,
		})
	}

	// Get userstyle.
	uc, err := usercss.ParseFromURL(r)
	if err != nil {
		log.Println("ParsingFromURL err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to fetch external userstyle",
			"User":  u,
		})
	}
	if valid, _ := usercss.BasicMetadataValidation(uc); !valid {
		return c.Render("err", fiber.Map{
			"Title": "Failed to validate external userstyle",
			"User":  u,
		})
	}

	s := models.Style{
		UserID:      u.ID,
		Name:        uc.Name,
		Code:        uc.SourceCode,
		License:     uc.License,
		Description: uc.Description,
		Homepage:    uc.HomepageURL,
		Preview:     c.FormValue("preview"),
		Category:    c.FormValue("category"),
		Original:    r,
	}

	s, err = models.CreateStyle(database.DB, s)
	if err != nil {
		log.Println("Style import failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect(fmt.Sprintf("/style/%d", int(s.ID)), fiber.StatusSeeOther)
}
