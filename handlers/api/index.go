package api

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

func getUSoIndex(c *fiber.Ctx) error {
	index, found := cache.Store.Get("index")
	if !found {
		var err error
		index, err = storage.GetStyleCompactIndex(database.Conn)
		if err != nil {
			log.Warn.Printf("Failed to get compact index: %s\n", err)
			return c.JSON(fiber.Map{"data": "index not found"})
		}

		// Set cache for index endpoint.
		cache.Store.Set("index", index, 0)
	}

	c.Set("Content-Type", "application/json")
	return c.Send(index.([]byte))
}

func getFullIndex(c *fiber.Ctx) error {
	styles, err := models.GetAllStylesForIndexAPI()
	if err != nil {
		return c.JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": styles,
	})
}

func GetStyleIndex(c *fiber.Ctx) error {
	switch c.Params("format") {
	case "uso-format":
		return getUSoIndex(c)
	default:
		return getFullIndex(c)
	}
}
