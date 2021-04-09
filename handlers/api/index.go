package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

type USoFormat struct {
	ID         uint   `json:"i"`
	Name       string `json:"n"`
	Category   string `json:"c"`
	Author     string `json:"an"`
	Screenshot string `json:"sn"`
}

func convertToUSoFormat(s models.APIStyle) USoFormat {
	return USoFormat{
		ID:         s.ID,
		Name:       s.Name,
		Category:   fixCategory(s.Category),
		Author:     s.Username,
		Screenshot: s.Preview,
	}
}

func fixCategory(cat string) string {
	if cat == "unset" {
		return "global"
	}
	cat = strings.TrimSuffix(cat, ".com")
	cat = strings.TrimSuffix(cat, ".org")
	return cat

}

func GetStyleIndex(c *fiber.Ctx) error {
	format := c.Params("format")

	styles, err := models.GetAllStyles(database.DB)
	if err != nil {
		return c.JSON(fiber.Map{
			"data": "styles not found",
		})
	}

	// Used by Stylus integration.
	if format == "uso-format" {
		formattedStyles := make([]USoFormat, len(*styles))
		for i, style := range *styles {
			formattedStyles[i] = convertToUSoFormat(style)
		}

		return c.JSON(fiber.Map{
			"data": formattedStyles,
		})
	}

	return c.JSON(fiber.Map{
		"data": styles,
	})
}
