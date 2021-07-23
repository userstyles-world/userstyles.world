package core

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

func readFile(f string) (s string) {
	if !strings.HasPrefix(f, "docs/") {
		f = "docs/" + f
	}

	if !strings.HasSuffix(f, ".md") {
		f = f + ".md"
	}

	b, err := os.ReadFile(f)
	if err != nil {
		return ""
	}

	return string(b)
}

func GetDocs(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	var title, content string
	switch c.Params("document") {
	case "":
		content = readFile("readme")
		title = "Documentation"
	case "changelog":
		content = readFile("changelog")
		title = "Changelog"
	case "faq":
		content = readFile("faq")
		title = "Frequently Asked Questions"
	case "privacy":
		content = readFile("privacy")
		title = "Privacy Policy"
	case "security":
		content = readFile("security")
		title = "Security Policy"
	case "crypto":
		content = readFile("crypto.md")
		title = "Cryptography Usages"
	}

	if content == "" {
		return c.Render("err", fiber.Map{
			"Title": "Failed to load the document",
			"User":  u,
		})
	}

	return c.Render("core/docs", fiber.Map{
		"Title":     title,
		"User":      u,
		"content":   content,
		"Canonical": "docs",
	})
}
