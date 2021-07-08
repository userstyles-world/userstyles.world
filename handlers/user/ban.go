package user

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
)

func Ban(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id") // TODO: Switch to int type.

	if !u.IsModOrAdmin() {
		return c.Render("err", fiber.Map{
			"Title": "Unauthorized",
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

	if u.ID == user.ID {
		return c.Render("err", fiber.Map{
			"Title": "You can't ban yourself",
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
	stringID := c.Params("id")
	id, _ := strconv.Atoi(stringID)
	reason := c.FormValue("reason")

	if !u.IsModOrAdmin() {
		return c.Render("err", fiber.Map{
			"Title": "Unauthorized",
			"User":  u,
		})
	}

	if int(u.ID) == id {
		return c.Render("err", fiber.Map{
			"Title": "You can't ban yourself",
			"User":  u,
		})
	}

	targetUser, err := models.FindUserByID(stringID)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User ID doesn't exist",
			"User":  u,
		})
	}

	err = database.Conn.
		Debug().
		Delete(&models.User{}, "id = ?", id).
		Error

	if err != nil {
		log.Printf("Failed to ban user %d, err: %s", id, err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	// Add banned user log entry.
	err = database.Conn.
		Debug().
		Create(&models.Log{
			UserID:         u.ID,
			TargetData:     "",
			Reason:         reason,
			Kind:           models.LogBanUser,
			TargetUserName: targetUser.Username,
		}).
		Error

	if err != nil {
		log.Printf("Failed to add log entry!!! user %d, err: %s", id, err)
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
		log.Printf("Failed to ban styles from user %d, err: %s", id, err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Redirect("/dashboard")
}
