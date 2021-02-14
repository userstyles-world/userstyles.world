package user

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func LoginGet(c *fiber.Ctx) error {
	s := sessions.State(c)

	if s.Fresh() == false {
		log.Printf("User %s has set session, redirecting.", s.Get("email"))
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

	if err := utils.Validate().StructExcept(form, "Username"); err != nil {
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

	s := sessions.State(c)
	defer s.Save()

	// Set session data.
	s.Set("id", user.ID)
	s.Set("name", user.Username)
	s.Set("email", user.Email)
	log.Println("Session:", s.Get("name"), s.Get("email"))

	return c.Redirect("/account", fiber.StatusSeeOther)
}
