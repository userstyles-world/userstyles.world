package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
	"userstyles.world/modules/util"
	"userstyles.world/utils"
)

func VerifyGet(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Info.Printf("User %d has set session, redirecting.\n", u.ID)
		return c.Redirect("/account", fiber.StatusSeeOther)
	}

	// Using unified Errors, won't give possible attackers any information.
	// If the process went good.
	base64Key := c.Params("key")
	if base64Key == "" {
		return c.Render("err", fiber.Map{
			"Title":  "Verifcation key not found",
			"Error:": "Key was not found",
		})
	}

	unSealedText, err := utils.DecryptText(base64Key, utils.AEADCrypto, config.ScrambleConfig)
	if err != nil {
		log.Warn.Printf("Failed to decode JWT text: %s\n", err.Error())
		return c.Render("err", fiber.Map{
			"Title":  "Verifcation key not found",
			"Error:": "Key was not found",
		})
	}

	token, err := jwt.Parse(unSealedText, utils.VerifyJwtKeyFunction)
	if err != nil || !token.Valid {
		log.Warn.Printf("Failed to decode JWT token: %s\n", err.Error())
		return c.Render("err", fiber.Map{
			"Title":  "Verifcation key not found",
			"Error:": "Key was not found",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Register failed",
				"Error": "Internal server error.",
			})
	}
	u := &models.User{
		Username: claims["username"].(string),
		Password: claims["password"].(string),
		Email:    claims["email"].(string),
	}

	pw, err := util.HashPassword(u.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("err", fiber.Map{
			"Title": "Failed to hash password",
		})
	}
	u.Password = pw

	if err = database.Conn.Create(u).Error; err != nil {
		log.Database.Printf("Failed to register %s: %s\n", u.Email, err)

		return c.Status(fiber.StatusInternalServerError).
			Render("err", fiber.Map{
				"Title": "Register failed",
				"Error": "Internal server error.",
			})
	}

	return c.Render("user/verification", fiber.Map{
		"Title":        "Successful verifcation",
		"Verification": "Successful email verification",
		"Reason":       "You've successfully verified your email address",
	})
}
