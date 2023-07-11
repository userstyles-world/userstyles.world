package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func viewPage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	i, err := c.ParamsInt("r")
	if err != nil || i < 1 {
		c.Locals("Title", "Invalid review ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}
	c.Locals("Title", "Review "+c.Params("r"))

	r, err := models.GetReview(i)
	if err != nil {
		c.Locals("Title", "Review not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Review", r)

	return c.Render("review/view", fiber.Map{})
}
