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

	clx, err := models.GetChangelogs(database.Conn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
			"Title": "Failed to get data",
		})
	}

	return c.Render("core/changelog", fiber.Map{
		"Title":      "Changelog",
		"Changelogs": clx,
	})
}

func createChangelogPage(c *fiber.Ctx) error {
	return c.Render("core/changelog-create", fiber.Map{
		"Title": "Create a changelog",
	})
}

func createChangelogForm(c *fiber.Ctx) error {
	cl := models.Changelog{UserID: int(c.Locals("User").(*models.APIUser).ID)}
	if err := c.BodyParser(&cl); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Failed to parse data",
			"Error": err,
		})
	}

	if err := models.CreateChangelog(database.Conn, cl); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
			"Title": "Failed to insert data",
		})
	}

	return c.Redirect("/changelog")
}

func editChangelogPage(c *fiber.Ctx) error {
	return c.Render("core/changelog-create", fiber.Map{
		"Title":     "Edit changelog",
		"Changelog": c.Locals("Changelog").(models.Changelog),
	})
}

func editChangelogForm(c *fiber.Ctx) error {
	cl := c.Locals("Changelog").(models.Changelog)

	if err := c.BodyParser(&cl); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Failed to parse data",
		})
	}

	if err := models.UpdateChangelog(database.Conn, cl); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
			"Title": "Failed to update data",
		})
	}

	return c.Redirect("/changelog")
}

func deleteChangelogPage(c *fiber.Ctx) error {
	return c.Render("core/changelog-delete", fiber.Map{
		"Title":     "Delete changelog",
		"Changelog": c.Locals("Changelog").(models.Changelog),
	})
}

func deleteChangelogForm(c *fiber.Ctx) error {
	cl := c.Locals("Changelog").(models.Changelog)
	if err := models.DeleteChangelog(database.Conn, cl); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
			"Title": "Failed to delete data",
		})
	}

	return c.Redirect("/changelog")
}
