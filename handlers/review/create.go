package review

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
)

func createPage(c *fiber.Ctx) error {
	return c.Render("review/create", fiber.Map{"Title": "Create review"})
}

func createForm(c *fiber.Ctx) error {
	c.Locals("Title", "Create review")

	u := c.Locals("User").(*models.APIUser)
	s := c.Locals("Style").(*models.Style)
	r := models.NewReview(u.ID, s.ID, c.FormValue("rating"), c.FormValue("comment"))
	c.Locals("Review", r)

	if err := r.Validate(); err != nil {
		return c.Render("review/create", fiber.Map{
			"Error": "Validation error: " + err.Error(),
		})
	}

	// Add review to database.
	if err := r.CreateForStyle(); err != nil {
		log.Warn.Printf("Failed to add review to style %d: %s\n", s.ID, err)
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
			"Title": "Failed to create your review",
		})
	}

	n := models.Notification{
		Kind:     models.KindReview,
		TargetID: int(s.UserID),
		UserID:   int(u.ID),
		StyleID:  int(s.ID),
		ReviewID: int(r.ID),
	}

	if err := models.CreateNotification(database.Conn, &n); err != nil {
		log.Warn.Printf("Failed to add notification to review %d: %s\n", r.ID, err)
	}

	a := models.NewSuccessAlert("Review has been created.")
	cache.Store.Add("alert "+u.Username, a, time.Minute)

	r.Style = *s      // assign *after* DB queries
	r.Style.Prepare() // generate slug for permalink

	return c.Redirect(r.Permalink())
}
