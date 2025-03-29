package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
)

func viewPage(c *fiber.Ctx) error {
	r := c.Locals("Review").(*models.Review)
	return c.Render("review/view", fiber.Map{
		"Title": "Review for " + r.Style.Name,
	})
}
