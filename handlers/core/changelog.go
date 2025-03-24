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
	c.Locals("Title", "Create a changelog")

	return c.Render("core/changelog-create", fiber.Map{})
}

func createChangelogForm(c *fiber.Ctx) error {
	cl := models.Changelog{UserID: int(c.Locals("User").(*models.APIUser).ID)}
	if err := c.BodyParser(&cl); err != nil {
		c.Locals("Error", err)
		c.Locals("Title", "Failed to parse data")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	if err := models.CreateChangelog(database.Conn, cl); err != nil {
		c.Locals("Title", "Failed to insert data")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	return c.Redirect("/changelog")
}

func editChangelogPage(c *fiber.Ctx) error {
	c.Locals("Title", "Edit changelog")

	cl, err := models.GetChangelog(database.Conn, c.Locals("id").(int))
	if err != nil || cl.ID == 0 {
		c.Locals("Title", "Changelog not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Changelog", cl)

	return c.Render("core/changelog-create", fiber.Map{})
}

func editChangelogForm(c *fiber.Ctx) error {
	cl := models.Changelog{
		UserID: int(c.Locals("User").(*models.APIUser).ID),
		ID:     c.Locals("id").(int),
	}

	if err := c.BodyParser(&cl); err != nil {
		c.Locals("Title", "Failed to parse data")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	if err := models.UpdateChangelog(database.Conn, cl); err != nil {
		c.Locals("Title", "Failed to update data")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	return c.Redirect("/changelog")
}

func deleteChangelogPage(c *fiber.Ctx) error {
	c.Locals("Title", "Delete changelog")

	cl, err := models.GetChangelog(database.Conn, c.Locals("id").(int))
	if err != nil || cl.ID == 0 {
		c.Locals("Title", "Changelog not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Changelog", cl)

	return c.Render("core/changelog-delete", fiber.Map{})
}

func deleteChangelogForm(c *fiber.Ctx) error {
	cl, err := models.GetChangelog(database.Conn, c.Locals("id").(int))
	if err != nil {
		c.Locals("Title", "Failed to get data")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	if err := models.DeleteChangelog(database.Conn, cl); err != nil {
		c.Locals("Title", "Failed to delete data")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	return c.Redirect("/changelog")
}
