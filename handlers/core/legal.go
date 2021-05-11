package core

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func GetLegal(c *fiber.Ctx) error {
	document := c.Params("document")

	var content []byte
	switch document {
	case "privacy-policy":
		content, _ = os.ReadFile("privacy-policy.md")
	case "terms-of-service":
		content, _ = os.ReadFile("terms-of-service.md")
	}

	if len(content) == 0 {
		return c.Render("err", fiber.Map{
			"Title": "Couldn't load the document.",
		})
	}

	return c.Render("legal", fiber.Map{
		"content": string(content),
	})
}
