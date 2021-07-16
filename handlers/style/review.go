package style

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func ReviewGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	return c.Render("style/review", fiber.Map{
		"Title": "Review style",
		"User":  u,
		"ID":    id,
	})
}

func ReviewPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	cmt := c.FormValue("comment")

	// Check if style exists.
	if _, err := models.GetStyleByID(c.Params("id")); err != nil {
		c.Status(fiber.StatusNotFound)
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	r, err := strconv.Atoi(c.FormValue("rating"))
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Invalid style ID",
			"User":  u,
		})
	}

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Invalid style ID",
			"User":  u,
		})
	}

	// Create review.
	review := models.Review{
		Rating:  r,
		Comment: cmt,
		StyleID: id,
		UserID:  int(u.ID),
	}

	// Add review to database.
	if err := review.CreateForStyle(id); err != nil {
		log.Printf("Failed to add review to style %v, err: %v", id, err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to add your review",
			"User":  u,
		})
	}

	return c.Redirect(fmt.Sprintf("/style/%d", id), fiber.StatusSeeOther)
}
