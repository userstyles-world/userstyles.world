package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
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

	// Adjust pagination numbers.
	p.CalcItems(int(styleCount), 40)

	// Set sort query in pagination.
	fv := c.Query("sort")
	p.Sort = fv

	// Set sorting method.
	var orderFunction string
	switch fv {
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
		"Title":     "Explore",
		"User":      u,
		"Styles":    s,
		"Sort":      fv,
		"P":         p,
		"Canonical": "explore",
	})
}
