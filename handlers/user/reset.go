package user

import (
	"log"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func hasAccount(c *fiber.Ctx) {
	if u, ok := jwtware.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		c.Redirect("/account", fiber.StatusSeeOther)
	}
}

func RecoverGet(c *fiber.Ctx) error {
	hasAccount(c)
	return c.Render("reset", fiber.Map{
		"Title": "Reset",
	})
}

func ResetGet(c *fiber.Ctx) error {
	hasAccount(c)
	renderError := c.Render("err", fiber.Map{
		"Title": "Reset key not found",
	})

	key := c.Params("key")
	if key == "" {
		return renderError
	}

	_, err := utils.DecodePreparedText(key)
	if err != nil {
		log.Printf("Couldn't decode key due to: %s\n", err.Error())
		return renderError
	}

	return c.Render("reset_password", fiber.Map{
		"Title": "Reset password",
		"Key":   key,
	})
}

// Todo: Send email that password has been changed.
func ResetPost(c *fiber.Ctx) error {
	hasAccount(c)

	renderError := c.Render("err", fiber.Map{
		"Title":  "Reset key not found",
		"Error:": "Key was not found",
	})

	password, key := c.FormValue("password"), c.Params("key")

	// Using unified Errors, won't give possible attackers any information.
	// If the process went good.
	if key == "" {
		return renderError
	}

	token, err := utils.DecodePreparedText(key)
	if err != nil {
		log.Printf("Couldn't decode key due to: %s\n", err.Error())
		return renderError
	}

	claims := token.Claims.(jwt.MapClaims)

	user, err := models.FindUserByEmail(database.DB, claims["email"].(string))
	if err != nil {
		return renderError
	}

	t := new(models.User)
	user.Password = utils.GenerateHashedPassword(password)

	err = database.DB.
		Model(t).
		Where("id", user.ID).
		Updates(user).
		Error

	if err != nil {
		log.Println("Updating user failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"Error": "Internal server error.",
		})
	}

	return c.Render("verification", fiber.Map{
		"Title":        "Successful reset",
		"Verification": "Successful verification reset",
		"Reason":       "Your successfully resseted your password",
	})
}

func RecoverPost(c *fiber.Ctx) error {
	hasAccount(c)

	u := models.User{
		Email: c.FormValue("email"),
	}

	if err := utils.Validate().StructPartial(u, "email"); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Println("Validation errors:", errors)

		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("reset", fiber.Map{
			"Title": "Reset failed",
			"Error": "Failed to send email. Make sure your input is correct.",
		})
	}

	if _, err := models.FindUserByEmail(database.DB, u.Email); err != nil {
		// We need to just say we have send an reset email.
		// So that we can't leak if we have such email in our database ;).
		return c.Render("send_email", fiber.Map{
			"Title":  "Password reset",
			"Reason": "We've sent an email to reset your password.",
		})
	}

	jwt, err := utils.NewJWTToken().
		SetClaim("email", u.Email).
		SetExpiration(time.Hour * 2).
		GetSignedString(utils.VerifySigningKey)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
		})
	}

	link := c.BaseURL() + "/reset/" + utils.PrepareText(jwt)

	PlainPart := utils.NewPart().
		SetBody("We have received a request to reset the password for your UserStyles.world account.\n\n" +
			"The link will expire in 2 hours\n\n" +
			link + "\n\n" +
			"You can safely ignore this e-mail if you didn't request to reset your password.")
	HTMLPart := utils.NewPart().
		SetBody("<p>We have received a request to reset the password for your UserStyles.world account.</p>\n" +
			"<b>The link will expire in 2 hours</b>\n" +
			"<br>\n" +
			"<a target=\"_blank\" clicktracking=\"off\" href=\"" + link + "\">Reset password link</a>\n" +
			"<br><br>\n" +
			"<p>You can safely ignore this e-mail if you didn't request to reset your password.</p>").
		SetContentType("text/html")

	emailErr := utils.NewEmail().
		SetTo(u.Email).
		SetSubject("Reset your password").
		AddPart(*PlainPart).
		AddPart(*HTMLPart).
		SendEmail()

	if emailErr != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
		})
	}

	return c.Render("send_email", fiber.Map{
		"Title":  "Password reset",
		"Reason": "We've sent an email to reset your password.",
	})
}
