package review

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
)

func deletePage(c *fiber.Ctx) error {
	return c.Render("review/delete", fiber.Map{
		"Title": "Delete review",
	})
}

func deleteForm(c *fiber.Ctx) error {
	u := c.Locals("User").(*models.APIUser)
	r := c.Locals("Review").(*models.Review)
	if err := models.DeleteReviewFromUser(int(r.ID), int(u.ID)); err != nil {
		c.Locals("Title", "Failed to find review")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	a := models.NewSuccessAlert("Review successfully deleted.")
	cache.Store.Add("alert "+u.Username, a, time.Minute)

	return c.Redirect(r.Style.Permalink())
}
