package api

import (
	"fmt"
	"regexp"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetStyleSource(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(database.DB, id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	// Check if source code is a link.
	r, err := regexp.Compile(`^https?://.*\.user\.(css|styl|less)$`)
	if err != nil {
		return c.JSON(fiber.Map{"data": "internal server error"})
	}

	// Redirect to external userstyle.
	if r.MatchString(style.Code) {
		// TODO: Add validation for external userstyle.
		return c.Redirect(style.Code, fiber.StatusTemporaryRedirect)
	}

	c.Set("Content-Type", "text/css")
	return c.SendString(fmt.Sprintf("%s", style.Code))
}
