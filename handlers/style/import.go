package style

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

func StyleImportGet(c *fiber.Ctx) error {
	u := sessions.User(c)

	if sessions.State(c).Fresh() {
		c.Status(fiber.StatusFound)
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to add new userstyle.",
		})
	}

	return c.Render("import", fiber.Map{
		"Title": "Add userstyle",
		"User":  u,
	})
}

func StyleImportPost(c *fiber.Ctx) error {
	u := sessions.User(c)
	r := c.FormValue("import")

	if sessions.State(c).Fresh() {
		return c.Render("login", fiber.Map{
			"Title": "Login is required",
			"Error": "You must log in to import a userstyle.",
		})
	}

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

	style := &models.Style{
		UserID:      u.ID,
		Name:        uc.Name,
		Code:        uc.SourceCode,
		Description: uc.Description,
		Homepage:    uc.HomepageURL,
		Preview:     c.FormValue("preview"),
		Category:    c.FormValue("category"),
		Original:    r,
	}

	err = database.DB.
		Debug().
		Create(style).
		Error

	if err != nil {
		log.Println("Style import failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect("/account", fiber.StatusSeeOther)
}
