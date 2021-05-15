package core

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
)

func GetDocs(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	document := c.Params("document")

	var content []byte
	switch document {
	case "changelog":
		content, _ = os.ReadFile("docs/changelog.md")
	}

	if len(content) == 0 {
		return c.Render("err", fiber.Map{
			"Title": "Couldn't load the document.",
			"User":  u,
		})
	}

	return c.Render("docs", fiber.Map{
		"Title":   "Changelog",
		"User":    u,
		"content": string(content),
	})
}
