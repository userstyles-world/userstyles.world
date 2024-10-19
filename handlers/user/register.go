package user

import (
	"errors"
	"time"

	val "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
	"userstyles.world/modules/validator"
)

func RegisterGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Info.Printf("User %d has set session, redirecting.\n", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	return c.Render("user/register", fiber.Map{
		"Title":     "Sign up",
		"Canonical": "signup",
	})
}

func RegisterPost(c *fiber.Ctx) error {
	bot := c.FormValue("bot") != "on"
	if bot {
		username := c.FormValue("username")
		password := c.FormValue("password")
		email := c.FormValue("email")
		log.Spam.Printf("kind=botRegisterPost ip=%q un=%q email=%q password=%d ua=%q",
			c.IP(), username, email, len(password), c.Context().UserAgent())
		return c.Render("err", fiber.Map{"Title": "Bots aren't allowed"})
	}

	password, confirm := c.FormValue("password"), c.FormValue("confirm")
	if confirm != password {
		return c.Status(fiber.StatusForbidden).
			Render("user/register", fiber.Map{
				"Title": "Sign up failed",
				"Error": "Your passwords don't match.",
			})
	}

	u := models.User{
		Username: c.FormValue("username"),
		Password: password,
		Email:    c.FormValue("email"),
	}

	err := validator.V.StructPartial(u, "Username", "Email", "Password")
	if err != nil {
		var validationError val.ValidationErrors
		if ok := errors.As(err, &validationError); ok {
			log.Info.Println("Validation errors:", validationError)
		}
		return c.Status(fiber.StatusInternalServerError).
			Render("user/register", fiber.Map{
				"Title": "Sign up failed",
				"Error": "Failed to sign up. Make sure your input is correct.",
			})
	}

	token, err := util.NewJWT().
		SetClaim("username", u.Username).
		SetClaim("password", u.Password).
		SetClaim("email", u.Email).
		SetExpiration(time.Now().Add(time.Hour * 4)).
		GetSignedString(util.VerifySigningKey)
	if err != nil {
		log.Warn.Println("Failed to create a JWT Token:", err.Error())
		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Internal server error",
			})
	}

	link := c.BaseURL() + "/verify/" + util.EncryptText(token, util.AEADCrypto, config.Secrets)
	args := fiber.Map{
		"User": u,
		"Link": link,
	}

	err = email.Send("user/register", u.Email, "Verify your email address", args)
	if err != nil {
		log.Warn.Printf("Failed to send an email: %s\n", err)
		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Internal server error",
				"Error": "Failed to send email.",
			})
	}

	return c.Render("user/email-sent", fiber.Map{
		"Title":  "Email verification sent",
		"Reason": "Verification link has been sent to your email address. Your account will be created when you visit the verification link.",
	})
}
