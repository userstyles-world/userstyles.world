package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/sessions"
	"userstyles.world/utils"
)

func Monitor(c *fiber.Ctx) error {
	u := sessions.User(c)

	// Only first user (admin) is allowed.
	if u.ID == 1 {
		return c.Redirect(utils.MonitorURL, fiber.StatusSeeOther)
	}

	return c.Render("err", fiber.Map{
		"Title": "Page not found",
	})
}
