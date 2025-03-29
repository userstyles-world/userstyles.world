package review

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
)

func removePage(c *fiber.Ctx) error {
	return c.Render("review/remove", fiber.Map{"Title": "Remove review"})
}

func removeForm(c *fiber.Ctx) error {
	u := c.Locals("User").(*models.APIUser)
	c.Locals("Title", "Remove review")

	r := c.Locals("Review").(*models.Review)
	if err := models.DeleteReviewFromUser(int(r.ID), int(r.UserID)); err != nil {
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

	if err := database.Conn.Create(&l).Error; err != nil {
		c.Locals("Title", "Failed to add mod log entry")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	n := models.Notification{
		Kind:     models.KindRemovedReview,
		TargetID: int(r.UserID),
		UserID:   int(u.ID),
		StyleID:  int(r.StyleID),
	}
	if err := models.CreateNotification(database.Conn, &n); err != nil {
		c.Locals("Title", "Failed to add notification")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	args := fiber.Map{
		"User": r.User,
		"Log":  l,
		"Link": config.App.BaseURL + "/modlog#id-" + strconv.Itoa(int(l.ID)),
	}

	title := "Your review has been removed"
	if err := email.Send("review/remove", r.User.Email, title, args); err != nil {
		log.Warn.Printf("Failed to email author for review %d: %s\n", r.ID, err)
	}

	a := models.NewSuccessAlert("Review successfully removed.")
	cache.Store.Add("alert "+u.Username, a, time.Minute)

	return c.Redirect(fmt.Sprintf("/modlog#id-%d", l.ID))
}
