package style

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetExplore(c *fiber.Ctx) error {
	t := &models.Style{}
	q := &[]models.APIStyle{}
	err := database.DB.
		Debug().
		Model(t).
		Select("styles.*, u.username").
		Joins("join users u on u.id = styles.user_id").
		Find(q).
		Error

	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Styles not found",
		})
	}

	return c.SendString(fmt.Sprintf("%#+v", q))
}
