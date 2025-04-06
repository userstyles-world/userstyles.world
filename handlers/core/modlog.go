package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
)

// GetModLog renders the modlog view.
// It will pass trough the relevant information from the database.
func GetModLog(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	var p models.Pagination
	if err := c.QueryParser(&p); err != nil {
		log.Info.Printf("Parsing pagination failed: %s", err)
		return c.Render("err", fiber.Map{"Title": "Invalid pagination data"})
	}

	if p.Kind < 0 || p.Kind >= int(models.LogCount) {
		return c.Render("err", fiber.Map{"Title": "Invalid sort method"})
	}

	total, err := models.GetModLogCount(database.Conn, p.Kind)
	if err != nil {
		log.Database.Printf("GetModLogCount failed: %s", err)
		return c.Render("err", fiber.Map{"Title": "Failed to get mod log data"})
	}

	url, ok := p.ModLogCheck(int(total))
	if !ok {
		return c.Redirect(url, 302)
	}
	c.Locals("P", p)

	l, err := models.GetModLogs(database.Conn, p.Now, config.App.PageMaxItems, p.Kind)
	if err != nil {
		log.Database.Printf("GetModLogs failed: %s: %v", err, p.Kind)
		return c.Render("err", fiber.Map{"Title": "Failed to get data"})
	}

	return c.Render("core/modlog", fiber.Map{
		"Logs":      l,
		"Title":     "Moderation log",
		"Canonical": "modlog",
	})
}
