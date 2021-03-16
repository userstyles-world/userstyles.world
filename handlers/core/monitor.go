package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	"userstyles.world/handlers/sessions"
)

var monitorhandler = monitor.New()

func Monitor(c *fiber.Ctx) error {
	u := sessions.User(c)

	// Only first user (admin) is allowed.
	if u.ID == 1 {
		return monitorhandler(c)
	}

	return c.Render("err", fiber.Map{
		"Title": "Page not found",
		"User":  u,
	})
}
