package api

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

type USoFormat struct {
	ID             uint   `json:"i"`
	Name           string `json:"n"`
	Category       string `json:"c"`
	UpdatedAt      int64  `json:"u"` // Requires Unix timestamp
	TotalInstalls  int64  `json:"t"`
	WeeklyInstalls int64  `json:"w"`
	Author         string `json:"an"`
	Screenshot     string `json:"sn"`
}

func convertToUSoFormat(s models.APIStyle) USoFormat {
	id := fmt.Sprintf("%d", s.ID) // Convert uint to string.

	var img string
	if s.Preview != "" {
		img = fmt.Sprintf("https://userstyles.world/api/style/preview/%d.webp", s.ID)
	}

	return USoFormat{
		ID:             s.ID,
		Name:           s.Name,
		Category:       fixCategory(s.Category),
		Author:         s.Username,
		Screenshot:     img,
		UpdatedAt:      s.UpdatedAt.Unix(),
		TotalInstalls:  models.GetTotalInstallsForStyle(database.DB, id),
		WeeklyInstalls: models.GetWeeklyInstallsForStyle(database.DB, id),
	}
}

func fixCategory(cat string) string {
	if cat == "unset" {
		return "global"
	}
	cat = strings.ToLower(cat)

	if strings.HasSuffix(cat, ".com") || strings.HasSuffix(cat, ".org") {
		cat = strings.TrimSuffix(cat, ".com")
		cat = strings.TrimSuffix(cat, ".org")
		// Remove any subdomain
		// web.whatsapp -> whatsapp
		if strings.Count(cat, ".") >= 1 {
			cat = strings.Split(cat, ".")[1]
		}
	} else {
		if strings.Count(cat, ".") >= 2 {
			cat = strings.Join(strings.Split(cat, ".")[1:], ".")
		}
	}

	return cat
}

func GetStyleIndex(c *fiber.Ctx) error {
	format := c.Params("format")

	styles, err := models.GetAllStylesForIndexAPI(database.DB)
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
