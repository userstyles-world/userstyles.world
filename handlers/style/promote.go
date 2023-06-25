package style

import (
	"bytes"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
)

func sendPromotionEmail(c *fiber.Ctx, userID uint, style *models.APIStyle, modName, baseURL string) {
	user, err := models.FindUserByID(strconv.Itoa(int(userID)))
	if err != nil {
		log.Warn.Printf("Couldn't find user %d: %s", userID, err.Error())
		return
	}

	modProfile := baseURL + "/user/" + modName

	args := fiber.Map{
		"ModName":    modName,
		"ModProfile": modProfile,
	}

	var bufText, bufHTML bytes.Buffer
	err = email.Render(&bufText, &bufHTML, "stylepromoted", args)
	if err != nil {
		log.Warn.Printf("Failed to render email template: %v\n", err)
		return
	}

	err = utils.NewEmail().
		SetTo(user.Email).
		SetSubject("Your style is being featured").
		AddPart(*utils.NewPart().SetBody(bufText.String())).
		AddPart(*utils.NewPart().SetBody(bufHTML.String()).HTML()).
		SendEmail(config.IMAPServer)
	if err != nil {
		log.Warn.Println("Failed to send email:", err.Error())
		return
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
		go sendPromotionEmail(c, style.UserID, style, u.Username, c.BaseURL())

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
