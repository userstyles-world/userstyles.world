package style

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/log"
)

func ReviewGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	// Check if style exists.
	style, err := models.GetStyleByID(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Check if user isn't style's author.
	if u.ID == style.UserID {
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
			"Title": "Cannot review own style",
			"User":  u,
		})
	}

	// Prevent spam.
	reviewSpam := new(models.Review)
	// Collecting of the error is not needed.
	// As we simply check "valid" data by checking if ID is a positive integer.
	if _ = reviewSpam.FindLastFromUser(id, u.ID); reviewSpam.ID > 0 {
		log.Info.Printf("User %d tried to review style %v more than once.\n", u.ID, id)
		week := time.Now().Add(-7 * 24 * time.Hour)
		if reviewSpam.CreatedAt.After(week) {
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

	// Check if user isn't style's author.
	if u.ID == style.UserID {
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
			"Title": "Cannot review own style",
			"User":  u,
		})
	}

	// Prevent spam.
	reviewSpam := new(models.Review)
	// Collecting of the error is not needed.
	// As we simply check "valid" data by checking if ID is a positive integer.
	if _ = reviewSpam.FindLastFromUser(id, u.ID); reviewSpam.ID > 0 {
		log.Warn.Printf("User %d tried to review style %v more than once.\n", u.ID, id)
		week := time.Now().Add(-7 * 24 * time.Hour)
		if reviewSpam.CreatedAt.After(week) {
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
	if r < 0 || r > 5 {
		return c.Render("style/review", fiber.Map{
			"Title":   "Review style",
			"User":    u,
			"ID":      id,
			"Error":   "Rating is out of range.",
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
			"Error":   "Comment can't be longer than 500 characters.",
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
	if err := review.CreateForStyle(); err != nil {
		log.Warn.Printf("Failed to add review to style %v: %v\n", id, err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to add your review",
			"User":  u,
		})
	}

	if err = review.FindLastForStyle(id, u.ID); err != nil {
		log.Warn.Printf("Failed to get review for style %v from user %v: %v\n", id, u.ID, err)
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
				log.Warn.Printf("Failed to create a notification for review %d: %v\n", review.ID, err)
			}
		}(notification)
	}

	return c.Redirect(fmt.Sprintf("/style/%d", id), fiber.StatusSeeOther)
}
