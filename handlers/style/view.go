package style

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
)

func GetStyle(c *fiber.Ctx) error {
	t := &models.Style{}
	q := &models.APIStyle{}
	err := database.DB.
		Debug().
		Model(t).
		Select("styles.*,  u.username").
		Joins("join users u on u.id = styles.user_id").
		Find(q, "styles.id = ?", c.Params("id")).
		Error

	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
		})
	}

	return c.Render("style", fiber.Map{
		"Title": q.Name,
		"Style": q,
	})
}
