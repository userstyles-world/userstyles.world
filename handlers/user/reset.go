package user

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/email"
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

	args := fiber.Map{"User": user}
	title := "Your password has been changed"
	go email.Send("passwordreset", user.Email, title, args)

	return c.Render("user/verification", fiber.Map{
		"Title":        "Successful reset",
		"Verification": "Successful password reset",
		"Reason":       "You've successfully changed your password",
	})
}
