package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func viewPage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	sid, err := c.ParamsInt("s")
	if err != nil || sid < 1 {
		c.Locals("Title", "Invalid style ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	rid, err := c.ParamsInt("r")
	if err != nil || rid < 1 {
		c.Locals("Title", "Invalid review ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	r, err := models.GetReview(rid)
	if err != nil {
		c.Locals("Title", "Review not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Review", r)

	s, err := models.GetStyle(c.Params("s"))
	c.Locals("Title", "Review for "+s.Name)

	return c.Render("review/view", fiber.Map{})
}
