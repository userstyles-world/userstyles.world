package style

import (
	"sort"
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

	s, err := models.GetAllAvailableStylesPaginated(pageNow)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
			"User":  u,
		})
	}

	fv := c.Query("sort")
	var sortFunction func(i, j int) bool
	switch fv {
	case "newest":
		sortFunction = func(i, j int) bool { return s[i].CreatedAt.Unix() > s[j].CreatedAt.Unix() }
	case "oldest":
		sortFunction = func(i, j int) bool { return s[i].CreatedAt.Unix() < s[j].CreatedAt.Unix() }
	case "recentlyupdated":
		sortFunction = func(i, j int) bool { return s[i].UpdatedAt.Unix() > s[j].UpdatedAt.Unix() }
	case "leastupdated":
		sortFunction = func(i, j int) bool { return s[i].UpdatedAt.Unix() < s[j].UpdatedAt.Unix() }
	case "mostinstalls":
		sortFunction = func(i, j int) bool { return s[i].Installs > s[j].Installs }
	case "leastinstalls":
		sortFunction = func(i, j int) bool { return s[i].Installs < s[j].Installs }
	case "mostviews":
		sortFunction = func(i, j int) bool { return s[i].Views > s[j].Views }
	case "leastviews":
		sortFunction = func(i, j int) bool { return s[i].Views < s[j].Views }
	}
	if sortFunction != nil {
		sort.Slice(s, sortFunction)
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
