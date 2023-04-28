package core

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/search"
)

func Search(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Search userstyles")
	c.Locals("Canonical", "search")

	keyword := c.Query("q")
	if keyword == "" {
		return c.Render("core/search", fiber.Map{})
	}
	c.Locals("Keyword", keyword)

	page, err := models.IsValidPage(c.Query("page"))
	if err != nil || page < 1 {
		c.Locals("Title", "Invalid page size")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	sort := c.Query("sort")
	c.Locals("Sort", sort)

	s, m, err := search.FindStylesByText(keyword, sort, page, config.AppPageMaxItems)
	switch {
	case errors.Is(err, search.ErrSearchNoResults):
		c.Locals("Title", "No results found")
		c.Locals("Error", "No results found for <b>"+keyword+"</b>.")
		return c.Status(fiber.StatusNotFound).Render("core/search", fiber.Map{})
	case err != nil:
		c.Locals("Title", "Bad search request")
		c.Locals("Error", "Bad search request. Please try again.")
		return c.Status(fiber.StatusBadRequest).Render("core/search", fiber.Map{})
	}
	c.Locals("Styles", s)
	c.Locals("Metrics", m)

	p := models.NewPagination(page, m.Total, sort, c.Path())
	p.Query = keyword
	if p.OutOfBounds() {
		return c.Redirect(p.URL(p.Now), 302)
	}
	c.Locals("Pagination", p)

	return c.Render("core/search", fiber.Map{})
}
