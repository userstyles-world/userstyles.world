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

func BanStyle(style models.APIStyle, u *models.APIUser, user *storage.User, i int, id string, c *fiber.Ctx) (models.Log, error) {

	event := models.Log{
		UserID:         u.ID,
		Username:       u.Username,
		Kind:           models.LogRemoveStyle,
		TargetUserName: style.Username,
		TargetData:     style.Name,
		Reason:         strings.TrimSpace(c.FormValue("reason")),
		Message:        strings.TrimSpace(c.FormValue("message")),
		Censor:         c.FormValue("censor") == "on",
	}

	notification := models.Notification{
		Kind:     models.KindBannedStyle,
		TargetID: int(event.ID),
		UserID:   int(user.ID),
		StyleID:  int(style.ID),
	}

	// INSERT INTO `logs`
	err := database.Conn.Transaction(func(tx *gorm.DB) error {
		if err := storage.DeleteUserstyle(tx, i); err != nil {
			return err
		}
		if err := models.DeleteStats(tx, i); err != nil {
			return err
		}
		if err := storage.DeleteSearchData(tx, i); err != nil {
			return err
		}
		if err := models.CreateLog(tx, &event); err != nil {
			return err
		}
		if err := models.CreateNotification(tx, &notification); err != nil {
			return err
		}
		return models.RemoveStyleCode(id)
	})
	if err != nil {
		return event, err
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
	id := c.Params("id")

	style, err := models.GetStyleByID(id)
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

	event, err := BanStyle(*style, u, user, i, id, c)
	if err != nil {
		log.Database.Printf("Failed to remove %d: %s\n", i, err)
		return c.Render("err", fiber.Map{
			"Title": "Failed to remove userstyle",
			"User":  u,
		})
	}

	go sendRemovalEmail(user, style, event)

	return c.Redirect("/modlog", fiber.StatusSeeOther)
}

func sendRemovalEmail(user *storage.User, style *models.APIStyle, entry models.Log) {
	args := fiber.Map{
		"User":  user,
		"Style": style,
		"Log":   entry,
		"Link":  config.BaseURL + "/modlog#id-" + strconv.Itoa(int(entry.ID)),
	}

	title := "Your style has been removed"
	if err := email.Send("style/ban", user.Email, title, args); err != nil {
		log.Warn.Printf("Failed to email %d: %s\n", user.ID, err)
	}
}
