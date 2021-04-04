package core

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

var monitorhandler = monitor.New()

func Monitor(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Only first user (admin) is allowed.
	if u.Role == models.Admin {
		return monitorhandler(c)
	}

	return c.Render("err", fiber.Map{
		"Title": "Page not found",
		"User":  u,
	})
}
