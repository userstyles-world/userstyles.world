package review

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
)

func removePage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	if !u.IsModOrAdmin() {
		c.Locals("Title", "You are not authorized to perform this action")
		return c.Status(fiber.StatusUnauthorized).Render("err", fiber.Map{})
	}
	c.Locals("Title", "Remove review")

	rid, err := c.ParamsInt("r")
	if err != nil || rid < 1 {
		c.Locals("Title", "Invalid review ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	r, err := models.GetReview(rid)
	if err != nil {
		c.Locals("Title", "Failed to find review")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Review", r)

	return c.Render("review/remove", fiber.Map{})
}

func removeForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	if !u.IsModOrAdmin() {
		c.Locals("Title", "You are not authorized to perform this action")
		return c.Status(fiber.StatusUnauthorized).Render("err", fiber.Map{})
	}
	c.Locals("Title", "Remove review")

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
		c.Locals("Title", "Failed to find review")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	if err = models.DeleteReviewFromUser(int(r.ID), int(u.ID)); err != nil {
		c.Locals("Title", "Failed to delete review")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	l := models.Log{
		UserID:         u.ID,
		Username:       u.Username,
		Kind:           models.LogRemoveReview,
		TargetUserName: r.User.Username,
		TargetData:     strconv.Itoa(int(r.ID)),
		Reason:         strings.TrimSpace(c.FormValue("reason")),
		Censor:         c.FormValue("censor") == "on",
	}

	if err = database.Conn.Create(&l).Error; err != nil {
		c.Locals("Title", "Failed to add mod log entry")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	return c.Redirect(fmt.Sprintf("/style/%d", sid), fiber.StatusSeeOther)
}
