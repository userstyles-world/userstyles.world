package user

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

func Profile(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	username := c.Params("name")
	profile, err := models.FindUserByName(username)
	if err != nil {
		c.Locals("Title", "User not found")
		return c.Render("err", fiber.Map{})
	}
	c.Locals("Profile", profile)
	c.Locals("Canonical", "user/"+username)

	// Always redirect to correct URL.
	if username != profile.Username {
		return c.Redirect("/user/"+profile.Username, fiber.StatusSeeOther)
	}
	c.Locals("Title", profile.Name()+"'s profile")

	page, err := models.IsValidPage(c.Query("page"))
	if err != nil {
		c.Locals("Title", "Invalid page size")
		return c.Render("err", fiber.Map{})
	}

	count, err := storage.CountStylesForUserID(profile.ID)
	if err != nil {
		c.Locals("Title", "Failed to count userstyles")
		return c.Render("err", fiber.Map{})
	}
	c.Locals("Count", count)

	size := config.App.PageMaxItems
	p := models.NewPagination(page, count, c.Query("sort"), c.Path())
	if p.OutOfBounds() {
		return c.Redirect(p.URL(p.Now), 302)
	}
	c.Locals("Pagination", p)
	c.Locals("Sort", p.Sort)

	styles, err := storage.FindStyleCardsPaginatedForUserID(
		p.Now, size, p.SortStyles(), profile.ID)
	if err != nil {
		log.Database.Println("Failed to get styles:", err)
		c.Locals("Title", "Styles not found")
		return c.Render("err", fiber.Map{})
	}
	c.Locals("Styles", styles)

	return c.Render("user/profile", fiber.Map{})
}
