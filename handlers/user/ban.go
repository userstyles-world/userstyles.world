package user

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
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

func sendBanEmail(baseURL string, user *models.User, modLogID uint) error {
	modLogEntry := baseURL + "/modlog#id-" + strconv.Itoa(int(modLogID))

	partPlain := utils.NewPart().
		SetBody("Hi " + user.Username + ",\n" +
			"We'd like to notice you about a recent action from our moderation team:\n\n" +
			"You have been banned from our platform.\n" +
			"You can check for more information about this action on the modlog: " + modLogEntry + "\n\n" +
			"If you'd like to come in touch with us,please email us at feedback@userstyles.world\n" +
			"Regards,\n" + "The Moderation Team")
	partHTML := utils.NewPart().
		SetBody("<p>Hi " + user.Username + ",</p>\n" +
			"<p>We'd like to notice you about a recent action from our moderation team:</p>\n" +
			"<br>\n" +
			"<p>You have been banned from our platform.</p>\n" +
			"<p>You can check for more information about this action on the " +
			"<a target=\"_blank\" clicktracking=\"off\" href=\"" + modLogEntry + "\">Modlog</a>.</p>\n" +
			"<p>If you'd like to come in touch with us, please email us at " +
			"<a href=\"mailto:feedback@userstyles.world\">feedback@userstyles.world</a>.<p>\n" +
			"<br>\n" +
			"<p>Regards,</p>\n" + "<p>The Moderation Team</p>").
		SetContentType("text/html")

	err := utils.NewEmail().
		SetTo(user.Email).
		SetSubject("You have been banned").
		AddPart(*partPlain).
		AddPart(*partHTML).
		SendEmail(config.IMAPServer)
	if err != nil {
		return err
	}
	return nil
}

func ConfirmBan(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	stringID := c.Params("id")
	id, _ := strconv.Atoi(stringID)

	if !u.IsModOrAdmin() {
		return c.Render("err", fiber.Map{
			"Title": "Unauthorized",
			"User":  u,
		})
	}

	if u.ID == uint(id) {
		return c.Render("err", fiber.Map{
			"Title": "You can't ban yourself",
			"User":  u,
		})
	}

	// Check if user exists.
	targetUser, err := models.FindUserByID(stringID)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User ID doesn't exist",
			"User":  u,
		})
	}

	// Delete from database.
	user := new(models.User)
	if err := user.DeleteWhereID(targetUser.ID); err != nil {
		log.Warn.Printf("Failed to ban user %d: %s\n", id, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	// Delete user's styles.
	styles := new(models.Style)
	if err := styles.BanWhereUserID(targetUser.ID); err != nil {
		log.Warn.Printf("Failed to ban styles from user %d: %s\n", id, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	// Initialize modlog data.
	logEntry := models.Log{
		UserID:         u.ID,
		Username:       u.Username,
		Kind:           models.LogBanUser,
		TargetUserName: targetUser.Username,
		Reason:         strings.TrimSpace(c.FormValue("reason")),
		Censor:         c.FormValue("censor") == "on",
	}

	// Add banned user log entry.
	modlog := new(models.Log)
	if err := modlog.AddLog(&logEntry); err != nil {
		log.Warn.Printf("Failed to add user %d to ModLog: %s\n", targetUser.ID, err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	go func(baseURL string, user *models.User, modLogID uint) {
		// Send a email about they've been banned.
		if err := sendBanEmail(baseURL, targetUser, modLogID); err != nil {
			log.Warn.Printf("Failed to send an email to user %d: %s", user.ID, err.Error())
		}
	}(c.BaseURL(), targetUser, logEntry.ID)

	return c.Redirect("/dashboard")
}
