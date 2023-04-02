package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
)

func DeleteGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	m := fiber.Map{"User": u}

	id, err := c.ParamsInt("id")
	if err != nil {
		m["Title"] = "Invalid user ID"
		return c.Render("err", m)
	}

	if u.ID != uint(id) {
		m["Title"] = "You can't delete other users"
		return c.Render("err", m)
	}

	m["Title"] = "Delete account"

	return c.Render("user/delete", m)
}

func DeletePost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	m := fiber.Map{"User": u, "Title": "Delete account"}

	id, err := c.ParamsInt("id")
	if err != nil {
		m["Title"] = "Invalid user ID"
		return c.Render("err", m)
	}

	if u.ID != uint(id) {
		m["Title"] = "You can't delete other users"
		return c.Render("err", m)
	}

	// Delete from database.
	err = database.Conn.Transaction(func(tx *gorm.DB) error {
		if err = tx.Debug().Delete(&models.User{}, "id = ?", id).Error; err != nil {
			return err
		}

		// TODO: Remove after introducing cascading deletes.
		if err = tx.Debug().Delete(&models.Style{}, "user_id = ?", id).Error; err != nil {
			return err
		}
		if err = tx.Debug().Delete(&models.Review{}, "user_id = ?", id).Error; err != nil {
			return err
		}
		if err = tx.Debug().Delete(&models.Notification{}, "user_id = ?", id).Error; err != nil {
			return err
		}
		if err = tx.Debug().Delete(&models.ExternalUser{}, "user_id = ?", id).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Database.Printf("Failed to delete %q: %s\n", u.Username, err)
		m["Title"] = "Failed to delete you account. Please try again"
		return c.Render("err", m)
	}

	return c.Redirect("/logout")
}
