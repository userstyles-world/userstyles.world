package user

import (
	"encoding/base64"
	"log"
	"net/url"

	"github.com/form3tech-oss/jwt-go"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func VerifyGet(c *fiber.Ctx) error {
	// Using unified Errors, won't give possible attackers any information.
	// If the process went good.

	base64Key := c.Params("key")
	if base64Key == "" {
		return c.Render("err", fiber.Map{
			"Title":  "Verifcation key not found",
			"Error:": "Key was not found",
		})
	}
	unescapeKey, err := url.PathUnescape(base64Key)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title":  "Verifcation key not found",
			"Error:": "Key was not found",
		})
	}

	key, err := base64.StdEncoding.DecodeString(unescapeKey)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title":  "Verifcation key not found",
			"Error:": "Key was not found",
		})
	}

	decryptedText, err := utils.OpenText(utils.FastByteToString(key))
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title":  "Verifcation key not found",
			"Error:": "Key was not found",
		})
	}

	token, err := jwt.Parse(utils.FastByteToString(decryptedText), utils.VerifyJwtKeyFunction)
	if err != nil || !token.Valid {
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
			"Error": "Failed to register. Make sure your credentials are valid.",
		})
	}

	return c.Render("verification", fiber.Map{
		"Title": "Succesfull email verifcation",
	})

}
