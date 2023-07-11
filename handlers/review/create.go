package review

import (
	"fmt"
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

	r := models.NewReview(u.ID, s.ID, c.FormValue("rating"), c.FormValue("comment"))
	c.Locals("Review", r)

	if err = r.Validate(); err != nil {
		c.Locals("Error", strings.ToTitle(err.Error()[:1])+err.Error()[1:]+".")
		return c.Render("review/create", fiber.Map{})
	}

	// Add review to database.
	if err = r.CreateForStyle(); err != nil {
		log.Warn.Printf("Failed to add review to style %d: %s\n", i, err)
		c.Locals("Title", "Failed to add your review")
		return c.Render("err", fiber.Map{})
	}

	err = r.FindLastForStyle(i, u.ID)
	if err != nil {
		log.Warn.Printf("Failed to get review for style %d from user %d: %s\n", i, u.ID, err)
	} else {
		// Create a notification.
		notification := models.Notification{
			Seen:     false,
			Kind:     models.KindReview,
			TargetID: int(s.UserID),
			UserID:   int(u.ID),
			StyleID:  i,
			ReviewID: int(r.ID),
		}

		if err := notification.Create(); err != nil {
			log.Warn.Printf("Failed to create a notification for review %d: %v\n", r.ID, err)
		}
	}

	return c.Redirect(fmt.Sprintf("/style/%d", i), fiber.StatusSeeOther)
}
