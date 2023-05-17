package core

import (
	"io"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/log"
	"userstyles.world/modules/markdown"
	"userstyles.world/web"
)

func GetDocs(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	doc := c.Params("document")
	if doc == "" {
		doc = "readme"
	}

	f, err := web.DocsDir.Open(doc + ".md")
	if err != nil {
		log.Info.Printf("Failed to load document %q: %s\n", doc, err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to load document",
			"User":  u,
		})
	}

	b, err := io.ReadAll(f)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Failed to read document",
			"User":  u,
		})
	}

	content, meta := markdown.RenderDocs(b)
	if len(content) == 0 {
		return c.Render("err", fiber.Map{
			"Title": "Failed to render document",
			"User":  u,
		})
	}

	return c.Render("core/docs", fiber.Map{
		"Title":     meta["Title"],
		"User":      u,
		"Content":   content,
		"Canonical": "docs",
		"meta":      meta,
	})
}
