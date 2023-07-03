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
	"userstyles.world/utils"
)

func RegisterGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Info.Printf("User %d has set session, redirecting.\n", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	return c.Render("user/register", fiber.Map{
		"Title":     "Register",
		"Canonical": "register",
	})
}

func RegisterPost(c *fiber.Ctx) error {
	password, confirm := c.FormValue("password"), c.FormValue("confirm")
	if confirm != password {
		return c.Status(fiber.StatusForbidden).
			Render("user/register", fiber.Map{
				"Title": "Register failed",
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
				"Title": "Register failed",
				"Error": "Failed to register. Make sure your input is correct.",
			})
	}

	token, err := utils.NewJWTToken().
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

	link := c.BaseURL() + "/verify/" + util.EncryptText(token, util.AEADCrypto, config.ScrambleConfig)
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
