package review

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func editPage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Edit review")

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

	return c.Render("review/create", fiber.Map{})
}

func editForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Edit review")

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

	rating, comment := c.FormValue("rating"), c.FormValue("comment")
	r := models.NewReviewUpdate(u.ID, uint(sid), uint(rid), rating, comment)
	c.Locals("Review", r)

	if err = r.Validate(); err != nil {
		c.Locals("Error", strings.ToTitle(err.Error()[:1])+err.Error()[1:]+".")
		return c.Render("review/create", fiber.Map{})
	}

	if err = r.UpdateFromUser(u.ID); err != nil {
		c.Locals("Title", "Failed to update your review")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	return c.Redirect(r.Permalink(), fiber.StatusSeeOther)
}
