package review

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func deletePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Delete review")

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

	r, err := models.GetReviewFromUser(rid, int(u.ID))
	if err != nil {
		c.Locals("Title", "Failed to find review")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Review", r)

	return c.Render("review/delete", fiber.Map{})
}

func deleteForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Delete review")

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

	if err = models.DeleteReviewFromUser(rid, int(u.ID)); err != nil {
		c.Locals("Title", "Failed to find review")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	return c.Redirect(fmt.Sprintf("/style/%d", sid), fiber.StatusSeeOther)
}
