package style

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
	"userstyles.world/modules/storage"
)

func BulkBanGet(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Check if logged-in user has permissions.
	if !u.IsModOrAdmin() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("err", fiber.Map{
			"Title": "You are not authorized to perform this action",
			"User":  u,
		})
	}

	userid, err := c.ParamsInt("userid")
	if err != nil || userid < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"User":  u,
			"Title": "Invalid user ID",
		})
	}

	// fixme: this check is not working
	user, _ := storage.FindUser(uint(userid))
	if user == nil {
		c.Status(fiber.StatusInternalServerError)
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"User":  u,
			"Title": "Could not find such user",
		})
	}

	return c.Render("style/bulkban", fiber.Map{
		"Title":  "Perform a bulk ban",
		"User":   u,
		"UserID": userid,
	})
}

func BulkBanPost(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	// Check if logged-in user has permissions.
	if !u.IsModOrAdmin() {
		c.Status(fiber.StatusUnauthorized)
		return c.Render("err", fiber.Map{
			"Title": "You are not authorized to perform this action",
			"User":  u,
		})
	}

	userid, err := c.ParamsInt("userid")
	if err != nil || userid < 1 {
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"User":  u,
			"Title": "Invalid user ID",
		})
	}

	// fixme: this check is not working
	user, _ := storage.FindUser(uint(userid))
	if user == nil {
		c.Status(fiber.StatusInternalServerError)
		return c.Status(fiber.StatusBadRequest).Render("err", fiber.Map{
			"User":  u,
			"Title": "Could not find such user",
		})
	}

	styles := []models.APIStyle{}
	ids := strings.Split(c.FormValue("ids"), ",")

	// Process all IDs for problems not to have any errors in between of removal
	for _, element := range ids {
		id := strings.TrimSpace(element)

		style, err := models.GetStyleByID(id)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.Render("err", fiber.Map{
				"Title":    "Operation failed",
				"ErrTitle": "Style " + id + " was not found",
				"User":     u,
			})
		}
		if int(style.UserID) != userid {
			c.Status(fiber.StatusNotFound)
			return c.Render("err", fiber.Map{
				"Title":    "Operation failed",
				"ErrTitle": "User " + strconv.Itoa(int(style.UserID)) + " is not the author of style " + id,
				"User":     u,
			})
		}
		styles = append(styles, *style)

		_, err = strconv.Atoi(id)
		if err != nil {
			c.Status(fiber.StatusNotFound)
			return c.Render("err", fiber.Map{
				"Title":    "Operation failed",
				"ErrTitle": id + " is not a string",
				"User":     u,
			})
		}
	}

	// lastevent is used to link to the newest event in the modlog
	// so the user will be presented with all of them on the screen.
	var lastevent models.Log
	for index, style := range styles {
		event, _ := BanStyle(style, u, user, int(style.ID), strconv.Itoa(int(style.ID)), c)
		if index == len(styles)-1 {
			lastevent = event
		}
	}

	go sendBulkRemovalEmail(user, styles, lastevent)

	return c.Redirect("/modlog", fiber.StatusSeeOther)
}

func sendBulkRemovalEmail(user *storage.User, styles []models.APIStyle, firstentry models.Log) {
	args := fiber.Map{
		"User":   user,
		"Styles": styles,
		"Log":    firstentry,
		"Link":   config.BaseURL + "/modlog#id-" + strconv.Itoa(int(firstentry.ID)),
	}

	title := strconv.Itoa(len(styles)) + " of your style have been removed"
	if err := email.Send("style/bulkban", user.Email, title, args); err != nil {
		log.Warn.Printf("Failed to email %d: %s\n", user.ID, err)
	}
}
