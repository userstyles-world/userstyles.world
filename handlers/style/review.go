package style

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func ReviewGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	// Prevent spam.
	reviewSpam := new(models.Review)
	if err := reviewSpam.FindLastFromUser(id, u.ID); err != nil {
		log.Printf("Failed to find last review for style %v and user %v\n", id, u.ID)
	}

	if reviewSpam.ID > 0 {
		now := time.Now().Add(-24 * 7 * time.Hour)
		if now.Before(reviewSpam.CreatedAt) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("err", fiber.Map{
				"Title": "You can post only one review per week",
				"User":  u,
			})
		}
	}

	return c.Render("style/review", fiber.Map{
		"Title": "Review style",
		"User":  u,
		"ID":    id,
	})
}

func ReviewPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	cmt := c.FormValue("comment")

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Invalid style ID",
			"User":  u,
		})
	}

	// Check if style exists.
	style, err := models.GetStyleByID(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Prevent spam.
	reviewSpam := new(models.Review)
	if err := reviewSpam.FindLastFromUser(id, u.ID); err != nil {
		fmt.Printf("Failed to find last review for style %v and user %v\n", id, u.ID)
	}

	if reviewSpam.ID > 0 {
		now := time.Now().Add(-7 * 24 * time.Hour)
		if now.Before(reviewSpam.CreatedAt) {
			c.Status(fiber.StatusUnauthorized)
			return c.Render("err", fiber.Map{
				"Title": "You can post only one review per week",
				"User":  u,
			})
		}
	}

	r, err := strconv.Atoi(c.FormValue("rating"))
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Invalid style ID",
			"User":  u,
		})
	}

	// Check if rating is out of range.
	if r < 1 || r > 5 {
		return c.Render("style/review", fiber.Map{
			"Title":   "Review style",
			"User":    u,
			"ID":      id,
			"Error":   "Rating is out of range",
			"Rating":  r,
			"Comment": cmt,
		})
	}

	// Prevent spam.
	if len(cmt) > 500 {
		return c.Render("style/review", fiber.Map{
			"Title":   "Review style",
			"User":    u,
			"ID":      id,
			"Error":   "Comment can't be longer than 500 characters",
			"Rating":  r,
			"Comment": cmt,
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

	if err = review.FindLastForStyle(id, u.ID); err != nil {
		log.Printf("Failed to find review for style %v, err: %v", id, err)
	} else {
		// Create a notification.
		notification := models.Notification{
			Seen:     false,
			Kind:     models.KindReview,
			TargetID: int(style.UserID),
			UserID:   int(u.ID),
			StyleID:  id,
			ReviewID: int(review.ID),
		}

		go func(notification models.Notification) {
			if err := notification.Create(); err != nil {
				log.Printf("Failed to create a notification for %d, err: %v", id, err)
			}
		}(notification)
	}

	return c.Redirect(fmt.Sprintf("/style/%d", id), fiber.StatusSeeOther)
}
