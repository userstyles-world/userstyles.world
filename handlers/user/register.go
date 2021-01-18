package user

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"userstyles.world/database"
	"userstyles.world/models"
)

func RegisterGet(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"Title": "UserStyles.world",
		"Body":  "Hello, World!",
	})
}

func RegisterPost(c *fiber.Ctx) error {
	u := models.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
		Email:    c.FormValue("email"),
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		log.Println(err)
	}

	regErr := database.DB.Create(&models.User{
		Username: u.Username,
		Password: string(hash),
		Email:    u.Email,
	})

	if regErr != nil {
		log.Printf("Failed to register %s, error: %s", u.Email, regErr.Error)

		c.SendStatus(fiber.StatusSeeOther)
		return c.Render("register", fiber.Map{
			"Error": "Failed to register. Make sure your credentials are valid.",
		})
	}

	return c.Redirect("/login", fiber.StatusFound)
}
