package user

import (
	"log"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/utils"
)

func LoginGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}
	arguments := fiber.Map{
		"Title": "Login",
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
		errors := err.(validator.ValidationErrors)
		log.Println("Validation errors:", errors)

		return c.Render("user/login", fiber.Map{
			"Title":  "Login failed",
			"Errors": errors,
		})
	}

	user, err := models.FindUserByEmail(form.Email)
	if err != nil {
		log.Printf("Failed to find %s, error: %s", form.Email, err)

		return c.Status(fiber.StatusUnauthorized).
			Render("user/login", fiber.Map{
				"Title": "Login failed",
				"Error": "Invalid credentials.",
			})
	}
	if user.OAuthProvider != "none" {
		return c.Status(fiber.StatusUnauthorized).
			Render("user/login", fiber.Map{
				"Title": "Login failed",
				"Error": "Login via OAuth provider",
			})
	}

	match := utils.CompareHashedPassword(user.Password, form.Password)
	if match != nil {
		log.Printf("Failed to match hash for user: %#+v\n", user.Email)

		return c.Status(fiber.StatusInternalServerError).
			Render("user/login", fiber.Map{
				"Title": "Login failed",
				"Error": "Invalid credentials.",
			})
	}

	var expiration time.Time
	if remember {
		// 1 months
		expiration = time.Now().Add(time.Hour * 24 * 31)
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
				"Title": "Internal server error.",
			})
	}

	c.Cookie(&fiber.Cookie{
		Name:     fiber.HeaderAuthorization,
		Value:    t,
		Path:     "/",
		Expires:  expiration,
		Secure:   config.Production,
		HTTPOnly: true,
		SameSite: "strict",
	})

	if r := c.Query("r"); r != "" {
		path, err := url.QueryUnescape(r)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				Render("err", fiber.Map{
					"Title": "Internal server error.",
				})
		}
		return c.Redirect(path, fiber.StatusSeeOther)
	}

	return c.Redirect("/account", fiber.StatusSeeOther)
}
