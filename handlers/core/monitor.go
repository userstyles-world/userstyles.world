package core

import (
	"github.com/userstyles-world/fiber/v2"
	"github.com/userstyles-world/fiber/v2/middleware/proxy"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
)

var addr = []byte(config.ProxyMonitor)

func Monitor(c *fiber.Ctx) error {
	u, ok := jwt.User(c)

	// Only admins are allowed here.
	if !ok || u.Role != models.Admin {
		return c.
			Status(fiber.StatusUnauthorized).
			Render("err", fiber.Map{
				"Title": "Access denied",
				"User":  u,
			})
	}

	// Add trailing slash, otherwise proxy won't work.
	if c.Path() == "/monitor" {
		return c.Redirect("/monitor/", fiber.StatusMovedPermanently)
	}

	// Proxy requests.
	url := addr
	url = append(url, c.Request().URI().Path()[8:]...)
	url = append(url, 63)
	url = append(url, c.Context().URI().QueryArgs().QueryString()...)
	if err := proxy.Do(c, string(url)); err != nil {
		return err
	}

	return nil
}
