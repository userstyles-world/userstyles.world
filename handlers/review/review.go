package review

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/handlers/middleware"
	"userstyles.world/models"
	"userstyles.world/modules/database"
)

// Routes provides routes for Fiber's router.
func Routes(app *fiber.App) {
	r := app.Group("/styles/:s-:slug/reviews")
	r.Get("/create", jwt.Protected, ensureStyle, createPage)
	r.Post("/create", jwt.Protected, ensureStyle, createForm)
	r.Get("/:r", middleware.Alert, ensureReview, viewPage)

	r = app.Group("/styles/:s-:slug/reviews/:r", jwt.Protected, ensureReview)
	r.Get("/edit", hasAuthor, editPage)
	r.Post("/edit", hasAuthor, editForm)
	r.Get("/delete", hasAuthor, deletePage)
	r.Post("/delete", hasAuthor, deleteForm)
	r.Get("/remove", hasStaff, removePage)
	r.Post("/remove", hasStaff, removeForm)
}

func ensureStyle(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	i, err := c.ParamsInt("s")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Style ID must be a positive number",
		})
	}

	s, err := models.GetStyleByID(i)
	if err != nil {
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{
			"Title": "Style not found",
		})
	}
	c.Locals("Style", s)

	if s.UserID == u.ID {
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
			"Title": "You can't review your own style",
		})
	}

	dur, ok := models.AbleToReview(u.ID, s.ID)
	if !ok {
		return c.Status(fiber.StatusTooManyRequests).Render("err", fiber.Map{
			"Title":    "Cannot review style",
			"ErrTitle": "You can review this style again in " + dur,
		})
	}

	return c.Next()
}

func ensureReview(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	s, err := c.ParamsInt("s")
	if err != nil || s < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Style ID must be a positive number",
		})
	}

	r, err := c.ParamsInt("r")
	if err != nil || r < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Review ID must be a positive number",
		})
	}

	if ok := models.MatchReviewForStyle(database.Conn, r, s); !ok {
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{
			"Title": "Review not found",
		})
	}

	review, err := models.GetReviewForStyle(database.Conn, r, s)
	if err != nil {
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
			"Title": "Failed to get the review",
		})
	}
	c.Locals("Review", review)

	return c.Next()
}

func hasAuthor(c *fiber.Ctx) error {
	u := c.Locals("User").(*models.APIUser)
	r := c.Locals("Review").(*models.Review)
	if r.UserID != u.ID {
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
			"Title": "You're not the author of this review",
		})
	}

	return c.Next()
}

func hasStaff(c *fiber.Ctx) error {
	u := c.Locals("User").(*models.APIUser)
	if !u.IsModOrAdmin() {
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
			"Title": "Unauthorized access",
		})
	}

	return c.Next()
}
