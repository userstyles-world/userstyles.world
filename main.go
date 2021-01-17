package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"userstyles.world/config"
	"userstyles.world/database"
	"userstyles.world/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html"
)

func main() {
	database.Connect()
	database.Prepare()

	var store = session.New()

	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "UserStyles.world",
			"Body":  "Hello, World!",
		})
	})

	app.Get("/register", func(c *fiber.Ctx) error {
		return c.Render("register", fiber.Map{
			"Title": "UserStyles.world",
			"Body":  "Hello, World!",
		})
	})

	app.Post("/register", func(c *fiber.Ctx) error {
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
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"Title": "UserStyles.world",
			"Body":  "Hello, World!",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		form := models.User{
			Email:    c.FormValue("email"),
			Password: c.FormValue("password"),
		}

		user, err := models.FindUserByEmail(database.DB, form.Email)
		if err != nil {
			log.Printf("Failed to find %s, error: %s", form.Email, err)

			c.SendStatus(fiber.StatusSeeOther)
			return c.Render("login", fiber.Map{
				"Error": "Invalid credentials.",
			})
		}

		match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		if match != nil {
			log.Println("Failed to match hash for user:", user.Email)
			c.SendStatus(fiber.StatusUnauthorized)
			return c.Render("login", fiber.Map{
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
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			log.Println(err)
		}
		sess.Destroy()

		return c.Redirect("/login", fiber.StatusFound)
	})

	app.Get("/account", func(c *fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			log.Println(err)
		}

		return c.Render("account", fiber.Map{
			"Name":  sess.Get("name"),
			"Title": "UserStyles.world",
			"Body":  "Hello, World!",
		})
	})

	log.Fatal(app.Listen(config.PORT))
}
