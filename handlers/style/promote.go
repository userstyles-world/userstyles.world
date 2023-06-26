package style

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
)

func sendPromotionEmail(style *models.APIStyle, mod string) {
	user, err := models.FindUserByID(strconv.Itoa(int(style.UserID)))
	if err != nil {
		log.Warn.Printf("Couldn't find user %d: %s\n", style.UserID, err)
		return
	}

	args := fiber.Map{
		"User":       user,
		"ModName":    mod,
		"ModProfile": config.BaseURL + "/user/" + mod,
	}

	title := "Your style has been featured"
	err = email.Send("style/promote", user.Email, title, args)
	if err != nil {
		log.Warn.Printf("Failed to send an email: %s\n", err)
	}
}

func Promote(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	p := c.Params("id")

	// Only moderator and above have permissions to promote styles.
	if u.Role < models.Moderator {
		return c.Render("err", fiber.Map{
			"Title": "You are not authorized to perform this action",
			"User":  u,
		})
	}

	id, err := strconv.Atoi(p)
	if err != nil {
		log.Info.Printf("Failed to convert %s to int: %s\n", p, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Couldn't convert style ID",
			"User":  u,
		})
	}

	style, err := models.GetStyleByID(p)
	if err != nil {
		log.Warn.Println("Failed to get the style:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	err = database.Conn.
		Model(models.Style{}).
		Where("id = ?", p).
		Update("featured", !style.Featured).
		Error

	if err != nil {
		log.Warn.Printf("Failed to promote style %d: %s\n", style.ID, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Failed to promote a style",
			"User":  u,
		})
	}

	// Ahem!!! We don't save the new value of Featured to the current style.
	// So we have to reverse check it ;)
	if !style.Featured {
		go sendPromotionEmail(style, u.Username)

		// Create a notification.
		notification := models.Notification{
			Seen:     false,
			Kind:     models.KindStylePromotion,
			TargetID: int(style.UserID),
			UserID:   int(u.ID),
			StyleID:  id,
		}

		go func(notification models.Notification) {
			if err := notification.Create(); err != nil {
				log.Warn.Printf("Failed to create a notification for %d, err: %v", id, err.Error())
			}
		}(notification)
	}

	return c.Redirect("/style/"+p, fiber.StatusSeeOther)
}
