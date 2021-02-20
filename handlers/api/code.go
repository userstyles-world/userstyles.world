package api

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetStyleSource(c *fiber.Ctx) error {
	id := c.Params("id")

	t, q := new(models.Style), new(models.APIStyle)
	err := database.DB.
		Debug().
		Model(t).
		Select("styles.*, u.username").
		Joins("join users u on u.id = styles.user_id").
		First(q, "styles.id = ?", id).
		Error

	if err != nil {
		log.Printf("Problem with style id %s, err: %v", id, err)
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	c.Set("Content-Type", "text/css")
	return c.SendString(fmt.Sprintf("%s", q.Code))
}
