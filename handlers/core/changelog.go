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

	return c.Render("core/changelog-create", fiber.Map{})
}

func createChangelogForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

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

func editChangelogPage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Edit changelog")

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		c.Locals("Title", "ID must be a positive number")
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{})
	}

	cl, err := models.GetChangelog(database.Conn, i)
	if err != nil || cl.ID == 0 {
		c.Locals("Title", "Changelog not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Changelog", cl)

	return c.Render("core/changelog-create", fiber.Map{})
}

func editChangelogForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	var cl models.Changelog
	if err := c.BodyParser(&cl); err != nil {
		c.Locals("Title", "Failed to parse data")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}
	cl.UserID = int(u.ID)

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		c.Locals("Title", "ID must be a positive number")
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{})
	}
	cl.ID = i

	if err := models.UpdateChangelog(database.Conn, cl); err != nil {
		c.Locals("Title", "Failed to update data")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	return c.Redirect("/changelog")
}

func deleteChangelogPage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		c.Locals("Title", "ID must be a positive number")
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{})
	}

	cl, err := models.GetChangelog(database.Conn, i)
	if err != nil || cl.ID == 0 {
		c.Locals("Title", "Changelog not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Changelog", cl)

	return c.Render("core/changelog-delete", fiber.Map{})
}

func deleteChangelogForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		c.Locals("Title", "ID must be a positive number")
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{})
	}

	cl, err := models.GetChangelog(database.Conn, i)
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
