package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
)

func changelogPage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Changelog")

	clx, err := models.GetChangelogs(database.Conn)
	if err != nil {
		c.Locals("Title", "Failed to get data")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}
	c.Locals("Changelogs", clx)

	return c.Render("core/changelog", fiber.Map{})
}
