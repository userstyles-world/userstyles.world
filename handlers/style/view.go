package style

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

func GetStyle(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	data, err := models.GetStyleByID(database.DB, id)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	var total int64
	database.DB.
		Debug().
		Model(models.Stats{}).
		Where("style_id = ?", id).
		Count(&total)

	var week int64
	database.DB.
		Debug().
		Model(models.Stats{}).
		Where("style_id = ? and updated_at > ?", id, time.Now().AddDate(0, 0, -7)).
		Count(&week)

	return c.Render("style", fiber.Map{
		"Title": data.Name,
		"User":  u,
		"Style": data,
		"Total": total,
		"Week":  week,
		"Url":   fmt.Sprintf("https://userstyles.world/style/%d", data.ID),
	})
}
