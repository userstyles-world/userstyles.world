package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"userstyles.world/config"
	"userstyles.world/database"
	"userstyles.world/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
)

func main() {
	database.Connect()
	database.Prepare()

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
			log.Panic(err)
		}

		database.DB.Create(&models.User{
			Username: u.Username,
			Password: string(hash),
			Email:    u.Email,
		})

		log.Printf("%+v\n", u)

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
			log.Fatalf("Error finding user: %v", err)
		}

		match := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		if match != nil {
			log.Println("Failed to match hash for user:", user.Email)
			c.SendStatus(fiber.StatusUnauthorized)
			return c.Render("login", fiber.Map{
				"Error": "Invalid login credentials",
			})
		}

		return c.Redirect("/", fiber.StatusFound)
	})

	log.Fatal(app.Listen(config.PORT))
}
