package user

import (
	"log"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	jwtware "userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func VerifyGet(c *fiber.Ctx) error {
	if u, ok := jwtware.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		c.Redirect("/account", fiber.StatusSeeOther)
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

	token, err := utils.DecodePreparedText(base64Key)
	if err != nil {
		log.Printf("Couldn't decode key due to: %s\n", err.Error())
		return c.Render("err", fiber.Map{
			"Title":  "Verifcation key not found",
			"Error:": "Key was not found",
		})
	}

	claims := token.Claims.(jwt.MapClaims)

	regErr := database.DB.Create(&models.User{
		Username: claims["username"].(string),
		Password: utils.GenerateHashedPassword(claims["password"].(string)),
		Email:    claims["email"].(string),
	})

	if regErr.Error != nil {
		log.Printf("Failed to register %s, error: %s", claims["email"].(string), regErr.Error)

		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Register failed",
			"Error": "Internal server error.",
		})
	}

	return c.Render("verification", fiber.Map{
		"Title":        "Successful verifcation",
		"Verification": "Successful email verification",
		"Reason":       "You've successfully verified your email address",
	})

}