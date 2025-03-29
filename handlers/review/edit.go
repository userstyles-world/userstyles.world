package review

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
)

func editPage(c *fiber.Ctx) error {
	return c.Render("review/create", fiber.Map{"Title": "Edit review"})
}

func editForm(c *fiber.Ctx) error {
	c.Locals("Title", "Edit review")

	r := c.Locals("Review").(*models.Review)
	r.Prepare(c.FormValue("rating"), c.FormValue("comment"))
	// s, _ := json.MarshalIndent(r, "", "    ")
	// fmt.Printf("r: %s\n", s)

	if err := r.Validate(); err != nil {
		// c.Locals("Error", strings.ToTitle(err.Error()[:1])+err.Error()[1:]+".")
		return c.Render("review/create", fiber.Map{
			"Error": "Validation error: " + err.Error(),
		})
	}

	if err := r.UpdateFromUser(); err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
			"Title": "Failed to update your review",
		})
	}

	a := models.NewSuccessAlert("Review has been updated.")
	cache.Store.Add("alert "+r.User.Username, a, time.Minute)

	return c.Redirect(r.Permalink())
}
