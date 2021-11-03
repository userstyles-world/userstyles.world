package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
)

func StyleGet(c *fiber.Ctx) error {
	u, _ := User(c)

	// /authorize_style tokens contains a positive StyleID
	if u.StyleID == 0 {
		return c.Status(400).
			JSON(fiber.Map{
				"data": "This token doesn't have permission to access this.",
			})
	}

	currentStyle, err := models.GetStyleByID(strconv.FormatUint(uint64(u.StyleID), 10))
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't find current style. HINT: Might be deleted.",
			})
	}

	return c.JSON(fiber.Map{
		"data": currentStyle,
	})
}
