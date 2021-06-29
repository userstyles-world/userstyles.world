package user

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
)

func Ban(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	if !u.IsAdmin() {
		return c.Render("err", fiber.Map{
			"Title": "Invalid",
			"User":  u,
		})
	}

	user, err := models.FindUserByID(id)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User ID doesn't exist",
			"User":  u,
		})
	}

	return c.Render("user/ban", fiber.Map{
		"Title":  "Ban user",
		"User":   u,
		"Params": user,
	})
}

func ConfirmBan(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	if !u.IsAdmin() {
		return c.Render("err", fiber.Map{
			"Title": "Invalid",
			"User":  u,
		})
	}

	err := database.Conn.
		Debug().
		Delete(&models.User{}, "id = ?", id).
		Error

	if err != nil {
		log.Printf("Failed to ban user %s, err: %s", id, err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	err = database.Conn.
		Debug().
		Delete(&models.Style{}, "user_id = ?", id).
		Error

	if err != nil {
		log.Printf("Failed to ban styles from user %s, err: %s", id, err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect("/dashboard")
}
