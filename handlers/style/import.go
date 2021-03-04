package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/sessions"
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
