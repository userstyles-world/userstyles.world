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

	i, err := c.ParamsInt("id")
	if err != nil || i < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"Title": "Invalid style ID",
			"User":  u,
		})
	}

	// Check if style exists.
	s, err := models.GetStyleByID(i)
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

func BanStyle(db *gorm.DB, style *models.Style, u *models.APIUser, user *storage.User, c *fiber.Ctx) (*models.Log, error) {
	event := &models.Log{
		UserID:         u.ID,
		Username:       u.Username,
		Kind:           models.LogRemoveStyle,
		TargetUserName: user.Username,
		TargetData:     style.Name,
		Reason:         strings.TrimSpace(c.FormValue("reason")),
		Message:        strings.TrimSpace(c.FormValue("message")),
		Censor:         c.FormValue("censor") == "on",
	}

	n := &models.Notification{
		Kind:     models.KindBannedStyle,
		TargetID: int(event.ID),
		UserID:   int(user.ID),
		StyleID:  int(style.ID),
	}

	i := int(style.ID)
	if err := storage.DeleteUserstyle(db, i); err != nil {
		return nil, err
	}
	if err := models.DeleteReviewsForStyle(db, i); err != nil {
		return nil, err
	}
	if err := models.DeleteNotificationsForStyle(db, i); err != nil {
		return nil, err
	}
	if err := models.DeleteStats(db, i); err != nil {
		return nil, err
	}
	if err := storage.DeleteSearchData(db, i); err != nil {
		return nil, err
	}
	if err := models.CreateLog(db, event); err != nil {
		return nil, err
	}
	if err := models.CreateNotification(db, n); err != nil {
		return nil, err
	}
	if err := models.RemoveStyleCode(strconv.Itoa(i)); err != nil {
		return nil, err
	}

	cache.Code.Remove(i)

	return event, nil
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

	style, err := models.GetStyleByID(i)
	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.Render("err", fiber.Map{
			"Title": "Style not found",
			"User":  u,
		})
	}

	user, err := storage.FindUser(style.UserID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	var event *models.Log
	err = database.Conn.Transaction(func(tx *gorm.DB) error {
		event, err = BanStyle(tx, style, u, user, c)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Database.Printf("Failed to remove %d: %s\n", i, err)
		c.Status(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Failed to remove userstyle",
			"User":  u,
		})
	}

	go sendRemovalEmail(user, style, event)

	return c.Redirect(event.Permalink())
}

func sendRemovalEmail(user *storage.User, style *models.Style, event *models.Log) {
	args := fiber.Map{
		"User":  user,
		"Style": style,
		"Log":   event,
		"Link":  config.App.BaseURL + event.Permalink(),
	}

	title := "Your style has been removed"
	if err := email.Send("style/ban", user.Email, title, args); err != nil {
		log.Warn.Printf("Failed to email %d: %s\n", user.ID, err)
	}
}
