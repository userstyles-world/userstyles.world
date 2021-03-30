package user

import (
	"encoding/base64"
	"log"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func RegisterGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		c.Redirect("/account", fiber.StatusSeeOther)
	}

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

	jwt, err := utils.NewJWTToken().
		SetClaim("username", u.Username).
		SetClaim("password", utils.GenerateHashedPassword(u.Password)).
		SetClaim("email", u.Email).
		SetExpiration(time.Hour * 2).
		GetSignedString()

	if err != nil {
		log.Fatal(err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	sealedText := base64.StdEncoding.EncodeToString(utils.SealText(jwt))

	link := c.BaseURL() + "/verify/" + url.PathEscape(sealedText)

	PlainPart := utils.NewPart().
		SetBody("Verify this Email-address for your UserStyles World account by clicking the link below.\n\n" +
			link + "\n\n" +
			"If you didn't request to verify an UserStyles World account, you can safely ignore this Email.")
	HTMLPart := utils.NewPart().
		SetBody("<p>Verify this Email-address for your UserStyles World account by clicking the link below.</p>\n" +
			"<br>\n" +
			"<a target=\"_blank\" clicktracking=\"off\" href=\"" + link + "\">Verifcation link</a>\n" +
			"<br><br>\n" +
			"<p>If you didn't request to verify an UserStyles World account, you can safely ignore this Email.</p>").
		SetContentType("text/html")

	emailErr := utils.NewEmail().
		SetTo(u.Email).
		SetSubject("Verify your email address").
		AddPart(*PlainPart).
		AddPart(*HTMLPart).
		SendEmail()

	if emailErr != nil {
		log.Fatalf("Couldn't send email due to %s", err)
	}

	return c.Render("await_verifcation", fiber.Map{
		"Title": "Email Verifcation",
	})
}
