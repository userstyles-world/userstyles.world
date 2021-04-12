package api

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/config"
	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func CallbackGet(c *fiber.Ctx) error {
	// Get the necessary information.
	redirectCode, tempCode, state := c.Params("rcode"), c.Query("code"), c.Query("state")
	if redirectCode == "" || tempCode == "" {
		// Give them the bad enpoint error.
		return c.Next()
	}
	var service string
	var rState string
	if redirectCode != "gitlab" && redirectCode != "codeberg" {
		// Decode the string so we get our actual information back.
		code, err := utils.DecodePreparedText(redirectCode, utils.AEAD_OAUTH)
		if err != nil {
			return c.Next()
		}
		// We added the service within the the information and use the '+'
		// As seperator so now unseperate them.
		if splitted := strings.Split(code, "+"); len(splitted) == 2 {
			service, rState = splitted[0], splitted[1]
		} else {
			return c.Next()
		}

		if rState != state {
			return c.Next()
		}
	} else {
		service = redirectCode
	}

	response := utils.CallbackOAuth(tempCode, rState, service)
	if response == (utils.OAuthResponse{}) {
		return c.Next()
	}
	if service == "codeberg" {
		response.Name = response.GiteaName
	}

	user, err := models.FindUserByEmail(database.DB, response.Email)
	if err != nil {
		fmt.Println(err.Error())
		if err.Error() == "User not found." || err.Error() == "record not found" {
			regErr := database.DB.Create(&models.User{
				Username:      response.Name,
				OAuthProvider: service,
				Email:         response.Email,
			})

			if regErr.Error != nil {
				log.Printf("Failed to register %s, error: %s", response.Email, regErr.Error)

				c.SendStatus(fiber.StatusInternalServerError)
				return c.Render("err", fiber.Map{
					"Title": "Register failed",
					"Error": "Internal server error.",
				})
			}
			user, err = models.FindUserByEmail(database.DB, response.Email)
			if err != nil {
				log.Printf("Failed to register %s, error: %s", response.Email, err.Error())

				c.SendStatus(fiber.StatusInternalServerError)
				return c.Render("err", fiber.Map{
					"Title": "Register failed",
					"Error": "Internal server error.",
				})
			}
		} else {
			return c.Next()
		}
	}
	if user.OAuthProvider == "none" || user.OAuthProvider != service {
		return c.Next()
	}

	expiration := time.Hour * 24 * 14
	t, err := utils.NewJWTToken().
		SetClaim("id", user.ID).
		SetClaim("name", user.Username).
		SetClaim("email", user.Email).
		SetClaim("role", user.Role).
		SetExpiration(expiration).
		GetSignedString(nil)

	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     fiber.HeaderAuthorization,
		Value:    t,
		Path:     "/",
		Expires:  time.Now().Add(expiration),
		Secure:   config.DB != "dev.db",
		HTTPOnly: true,
		SameSite: "lax",
	})

	return c.Redirect("/account", fiber.StatusSeeOther)

}
