package core

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetDocs(c *fiber.Ctx) error {
	document := c.Params("document")

	var content []byte
	switch document {
	case "changelog":
		content, _ = os.ReadFile("docs/changelog.md")
	}

	if len(content) == 0 {
		return c.Render("err", fiber.Map{
			"Title": "Couldn't load the document.",
		})
	}

	return c.Render("docs", fiber.Map{
		"Title":   "Changelog",
		"content": string(content),
	})
}
