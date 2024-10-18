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

func sendPromotionEmail(style *models.Style, user *models.User, mod string) {
	args := fiber.Map{
		"User":      user,
		"Style":     style,
		"StyleLink": config.App.BaseURL + "/style/" + strconv.Itoa(int(style.ID)),
		"ModName":   mod,
		"ModLink":   config.App.BaseURL + "/user/" + mod,
	}

	title := "Your style has been featured"
	if err := email.Send("style/promote", user.Email, title, args); err != nil {
		log.Warn.Printf("Failed to email %d: %s\n", user.ID, err)
	}
}

func Promote(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Only moderator and above have permissions to promote styles.
	if u.Role < models.Moderator {
		return c.Render("err", fiber.Map{
			"Title": "You are not authorized to perform this action",
			"User":  u,
		})
	}

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Invalid style ID",
			"User":  u,
		})
	}

	style, err := models.GetStyleByID(i)
	if err != nil {
		log.Warn.Println("Failed to get the style:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	user, err := models.FindUserByID(strconv.Itoa(int(style.UserID)))
	if err != nil {
		log.Warn.Printf("Failed to find user %d: %s\n", style.UserID, err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	err = database.Conn.
		Model(models.Style{}).
		Where("id = ?", i).
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
		go sendPromotionEmail(style, user, u.Username)

		n := models.Notification{
			Seen:     false,
			Kind:     models.KindStylePromotion,
			TargetID: int(style.UserID),
			UserID:   int(u.ID),
			StyleID:  i,
		}

		if err := models.CreateNotification(database.Conn, &n); err != nil {
			log.Warn.Printf("Failed to create a notification for %d: %s\n", i, err)
		}
	}

	return c.Redirect("/style/"+c.Params("id"), fiber.StatusSeeOther)
}
