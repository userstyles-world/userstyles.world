package user

import (
	"log"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/database"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func RecoverGet(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}
	return c.Render("reset", fiber.Map{
		"Title": "Reset",
	})
}

func ResetGet(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	renderError := c.Render("err", fiber.Map{
		"Title": "Reset key not found",
	})

	key := c.Params("key")
	if key == "" {
		return renderError
	}

	_, err := utils.DecodePreparedText(key, utils.AEAD_CRYPTO)
	if err != nil {
		log.Printf("Couldn't decode key due to: %s\n", err.Error())
		return renderError
	}

	return c.Render("reset_password", fiber.Map{
		"Title": "Reset password",
		"Key":   key,
	})
}

// TODO: Send email that password has been changed.
func ResetPost(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

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

	unSealedText, err := utils.DecodePreparedText(key, utils.AEAD_CRYPTO)
	if err != nil {
		log.Printf("Couldn't decode key due to: %s\n", err.Error())
		return renderError
	}
	token, err := jwt.Parse(unSealedText, utils.VerifyJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Printf("Couldn't decode key due to: %s\n", err.Error())
		return renderError
	}

	claims := token.Claims.(jwt.MapClaims)

	user, err := models.FindUserByEmail(claims["email"].(string))
	if err != nil {
		return renderError
	}

	t := new(models.User)
	user.Password = utils.GenerateHashedPassword(password)

	err = database.Conn.
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
		"Verification": "Successful password reset",
		"Reason":       "You've successfully changed your password",
	})
}

func RecoverPost(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	u := models.User{
		Email: c.FormValue("email"),
	}

	if err := utils.Validate().StructPartial(u, "email"); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Println("Validation errors:", errors)

		return c.Status(fiber.StatusInternalServerError).
			Render("reset", fiber.Map{
				"Title": "Reset failed",
				"Error": "Failed to send email. Make sure your input is correct.",
			})
	}

	if _, err := models.FindUserByEmail(u.Email); err != nil {
		// We need to just say we have send an reset email.
		// So that we can't leak if we have such email in our database ;).
		return c.Render("send_email", fiber.Map{
			"Title":  "Password reset",
			"Reason": "We've sent an email to reset your password.",
		})
	}

	jwt, err := utils.NewJWTToken().
		SetClaim("email", u.Email).
		SetExpiration(time.Now().Add(time.Hour * 2)).
		GetSignedString(utils.VerifySigningKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Internal server error",
			})
	}

	link := c.BaseURL() + "/reset/" + utils.PrepareText(jwt, utils.AEAD_CRYPTO)

	partPlain := utils.NewPart().
		SetBody("We have received a request to reset the password for your UserStyles.world account.\n\n" +
			"The link will expire in 2 hours\n\n" +
			link + "\n\n" +
			"You can safely ignore this e-mail if you didn't request to reset your password.")
	partHTML := utils.NewPart().
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
		AddPart(*partPlain).
		AddPart(*partHTML).
		SendEmail()

	if emailErr != nil {
		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Internal server error",
			})
	}

	return c.Render("send_email", fiber.Map{
		"Title":  "Password reset",
		"Reason": "We've sent an email to reset your password.",
	})
}
