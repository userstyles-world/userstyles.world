package user

import (
	"errors"
	"time"

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

func RecoverGet(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Info.Printf("User %d has set session, redirecting.\n", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	return c.Render("user/recover", fiber.Map{
		"Title":     "Reset password",
		"Canonical": "recover",
	})
}

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
		partPlain := utils.NewPart().
			SetBody("Hi " + user.Username + ",\n" +
				"We'd like to notice you about a recent action of your account:\n\n" +
				"We've authorized a password change. Which was verified by email.\n" +
				"If you don't recognize this account, please email us at feedback@userstyles.world.\n\n" +
				"Regards,\n" + "The UserStyles.world team")
		partHTML := utils.NewPart().
			SetBody("<p>Hi " + user.Username + ",</p>\n" +
				"<p>We'd like to notice you about a recent action of your account:</p>\n" +
				"<br>\n" +
				"<p>We've authorized a password change. Which was verified by email.</p>\n" +
				"<br>\n" +
				"<p>If you don't recognize this account, " +
				"please email us at <a href=\"mailto:feedback@userstyles.world\">feedback@userstyles.world</a>.</p>\n" +
				"<br>\n" +
				"<p>Regards,</p>\n" + "<p>The UserStyles.world team</p>").
			SetContentType("text/html")

		err := utils.NewEmail().
			SetTo(user.Email).
			SetSubject("Your password has been changed").
			AddPart(*partPlain).
			AddPart(*partHTML).
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

// Only allow a email request to happen every 5 minutes.
const emailRequestLimit = 5 * time.Minute

func RecoverPost(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Info.Printf("User %d has set session, redirecting.\n", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	u := models.User{
		Email: c.FormValue("email"),
	}

	if err := utils.Validate().StructPartial(u, "email"); err != nil {
		var validationError validator.ValidationErrors
		if ok := errors.As(err, &validationError); ok {
			log.Warn.Println("Validation errors:", validationError)
		}

		return c.Status(fiber.StatusInternalServerError).
			Render("user/recover", fiber.Map{
				"Title": "Reset failed",
				"Error": "Failed to send email. Make sure your input is correct.",
			})
	}

	go func(u *models.User) {
		user, err := models.FindUserByEmail(u.Email)
		// Return early if we got a error, or when the LastPasswordReset isn't zero
		// And LastPasswordReset + 5 minutes is later than time.Now(). So we only
		// allow to request a new password token every 5 minutes, also to prevent
		// spamming a user's mail.
		if err != nil || (!user.LastPasswordReset.IsZero() && user.LastPasswordReset.Add(emailRequestLimit).After(time.Now())) {
			return
		}

		if err := user.UpdateLastPasswordRequest(); err != nil {
			log.Warn.Printf("Not able to update user's last password reset: %v\n", err)
			return
		}

		jwtToken, err := utils.NewJWTToken().
			SetClaim("email", u.Email).
			SetExpiration(time.Now().Add(time.Hour * 2)).
			GetSignedString(utils.VerifySigningKey)
		if err != nil {
			return
		}

		link := c.BaseURL() + "/reset/" + utils.EncryptText(jwtToken, utils.AEADCrypto, config.ScrambleConfig)

		partPlain := utils.NewPart().
			SetBody("Hi " + user.Username + ",\n" +
				"We have received a request to reset the password for your UserStyles.world account.\n\n" +
				"Follow the link bellow to reset your password. The link will expire in 4 hours.\n" +
				link + "\n\n" +
				"You can safely ignore this email if you didn't request to reset your password.")
		partHTML := utils.NewPart().
			SetBody("<p>Hi " + user.Username + ",</p>\n" +
				"<p>We have received a request to reset the password for your UserStyles.world account.</p>\n" +
				"<br>\n" +
				"<p>Click the link bellow to reset your password. <b>The link will expire in 4 hours.</b><br>\n" +
				"<a target=\"_blank\" clicktracking=\"off\" href=\"" + link + "\">Reset your password</a></p>\n" +
				"<br>\n" +
				"<p>You can safely ignore this email if you didn't request to reset your password.</p>").
			SetContentType("text/html")

		emailErr := utils.NewEmail().
			SetTo(u.Email).
			SetSubject("Reset your password").
			AddPart(*partPlain).
			AddPart(*partHTML).
			SendEmail(config.IMAPServer)
	}

	// We need to just say we have send an reset email.
	// So that we can't leak if we have such email in our database ;).

	return c.Render("user/email-sent", fiber.Map{
		"Title":  "Password reset",
		"Reason": "We've sent an email to reset your password.",
	})
}
