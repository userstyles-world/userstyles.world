package user

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	val "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
	"userstyles.world/modules/validator"
)

func Account(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	user, err := models.FindUserByName(u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	return c.Render("user/account", fiber.Map{
		"Title":  "Account",
		"User":   u,
		"Params": user,
	})
}

func EditAccount(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	user, err := models.FindUserByName(u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	record := map[string]any{"id": user.ID}

	form := c.Params("form")
	switch form {
	case "name":
		name := strings.TrimSpace(c.FormValue("name"))
		prev := user.DisplayName
		user.DisplayName = name

		if name == "" {
			log.Info.Printf("User %v cleared their name\n", u.ID)
			record["display_name"] = name
			break
		}

		if err := validator.V.StructPartial(user, "DisplayName"); err != nil {
			var validationError val.ValidationErrors
			if ok := errors.As(err, &validationError); ok {
				log.Info.Printf("Validation errors for user %d: %v\n", u.ID, validationError)
			}
			user.DisplayName = prev

			l := len(name)
			var e string
			switch {
			case l < 3 || l > 32:
				e = "Display name must be between 3 and 32 characters."
			default:
				e = "Make sure your input contains valid characters."
			}

			return c.Render("user/account", fiber.Map{
				"Title":  "Validation Error",
				"User":   u,
				"Params": user,
				"Error":  e,
			})
		}

		record["display_name"] = name

	case "password":
		current := c.FormValue("current")
		if user.Password != "" {
			err := util.VerifyPassword(user.Password, current)
			if err != nil {
				return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
					"Title": "Failed to match current password",
					"User":  u,
				})
			}
		}

		newPassword, confirmPassword := c.FormValue("new_password"), c.FormValue("confirm_password")
		if confirmPassword != newPassword {
			return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
				"Title": "Failed to match new passwords",
				"User":  u,
			})
		}

		user.Password = newPassword
		if err := validator.V.StructPartial(user, "Password"); err != nil {
			var validationError val.ValidationErrors
			if ok := errors.As(err, &validationError); ok {
				log.Info.Println("Password change error:", validationError)
			}
			return c.Status(fiber.StatusForbidden).Render("err", fiber.Map{
				"Title": "Failed to validate inputs",
				"User":  u,
			})
		}

		pw, err := util.HashPassword(newPassword, config.Secrets)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
				"Title": "Failed to hash password",
				"User":  u,
			})
		}
		user.Password = pw
		record["password"] = user.Password

	case "biography":
		user.Biography = strings.TrimSpace(c.FormValue("biography"))

		if err := validator.V.StructPartial(user, "Biography"); err != nil {
			var validationError val.ValidationErrors
			if ok := errors.As(err, &validationError); ok {
				log.Info.Println("Validation errors:", validationError)
			}

			return c.Render("user/account", fiber.Map{
				"Title":  "Validation Error",
				"User":   u,
				"Params": user,
				"errBio": "Biography is too long. It must be up to 1000 characters.",
			})
		}

		record["biography"] = user.Biography

	case "socials":
		record["github"] = strings.TrimSpace(c.FormValue("github"))
		record["gitlab"] = strings.TrimSpace(c.FormValue("gitlab"))
		record["codeberg"] = strings.TrimSpace(c.FormValue("codeberg"))

	case "flags":
		b, err := json.Marshal(models.Flags{
			Welcome:         c.FormValue("welcomeFlag") == "on",
			Sidebar:         c.FormValue("sidebarFlag") == "on",
			SearchAutofocus: c.FormValue("autofocusFlag") == "on",
			ViewRedesign:    c.FormValue("viewRedesignFlag") == "on",
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
				"Title": "Failed to create 'flags' cookie",
				"User":  u,
			})
		}

		v := string(b)
		log.Info.Printf("kind=flags username=%s data=%s", u.Username, v)

		cookie := &fiber.Cookie{
			Name:     "flags",
			Value:    v,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour * 24 * 30 * 6),
			Secure:   config.App.Production,
			HTTPOnly: true,
			SameSite: fiber.CookieSameSiteLaxMode,
		}
		c.Cookie(cookie)

	default:
		return c.Render("err", fiber.Map{
			"Title": "Invalid form",
			"User":  u,
		})
	}

	dbErr := database.Conn.
		Model(models.User{}).
		Where("id", record["id"]).
		Updates(record).
		Error

	if dbErr != nil {
		log.Warn.Println("Updating user profile failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"User":  u,
		})
	}

	return c.Status(fiber.StatusSeeOther).Redirect("/account#" + form)
}
