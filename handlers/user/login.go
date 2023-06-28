package user

import (
	"errors"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
	"userstyles.world/utils"
)

func LoginGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Warn.Printf("User %d has set session, redirecting.", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}
	arguments := fiber.Map{
		"Title":     "Login",
		"Canonical": "login",
	}

	if r := c.Query("r"); r != "" {
		arguments["Redirect"] = "?r=" + url.QueryEscape(r)
		arguments["Error"] = "You must log in to do this action."
	}

	return c.Render("user/login", arguments)
}

func LoginPost(c *fiber.Ctx) error {
	form := models.User{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}
	remember := c.FormValue("remember") == "on"

	err := utils.Validate().StructPartial(form, "Email", "Password")
	if err != nil {
		var validationError validator.ValidationErrors
		if ok := errors.As(err, &validationError); ok {
			log.Warn.Println("Validation errors:", validationError)
		}

		return c.Render("user/login", fiber.Map{
			"Title": "Login failed",
			"Error": "Failed to login. Please provide valid data.",
		})
	}

	user, err := models.FindUserByEmail(form.Email)
	if err != nil {
		log.Warn.Printf("Failed to find %s: %s", form.Email, err.Error())

		return c.Status(fiber.StatusUnauthorized).
			Render("user/login", fiber.Map{
				"Title": "Login failed",
				"Error": "Invalid credentials.",
			})
	}

	match := util.CompareHashedPassword(user.Password, form.Password)
	if match != nil {
		log.Warn.Println("Failed to match hash for user:", user.Email)

		return c.Status(fiber.StatusInternalServerError).
			Render("user/login", fiber.Map{
				"Title": "Login failed",
				"Error": "Invalid credentials.",
			})
	}

	var expiration time.Time
	if remember {
		// 3 months
		expiration = time.Now().Add(time.Hour * 24 * 31 * 3)
	}
	t, err := utils.NewJWTToken().
		SetClaim("id", user.ID).
		SetClaim("name", user.Username).
		SetClaim("email", user.Email).
		SetClaim("role", user.Role).
		SetExpiration(expiration).
		GetSignedString(nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Internal server error",
			})
	}

	// Create and set cookie.
	cookie := &fiber.Cookie{
		Name:     fiber.HeaderAuthorization,
		Value:    t,
		Path:     "/",
		Expires:  expiration,
		Secure:   config.Production,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteLaxMode,
	}
	c.Cookie(cookie)

	if r := c.Query("r"); r != "" {
		path, err := url.QueryUnescape(r)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				Render("err", fiber.Map{
					"Title": "Internal server error",
				})
		}
		return c.Redirect(path, fiber.StatusSeeOther)
	}

	// Update last login date.
	if err = user.UpdateLastLogin(); err != nil {
		log.Warn.Println("Failed to update last login date for id:", user.ID)
	}

	args := fiber.Map{"User": user}
	go email.Send("user/login", user.Email, "New sign-in", args)

	return c.Redirect("/account", fiber.StatusSeeOther)
}
