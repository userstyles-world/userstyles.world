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

func createChangelogPage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Create a changelog")

	if !u.IsAdmin() {
		c.Locals("Title", "Unauthorized")
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{})
	}

	return c.Render("core/changelog-create", fiber.Map{})
}

func createChangelogForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	if !u.IsAdmin() {
		c.Locals("Title", "Unauthorized")
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{})
	}

	var cl models.Changelog
	if err := c.BodyParser(&cl); err != nil {
		c.Locals("Error", err)
		c.Locals("Title", "Failed to parse data")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}
	cl.UserID = int(u.ID)

	if err := models.CreateChangelog(database.Conn, cl); err != nil {
		c.Locals("Title", "Failed to insert data")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	return c.Redirect("/changelog")
}
