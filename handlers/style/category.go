package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/storage"
)

func GetCategory(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	c.Locals("Title", "Categories")
	c.Locals("User", u)

	page, err := models.IsValidPage(c.Query("page"))
	if err != nil {
		c.Locals("Title", "Invalid page size")
		return c.Render("err", fiber.Map{})
	}

	count, err := storage.CountStyleCategories()
	if err != nil {
		c.Locals("Title", "Failed to count categories")
		return c.Render("err", fiber.Map{})
	}

	p := models.NewPagination(page, count, "", c.Path())
	if p.OutOfBounds() {
		return c.Redirect(p.URL(p.Now), 302)
	}
	c.Locals("Pagination", p)

	cat, err := storage.GetStyleCategories(page, config.AppPageMaxItems)
	if err != nil {
		c.Locals("Title", "Failed to find categories")
		return c.Render("err", fiber.Map{})
	}
	c.Locals("Categories", cat)

	return c.Render("style/category", fiber.Map{})
}
