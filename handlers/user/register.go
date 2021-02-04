package user

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

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

	if err := validator.New().Struct(u); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Println("Validation errors:", errors)

		return c.Render("register", fiber.Map{
			"Title":  "Register failed",
			"Errors": errors,
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

		c.SendStatus(fiber.StatusSeeOther)
		return c.Render("register", fiber.Map{
			"Title": "Register failed",
			"Error": "Failed to register. Make sure your credentials are valid.",
		})
	}

	return c.Redirect("/login", fiber.StatusFound)
}
