package core

import (
	"io"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/log"
	"userstyles.world/modules/markdown"
	"userstyles.world/web"
)

func readFile(f string) (s string) {
	if !strings.HasPrefix(f, "docs/") {
		f = "docs/" + f
	}

	if !strings.HasSuffix(f, ".md") {
		f += ".md"
	}

	b, err := os.ReadFile(f)
	if err != nil {
		return ""
	}

	return string(b)
}

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

	// TODO: Extract metadata.
	var title string
	switch doc {
	case "":
		title = "Documentation"
	case "changelog":
		title = "Changelog"
	case "faq":
		title = "Frequently Asked Questions"
	case "privacy":
		title = "Privacy Policy"
	case "security":
		title = "Security Policy"
	case "crypto":
		title = "Cryptography Usages"
	case "content-guidelines":
		title = "Content Guidelines"
	}

	b, _ := io.ReadAll(f)
	content := markdown.RenderSafe(b)

	return c.Render("core/docs", fiber.Map{
		"Title":     title,
		"User":      u,
		"content":   content,
		"Canonical": "docs",
	})
}
