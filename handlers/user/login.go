package user

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"userstyles.world/database"
	"userstyles.world/handlers/sessions"
	"userstyles.world/models"
)

var (
	store = sessions.GetStore()
)

func LoginGet(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}

	if s.Fresh() == false {
		log.Printf("User %s has set session, redirecting.", s.Get("email"))
		c.Redirect("/account", fiber.StatusFound)
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

	if err := validator.New().StructExcept(form, "Username"); err != nil {
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

		c.SendStatus(fiber.StatusSeeOther)
		return c.Render("login", fiber.Map{
			"Title": "Login failed",
			"Error": "Invalid credentials.",
		})
	}

	match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if match != nil {
		log.Println("Failed to match hash for user:", user.Email)
		c.SendStatus(fiber.StatusUnauthorized)
		return c.Render("login", fiber.Map{
			"Title": "Login failed",
			"Error": "Invalid credentials.",
		})
	}

	sess, err := store.Get(c)
	if err != nil {
		log.Println(err)
	}
	defer sess.Save()

	sess.Set("name", user.Username)
	sess.Set("email", user.Email)

	log.Println("Session:", sess.Get("name"), sess.Get("email"))
	return c.Redirect("/account", fiber.StatusFound)
}
