package style

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/modules/search"
	"userstyles.world/utils"
)

func BanGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Check if logged-in user has permissions.
	if !u.IsModOrAdmin() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("err", fiber.Map{
			"Title": "You are not authorized to perform this action",
			"User":  u,
		})
	}

	// Check if style exists.
	s, err := models.GetStyleByID(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	return c.Render("style/ban", fiber.Map{
		"Title": "Confirm ban",
		"User":  u,
		"Style": s,
	})
}

func sendBanEmail(baseURL string, user *models.User, style *models.APIStyle, modLogID uint) error {
	modLogEntry := baseURL + "/modlog#id-" + strconv.Itoa(int(modLogID))

	partPlain := utils.NewPart().
		SetBody("Hi " + user.Username + ",\n" +
			"We'd like to notice you about a recent action from our moderation team:\n\n" +
			"Your style \"" + style.Name + "\" has been removed from our platform.\n" +
			"You can check for more information about this action on the modlog: " + modLogEntry + "\n\n" +
			"If you'd like to come in touch with us, please email us at feedback@userstyles.world\n" +
			"Regards,\n" + "The Moderation Team")
	partHTML := utils.NewPart().
		SetBody("<p>Hi " + user.Username + ",</p>\n" +
			"<br>" +
			"<p>We'd like to notice you about a recent action from our moderation team:</p>\n" +
			"<br><br>" +
			"<p>Your style \"<b>" + style.Name + "</b>\" has been removed from our platform.</p>\n" +
			"<p>You can check for more information about this action on the " +
			"<a target=\"_blank\" clicktracking=\"off\" href=\"" + modLogEntry + "\">Modlog</a>.</p>\n" +
			"<p>If you'd like to come in touch with us, " +
			"please email us at <a href=\"mailto:feedback@userstyles.world\">feedback@userstyles.world</a>.<p>\n" +
			"<br><br>" +
			"<p>Regards,</p>\n" + "<p>The Moderation Team</p>").
		SetContentType("text/html")

	err := utils.NewEmail().
		SetTo(user.Email).
		SetSubject("Moderation notice").
		AddPart(*partPlain).
		AddPart(*partHTML).
		SendEmail(config.IMAPServer)
	if err != nil {
		return err
	}
	return nil
}

func BanPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	id := c.Params("id")

	// Check if logged-in user has permissions.
	if !u.IsModOrAdmin() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("err", fiber.Map{
			"Title": "You are not authorized to perform this action",
			"User":  u,
		})
	}

	// Check if style exists.
	s, err := models.GetStyleByID(id)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	// Initialize modlog data.
	logEntry := models.Log{
		UserID:         u.ID,
		Username:       u.Username,
		Kind:           models.LogRemoveStyle,
		TargetUserName: s.Username,
		TargetData:     s.Name,
		Reason:         strings.TrimSpace(c.FormValue("reason")),
		Censor:         c.FormValue("censor") == "on",
	}

	// Add banned style log entry.
	modlog := new(models.Log)
	if err := modlog.AddLog(&logEntry); err != nil {
		log.Warn.Printf("Failed to add style %d to ModLog: %s", s.ID, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	// Delete style from database.
	q := new(models.Style)
	if err = database.Conn.Delete(q, "styles.id = ?", s.ID).Error; err != nil {
		log.Warn.Printf("Failed to delete style %d: %s\n", s.ID, err.Error())
		c.Status(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	// Delete stats from database.
	if err = new(models.Stats).Delete(s.ID); err != nil {
		log.Warn.Printf("Failed to delete stats for style %d: %s\n", s.ID, err.Error())
		c.Status(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	if err = search.DeleteStyle(s.ID); err != nil {
		log.Warn.Printf("Failed to delete style %d from index: %s", s.ID, err)
	}

	go func(baseURL string, style *models.APIStyle, modLogID uint) {
		// Add notification to database.
		notification := models.Notification{
			Seen:     false,
			Kind:     models.KindBannedStyle,
			TargetID: int(modLogID),
			UserID:   int(u.ID),
			StyleID:  int(s.ID),
		}

		if err := notification.Create(); err != nil {
			log.Warn.Printf("Failed to create a notification for ban removal %d: %v\n", style.ID, err)
		}

		targetUser, err := models.FindUserByID(strconv.Itoa(int(style.UserID)))
		if err != nil {
			log.Warn.Printf("Failed to find user %d: %s", style.UserID, err.Error())
			return
		}

		// Notify the author about style removal.
		if err := sendBanEmail(baseURL, targetUser, style, modLogID); err != nil {
			log.Warn.Printf("Failed to mail author for style %d: %s", style.ID, err.Error())
		}
	}(c.BaseURL(), s, logEntry.ID)

	return c.Redirect("/modlog", fiber.StatusSeeOther)
}
