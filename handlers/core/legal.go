package core

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetLegal(c *fiber.Ctx) error {
	document := c.Params("document")

	var title string
	var content []byte
	switch document {
	case "privacy-policy":
		content, _ = os.ReadFile("docs/privacy-policy.md")
		title = "Privacy Policy"
	case "terms-of-service":
		content, _ = os.ReadFile("docs/terms-of-service.md")
		title = "Terms of Service"
	}

	if len(content) == 0 {
		return c.Render("err", fiber.Map{
			"Title": "Couldn't load the document.",
		})
	}

	return c.Render("legal", fiber.Map{
		"Title":   title,
		"content": string(content),
	})
}
