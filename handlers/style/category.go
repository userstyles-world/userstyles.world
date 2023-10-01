package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/storage"
)

func GetCategory(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	c.Locals("Title", "Categories")
	c.Locals("User", u)

	cat, err := storage.GetStyleCategories()
	if err != nil {
		c.Locals("Title", "Failed to find categories")
		return c.Render("err", fiber.Map{})
	}
	c.Locals("Categories", cat)

	return c.Render("style/category", fiber.Map{})
}
