package style

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

func GetExplore(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	var p models.Pagination
	if err := p.ConvPage(c.Query("page")); err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Invalid page size",
			"User":  u,
		})
	}

	styleCount, err := models.GetStyleCount()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Failed to add pagination",
			"User":  u,
		})
	}

	// Set pagination.
	p.CalcItems(int(styleCount), config.AppPageMaxItems)
	p.Sort = c.Query("sort")
	if p.OutOfBounds() {
		r := fmt.Sprintf("/explore?page=%d", p.Now)
		if p.Sort != "" {
			r += "&sort=" + p.Sort
		}

		return c.Redirect(r, 302)
	}

	// Set sorting method.
	var orderFunction string
	switch p.Sort {
	case "newest":
		orderFunction = "styles.created_at DESC"
	case "oldest":
		orderFunction = "styles.created_at ASC"
	case "recentlyupdated":
		orderFunction = "styles.updated_at DESC"
	case "leastupdated":
		orderFunction = "styles.updated_at ASC"
	case "mostinstalls":
		orderFunction = "installs DESC"
	case "leastinstalls":
		orderFunction = "installs ASC"
	case "mostviews":
		orderFunction = "views DESC"
	case "leastviews":
		orderFunction = "views ASC"
	default:
		orderFunction = "styles.id ASC"
	}

	s, err := models.GetAllAvailableStylesPaginated(p.Now, orderFunction)
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
