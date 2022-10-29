package user

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/storage"
)

func Profile(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	username := c.Params("name")

	profile, err := models.FindUserByName(username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	// Always redirect to correct URL.
	if username != profile.Username {
		return c.Redirect("/user/"+strings.ToLower(username), fiber.StatusSeeOther)
	}

	var p models.Pagination
	if err := p.ConvPage(c.Query("page")); err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Invalid page size",
			"User":  u,
		})
	}

	count, err := storage.CountStylesForUserID(profile.ID)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Failed to count userstyles",
			"User":  u,
		})
	}

	// Set pagination.
	size := config.AppPageMaxItems
	p.CalcItems(count, size)
	p.Sort = c.Query("sort")
	p.Path = c.Path()
	styles, err := storage.FindStyleCardsPaginatedForUserID(
		p.Now, size, p.SortStyles(), profile.ID)
	if err != nil {
		return c.Render("err", fiber.Map{
			"User":  u,
			"Title": "Server error",
		})
	}

	// Render Account template if current user matches active session.
	/*if u.Username == username {
		return c.Render("user/account", fiber.Map{
			"Title":  "Account",
			"User":   u,
			"Styles": styles,
		})
	}*/

	return c.Render("user/profile", fiber.Map{
		"Title":     profile.Name() + "'s profile",
		"User":      u,
		"Profile":   profile,
		"Styles":    styles,
		"Count":     count,
		"Canonical": "user/" + username,
		"Sort":      p.Sort,
		"P":         p,
	})
}
