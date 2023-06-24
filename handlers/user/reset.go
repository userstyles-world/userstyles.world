package user

import (
	"bytes"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
)

func ResetGet(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Info.Printf("User %d has set session, redirecting.\n", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	renderError := c.Render("err", fiber.Map{
		"Title": "Reset key not found",
	})

	key := c.Params("key")
	if key == "" {
		return renderError
	}

	_, err := utils.DecryptText(key, utils.AEADCrypto, config.ScrambleConfig)
	if err != nil {
		log.Warn.Println("Failed to unseal JWT text:", err.Error())
		return renderError
	}

	return c.Render("user/reset-password", fiber.Map{
		"Title": "Reset password",
		"Key":   key,
	})
}

func ResetPost(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Info.Printf("User %d has set session, redirecting.\n", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	renderError := c.Render("err", fiber.Map{
		"Title":  "Reset key not found",
		"Error:": "Key was not found",
	})

	// Using unified Errors, won't give possible attackers any information.
	// If the process went good.
	key := c.Params("key")
	if key == "" {
		return renderError
	}

	newPassword, confirmPassword := c.FormValue("new_password"), c.FormValue("confirm_password")
	if newPassword != confirmPassword {
		return c.Status(fiber.StatusBadRequest).Render("user/reset-password", fiber.Map{
			"Title": "Passwords don't match",
			"Error": "Passwords don't match.",
			"Key":   key,
		})
	}

	unSealedText, err := utils.DecryptText(key, utils.AEADCrypto, config.ScrambleConfig)
	if err != nil {
		log.Warn.Println("Failed to unseal JWT text:", err.Error())
		return renderError
	}

	token, err := jwt.Parse(unSealedText, utils.VerifyJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Warn.Println("Failed to unseal JWT token:", err.Error())
		return renderError
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Warn.Println("Failed to parse JWT claims.")
		return renderError
	}

	user, err := models.FindUserByEmail(claims["email"].(string))
	if err != nil {
		return renderError
	}

	t := new(models.User)
	user.Password = newPassword
	if err := utils.Validate().StructPartial(user, "Password"); err != nil {
		var validationError validator.ValidationErrors
		if ok := errors.As(err, &validationError); ok {
			log.Info.Println("Password change error:", validationError)
		}
		return c.Status(fiber.StatusForbidden).Render("user/reset-password", fiber.Map{
			"Title":  "Failed to validate inputs",
			"Errors": validationError,
			"Key":    key,
		})
	}
	user.Password = utils.GenerateHashedPassword(newPassword)

	err = database.Conn.
		Model(t).
		Where("id", user.ID).
		Updates(user).
		Error

	if err != nil {
		log.Warn.Println("Failed to update user:", err.Error())
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
			"Error": "Internal server error.",
		})
	}

	// Sends email that the password has been changed.
	// But we do it in a separate routine, so we can render the view for the user.
	go func(user *models.User) {

		args := fiber.Map{}

		var bufText bytes.Buffer
		var bufHTML bytes.Buffer
		errText := c.App().Config().Views.Render(&bufText, "email/passwordreset.text", args)
		errHTML := c.App().Config().Views.Render(&bufHTML, "email/passwordreset.html", args)
		if errText != nil || errHTML != nil {
			log.Warn.Printf("Failed to render email template: %v\n", err)
			return
		}

		err := utils.NewEmail().
			SetTo(user.Email).
			SetSubject("Your password has been changed").
			AddPart(*utils.NewPart().SetBody(bufText.String())).
			AddPart(*utils.NewPart().SetBody(bufHTML.String()).HTML()).
			SendEmail(config.IMAPServer)
		if err != nil {
			log.Warn.Println("Failed to send an email:", err.Error())
		}
	}(user)

	return c.Render("user/verification", fiber.Map{
		"Title":        "Successful reset",
		"Verification": "Successful password reset",
		"Reason":       "You've successfully changed your password",
	})
}
