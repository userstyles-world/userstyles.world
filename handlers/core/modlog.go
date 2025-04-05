package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
)

// GetModLog renders the modlog view.
// It will pass trough the relevant information from the database.
func GetModLog(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	l, err := models.GetModLogs(database.Conn)
	if err != nil {
		return c.Render("err", fiber.Map{"Title": "Failed to get data"})
	}

	return c.Render("core/modlog", fiber.Map{
		"Logs":      l,
		"Title":     "Moderation log",
		"Canonical": "modlog",
	})
}
