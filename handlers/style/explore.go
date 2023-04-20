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

	page, err := models.IsValidPage(c.Query("page"))
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Invalid page size",
			"User":  u,
		})
	}

	styleCount, err := models.GetStyleCount()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Failed to count userstyles",
			"User":  u,
		})
	}

	// Set pagination.
	p := models.NewPagination(page, c.Query("sort"), c.Path())
	p.CalcItems(int(styleCount))
	if p.OutOfBounds() {
		return c.Redirect(p.URL(p.Now), 302)
	}

	// Query for [sorted] styles.
	s, err := storage.FindStyleCardsPaginated(p.Now, config.AppPageMaxItems, p.SortStyles())
	if err != nil {
		log.Warn.Println("Failed to get paginated styles:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

	return c.Render("core/explore", fiber.Map{
		"Title":     "Explore website themes",
		"User":      u,
		"Styles":    s,
		"Sort":      p.Sort,
		"P":         p,
		"Canonical": "explore",
	})
}
