package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetStyleSource(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(database.DB, id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	// Override updateURL field for Stylus integration.
	// TODO: Also override it in the database on demand?
	uc := usercss.ParseFromString(style.Code)
	url := "https://userstyles.world/api/style/" + id + ".user.css"
	uc.OverrideUpdateURL(url)

	c.Set("Content-Type", "text/css")
	return c.SendString(uc.SourceCode)
}
