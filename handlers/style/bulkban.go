package style

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

type bulkReq struct {
	IDs []string
}

func BulkBanGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	if !u.IsModOrAdmin() {
		c.Locals("Title", "You are not authorized to perform this action")
		return c.Status(fiber.StatusUnauthorized).Render("err", fiber.Map{})
	}

	id, err := c.ParamsInt("userid")
	if err != nil || id < 1 {
		c.Locals("Title", "Invalid user ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	if _, err = storage.FindUser(uint(id)); err != nil {
		c.Locals("Title", "Could not find such user")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	var styles []models.APIStyle
	err = database.Conn.Find(&styles, "user_id = ? AND deleted_at IS NULL", id).Error
	if err != nil || len(styles) == 0 {
		c.Locals("Title", "Could not find any userstyles")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}
	c.Locals("Styles", styles)

	c.Locals("UserID", id)
	c.Locals("Title", "Perform a bulk userstyle removal")

	return c.Render("style/bulkban", fiber.Map{})
}

func BulkBanPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)
	c.Locals("User", u)

	if !u.IsModOrAdmin() {
		c.Locals("Title", "You are not authorized to perform this action")
		return c.Status(fiber.StatusUnauthorized).Render("err", fiber.Map{})
	}

	uid, err := c.ParamsInt("userid")
	if err != nil || uid < 1 {
		c.Locals("Title", "Invalid user ID")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	user, err := storage.FindUser(uint(uid))
	if err != nil {
		c.Locals("Title", "Could not find such user")
		return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
	}

	var req bulkReq
	if err = c.BodyParser(&req); err != nil {
		c.Locals("Title", "Failed to process request body")
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
	}

	var styles []*models.Style

	// Process all IDs for problems not to have any errors in between of removal
	for _, val := range req.IDs {
		id, err := strconv.Atoi(val)
		if err != nil {
			c.Locals("Title", "Operation failed")
			c.Locals("ErrTitle", val+" is not a valid number")
			return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{})
		}

		style, err := models.GetStyleFromAuthor(id, uid)
		if err != nil {
			c.Locals("Title", "Operation failed")
			c.Locals("ErrTitle", "User isn't the author of style with ID "+val)
			return c.Status(fiber.StatusNotFound).Render("err", fiber.Map{})
		}

		styles = append(styles, &style)
	}

	// lastEvent is used to link to the newest event in the modlog
	// so the user will be presented with all of them on the screen.
	var lastEvent *models.Log
	err = database.Conn.Transaction(func(tx *gorm.DB) error {
		for index, style := range styles {
			event, err := BanStyle(tx, style, u, user, c)
			if err != nil {
				return err
			}

			if index == len(styles)-1 {
				lastEvent = event
			}
		}

		return nil
	})
	if err != nil {
		c.Locals("Title", "Failed to ban styles")
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{})
	}

	go sendBulkRemovalEmail(user, styles, lastEvent)

	return c.Redirect("/modlog", fiber.StatusSeeOther)
}

func sendBulkRemovalEmail(user *storage.User, styles []*models.Style, event *models.Log) {
	args := fiber.Map{
		"User":   user,
		"Styles": styles,
		"Log":    event,
		"Link":   config.BaseURL + "/modlog#id-" + strconv.Itoa(int(event.ID)),
	}

	title := strconv.Itoa(len(styles)) + " of your style have been removed"
	if err := email.Send("style/bulkban", user.Email, title, args); err != nil {
		log.Warn.Printf("Failed to email %d: %s\n", user.ID, err)
	}
}
