package style

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
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

func BanPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Check if logged-in user has permissions.
	if !u.IsModOrAdmin() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("err", fiber.Map{
			"Title": "You are not authorized to perform this action",
			"User":  u,
		})
	}

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"User":  u,
			"Title": "Invalid style ID",
		})
	}
	id := c.Params("id")

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
		Message:        strings.TrimSpace(c.FormValue("message")),
		Censor:         c.FormValue("censor") == "on",
	}

	err = database.Conn.Transaction(func(tx *gorm.DB) error {
		if err = storage.DeleteUserstyle(tx, i); err != nil {
			return err
		}
		if err = models.DeleteStats(tx, i); err != nil {
			return err
		}
		if err = storage.DeleteSearchData(tx, i); err != nil {
			return err
		}
		if err = models.CreateLog(tx, &logEntry); err != nil {
			return err
		}
		return models.RemoveStyleCode(id)
	})
	if err != nil {
		log.Database.Printf("Failed to remove %d: %s\n", i, err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to remove userstyle",
			"User":  u,
		})
	}

	cache.Code.Remove(i)

	go func(style *models.APIStyle, entry models.Log) {
		user, err := models.FindUserByID(strconv.Itoa(int(style.UserID)))
		if err != nil {
			log.Warn.Printf("Failed to find user %d: %s", style.UserID, err.Error())
			return
		}

		// Add notification to database.
		notification := models.Notification{
			Seen:     false,
			Kind:     models.KindBannedStyle,
			TargetID: int(entry.ID),
			UserID:   int(user.ID),
			StyleID:  int(style.ID),
		}

		if err := notification.Create(); err != nil {
			log.Warn.Printf("Failed to create a notification for ban removal %d: %v\n", style.ID, err)
		}

		args := fiber.Map{
			"User":  user,
			"Style": style,
			"Log":   entry,
			"Link":  config.BaseURL + "/modlog#id-" + strconv.Itoa(int(entry.ID)),
		}

		title := "Your style has been removed"
		if err := email.Send("style/ban", user.Email, title, args); err != nil {
			log.Warn.Printf("Failed to email author for style %d: %s\n", style.ID, err)
		}
	}(s, logEntry)

	return c.Redirect("/modlog", fiber.StatusSeeOther)
}
