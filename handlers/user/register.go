package user

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/config"
	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func RegisterGet(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title": "Register",
	})
}

func RegisterPost(c *fiber.Ctx) error {
	u := models.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
		Email:    c.FormValue("email"),
	}

	if err := utils.Validate().Struct(u); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Println("Validation errors:", errors)

		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("register", fiber.Map{
			"Title": "Register failed",
			"Error": "Failed to register. Make sure you've correct inputs.",
		})
	}

	password := utils.GenerateHashedPassword(u.Password)
	regErr := database.DB.Create(&models.User{
		Username: u.Username,
		Password: password,
		Email:    u.Email,
	})

	if regErr.Error != nil {
		log.Printf("Failed to register %s, error: %s", u.Email, regErr.Error)

		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("register", fiber.Map{
			"Title": "Register failed",
			"Error": "Failed to register. Make sure your credentials are valid.",
		})
	}

	user, err := models.FindUserByName(database.DB, u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
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

	return c.Redirect("/login", fiber.StatusSeeOther)
}
