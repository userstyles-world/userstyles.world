package user

import (
	"bytes"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/email"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
)

// Only allow an email request to happen every 5 minutes.
const emailRequestLimit = 5 * time.Minute

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

	go func(u models.User) {
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
			SetExpiration(time.Now().Add(time.Hour * 4)).
			GetSignedString(utils.VerifySigningKey)
		if err != nil {
			log.Warn.Printf("Not able to generate JWT token: %v\n", err)
			return
		}

		link := config.BaseURL + "/reset/" + utils.EncryptText(jwtToken, utils.AEADCrypto, config.ScrambleConfig)

		args := fiber.Map{
			"Link": link,
		}

		var bufText, bufHTML bytes.Buffer
		err = email.Render(&bufText, &bufHTML, "passwordrecovery", args)
		if err != nil {
			log.Warn.Printf("Failed to render email template: %v\n", err)
			return
		}

		utils.NewEmail().
			SetTo(u.Email).
			SetSubject("Reset your password").
			AddPart(*utils.NewPart().SetBody(bufText.String())).
			AddPart(*utils.NewPart().SetBody(bufHTML.String()).HTML()).
			SendEmail(config.IMAPServer)
	}(u)

	// We need to just say we have send an reset email.
	// So that we can't leak if we have such email in our database ;).

	return c.Render("user/email-sent", fiber.Map{
		"Title":  "Password reset",
		"Reason": "If there is an account associated with this email address, we'll send a password reset link to it.",
	})
}
