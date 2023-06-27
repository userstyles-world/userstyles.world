package core

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

func Search(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)
	c.Locals("Title", "Search userstyles")
	c.Locals("Canonical", "search")

	keyword := strings.TrimSpace(c.Query("q"))
	if len(keyword) < 3 {
		c.Locals("Error", "Keywords need to be in a group of three or more characters.")
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

	t := time.Now()

	total, err := storage.TotalSearchStyles(keyword, sort)
	if err != nil {
		log.Database.Println(err)
		c.Locals("Title", "Failed to count userstyles")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	p := models.NewPagination(page, total, sort, c.Path())
	p.Query = keyword
	if p.OutOfBounds() {
		return c.Redirect(p.URL(p.Now), 302)
	}
	c.Locals("Pagination", p)

	s, err := storage.FindSearchStyles(keyword, p.SortStyles(), page)
	if err != nil {
		log.Database.Println(err)
		c.Locals("Title", "Failed to search for userstyles")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}
	c.Locals("Styles", s)

	m := struct {
		Total     int
		TimeSpent time.Duration
	}{
		Total:     total,
		TimeSpent: time.Since(t),
	}
	c.Locals("Metrics", m)

	return c.Render("core/search", fiber.Map{})
}
