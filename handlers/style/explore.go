package style

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func GetExplore(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	page := c.Query("page")

	var pageNow int
	if page != "" {
		i, err := strconv.Atoi(page)
		if err != nil {
			return c.Render("err", fiber.Map{
				"Title": "Invalid page size",
				"User":  u,
			})
		}
		pageNow = i
	} else {
		pageNow = 1
	}

	styleCount, err := models.GetStyleCount()
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Failed to add pagination",
			"User":  u,
		})
	}

	// Adjust max pages.
	maxPages, remainder := int(styleCount)/40, styleCount%40
	if remainder > 0 {
		maxPages++
	}

	// If the page is greater than the max pages, display the last page.
	// Or if the page is less than 1, display the first page.
	if pageNow > maxPages {
		pageNow = maxPages
	}
	if pageNow < 1 {
		pageNow = 1
	}
	fv := c.Query("sort")
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

	s, err := models.GetAllAvailableStylesPaginated(pageNow, orderFunction)
	if err != nil {
		log.Println("Couldn't get paginated styles, ", err)
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

	return c.Render("core/explore", fiber.Map{
		"Title":    "Explore",
		"User":     u,
		"Styles":   s,
		"Sort":     fv,
		"PageMax":  maxPages,
		"PageNow":  pageNow,
		"PageBack": pageNow - 1,
		"PageNext": pageNow + 1,
	})
}
