package review

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/log"
)

func createPage(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Review style")

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		c.Locals("Title", "Invalid style ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}
	c.Locals("ID", i)

	s, err := models.GetStyleByID(c.Params("id"))
	if err != nil {
		c.Locals("Title", "Style not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	// Prevent authors reviewing their own userstyles.
	if u.ID == s.UserID {
		c.Locals("Title", "Can't review your own style")
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{})
	}

	// Prevent spamming reviews.
	dur, ok := models.AbleToReview(u.ID, s.ID)
	if !ok {
		c.Locals("Title", "Cannot review style")
		c.Locals("ErrTitle", "You can review this style again in "+dur)
		return c.Status(fiber.StatusTooManyRequests).Render("err", fiber.Map{})
	}

	return c.Render("review/create", fiber.Map{})
}

func createForm(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Review style")

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		c.Locals("Title", "Invalid style ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}
	c.Locals("ID", i)

	s, err := models.GetStyleByID(c.Params("id"))
	if err != nil {
		c.Locals("Title", "Style not found")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	// Prevent authors reviewing their own userstyles.
	if u.ID == s.UserID {
		c.Locals("Title", "Can't review your own style")
		return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{})
	}

	// Prevent spamming reviews.
	dur, ok := models.AbleToReview(u.ID, s.ID)
	if !ok {
		c.Locals("Title", "Cannot review style")
		c.Locals("ErrTitle", "You can review this style again in "+dur)
		return c.Status(fiber.StatusTooManyRequests).Render("err", fiber.Map{})
	}

	r, err := strconv.Atoi(c.FormValue("rating"))
	if err != nil {
		c.Locals("Error", "Rating should be a number")
		return c.Render("review/create", fiber.Map{})
	}
	c.Locals("Rating", r)

	cmt := strings.TrimSpace(c.FormValue("comment"))
	c.Locals("Comment", cmt)

	// Check if rating is out of range.
	if r < 0 || r > 5 {
		c.Locals("Error", "Rating is out of range.")
		return c.Status(fiber.StatusBadRequest).Render("review/create", fiber.Map{})
	}

	// Prevent spam.
	if len(cmt) > 500 {
		c.Locals("Error", "Comment can't be longer than 500 characters.")
		return c.Status(fiber.StatusBadRequest).Render("review/create", fiber.Map{})
	}

	if r == 0 && len(cmt) == 0 {
		c.Locals("Error", "You can't make empty reviews. Please insert a rating and/or a comment.")
		return c.Status(fiber.StatusBadRequest).Render("review/create", fiber.Map{})
	}

	// Create review.
	review := models.Review{
		Rating:  r,
		Comment: cmt,
		StyleID: i,
		UserID:  int(u.ID),
	}

	// Add review to database.
	if err := review.CreateForStyle(); err != nil {
		log.Warn.Printf("Failed to add review to style %d: %s\n", i, err)
		c.Locals("Title", "Failed to add your review")
		return c.Render("err", fiber.Map{})
	}

	if err = review.FindLastForStyle(i, u.ID); err != nil {
		log.Warn.Printf("Failed to get review for style %d from user %d: %s\n", i, u.ID, err)
	} else {
		// Create a notification.
		notification := models.Notification{
			Seen:     false,
			Kind:     models.KindReview,
			TargetID: int(s.UserID),
			UserID:   int(u.ID),
			StyleID:  i,
			ReviewID: int(review.ID),
		}

		if err := notification.Create(); err != nil {
			log.Warn.Printf("Failed to create a notification for review %d: %v\n", review.ID, err)
		}
	}

	return c.Redirect(fmt.Sprintf("/style/%d", i), fiber.StatusSeeOther)
}
