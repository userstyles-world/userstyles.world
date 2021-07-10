package core

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

func readFile(f string) (s string) {
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
	case "changelog":
		content = readFile("docs/changelog.md")
		title = "Changelog"
	case "privacy":
		content = readFile("docs/privacy.md")
		title = "Privacy Policy"
	}

	if content == "" {
		return c.Render("err", fiber.Map{
			"Title": "Couldn't load the document.",
			"User":  u,
		})
	}

	return c.Render("core/docs", fiber.Map{
		"Title":   title,
		"User":    u,
		"content": string(content),
	})
}
