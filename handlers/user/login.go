package user

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/config"
	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func LoginGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok == true {
		log.Printf("User %d has set session, redirecting.", u.ID)
		c.Redirect("/account", fiber.StatusSeeOther)
	}

	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}

func LoginPost(c *fiber.Ctx) error {
	form := models.User{
		Email:    c.FormValue("email"),
		Password: c.FormValue("password"),
	}

	if err := utils.Validate().StructExcept(form, "Username", "Biography"); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Println("Validation errors:", errors)

		return c.Render("login", fiber.Map{
			"Title":  "Login failed",
			"Errors": errors,
		})
	}

	user, err := models.FindUserByEmail(database.DB, form.Email)
	if err != nil {
		log.Printf("Failed to find %s, error: %s", form.Email, err)

		c.SendStatus(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login failed",
			"Error": "Invalid credentials.",
		})
	}

	match := utils.CompareHashedPassword(user.Password, form.Password)
	if match != nil {
		log.Printf("Failed to match hash for user: %#+v\n", user.Email)

		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("login", fiber.Map{
			"Title": "Login failed",
			"Error": "Invalid credentials.",
		})
	}

	t, err := utils.NewJWTToken().
		SetClaim("id", user.ID).
		SetClaim("name", user.Username).
		SetClaim("email", user.Email).
		SetExpiration(time.Hour * 24 * 14).
		GetSignedString()

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(&fiber.Cookie{
		Name:     fiber.HeaderAuthorization,
		Value:    t,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 14),
		Secure:   config.DB != "dev.db",
		HTTPOnly: true,
		SameSite: "strict",
	})

	return c.Redirect("/account", fiber.StatusSeeOther)
}
