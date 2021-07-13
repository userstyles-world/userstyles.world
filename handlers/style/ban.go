package style

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/search"
)

func BanGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Check if logged-in user has permissions.
	if !u.IsModOrAdmin() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("err", fiber.Map{
			"Title": "Can't do that",
			"User":  u,
		})
	}

	// Check if style exists.
	s, err := models.GetStyleByID(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	return c.Render("style/ban", fiber.Map{
		"Title": "Confirm ban",
		"User":  u,
		"Style": s,
	})
}

func BanPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	// Check if logged-in user has permissions.
	if !u.IsModOrAdmin() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("err", fiber.Map{
			"Title": "Can't do that",
			"User":  u,
		})
	}

	// Check if style exists.
	s, err := models.GetStyleByID(id)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Delete from database.
	q := new(models.Style)
	if err = database.Conn.Delete(q, "styles.id = ?", id).Error; err != nil {
		log.Printf("Failed to delete style, err: %#+v\n", err)
		c.Status(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	// Delete from search index.
	if err = search.DeleteStyle(s.ID); err != nil {
		log.Printf("Couldn't delete style %d failed, err: %s", s.ID, err.Error())
	}

	return c.Redirect("/account", fiber.StatusSeeOther)
}
