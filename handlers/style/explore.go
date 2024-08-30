package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

func GetExplore(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Explore custom website themes")
	c.Locals("Canonical", "explore")
	c.Locals("ExplorePage", true)

	page, err := models.IsValidPage(c.Query("page"))
	if err != nil {
		c.Locals("Title", "Invalid page size")
		return c.Render("err", fiber.Map{})
	}

	count, err := models.GetStyleCount()
	if err != nil {
		c.Locals("Title", "Failed to count userstyles")
		return c.Render("err", fiber.Map{})
	}

	sort := c.Query("sort")
	c.Locals("Sort", sort)

	p := models.NewPagination(page, count, sort, c.Path())
	if p.OutOfBounds() {
		return c.Redirect(p.URL(p.Now), 302)
	}
	c.Locals("Pagination", p)

	// Query for [sorted] styles.
	s, err := storage.FindStyleCardsPaginated(p.Now, config.App.PageMaxItems, p.SortStyles())
	if err != nil {
		log.Database.Println("Failed to get styles:", err)
		c.Locals("Title", "Styles not found")
		return c.Render("err", fiber.Map{})
	}
	c.Locals("Styles", s)

	return c.Render("core/explore", fiber.Map{})
}
