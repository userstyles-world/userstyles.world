package review

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
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
		Message:        strings.TrimSpace(c.FormValue("message")),
		Censor:         c.FormValue("censor") == "on",
	}

	if err = database.Conn.Create(&l).Error; err != nil {
		c.Locals("Title", "Failed to add mod log entry")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	n := models.Notification{
		Kind:     models.KindRemovedReview,
		TargetID: int(r.UserID),
		UserID:   int(u.ID),
		StyleID:  sid,
	}
	if err = n.Create(); err != nil {
		c.Locals("Title", "Failed to add notification")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	args := fiber.Map{
		"User": r.User,
		"Log":  l,
		"Link": config.BaseURL + "/modlog#id-" + strconv.Itoa(int(l.ID)),
	}

	title := "Your review has been removed"
	if err := email.Send("review/remove", r.User.Email, title, args); err != nil {
		log.Warn.Printf("Failed to email author for review %d: %s\n", rid, err)
	}

	return c.Redirect(fmt.Sprintf("/style/%d", sid), fiber.StatusSeeOther)
}
