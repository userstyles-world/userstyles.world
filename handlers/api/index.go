package api

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

type USOFormat struct {
	Id         uint   `json:"i"`
	Name       string `json:"n"`
	Category   string `json:"c"`
	Author     string `json:"an"`
	Screenshot string `json:"sn"`
	IsUSWStyle bool   `json:"uw"`
}

func convertToUsoFormat(ob1 models.APIStyle) USOFormat {
	return USOFormat{
		Id:         ob1.ID,
		Name:       ob1.Name,
		Category:   fixCategory(ob1.Category),
		Author:     ob1.Username,
		Screenshot: ob1.Preview,
		IsUSWStyle: true,
	}
}

func fixCategory(cat string) string {
	if cat == "unset" {
		return "global"
	}
	if !strings.Contains(cat, ".") {
		return cat
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

	if format == "uso-format" {
		formattedStyles := make([]USOFormat, len(*styles))
		for i, style := range *styles {
			formattedStyles[i] = convertToUsoFormat(style)
		}

		return c.JSON(fiber.Map{
			"data": formattedStyles,
		})
	}

	return c.JSON(fiber.Map{
		"data": styles,
	})
}
