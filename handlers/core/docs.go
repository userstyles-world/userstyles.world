package core

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

func GetDocs(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	var title string
	var content []byte

	switch c.Params("document") {
	case "changelog":
		content, _ = os.ReadFile("docs/changelog.md")
		title = "Changelog"
	case "privacy-policy":
		content, _ = os.ReadFile("docs/privacy-policy.md")
		title = "Privacy Policy"
	}

	if len(content) == 0 {
		return c.Render("err", fiber.Map{
			"Title": "Couldn't load the document.",
			"User":  u,
		})
	}

	return c.Render("docs", fiber.Map{
		"Title":   title,
		"User":    u,
		"content": string(content),
	})
}
