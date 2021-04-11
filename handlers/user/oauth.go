package user

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func AuthLoginGet(c *fiber.Ctx) error {
	if u, ok := jwt.User(c); ok {
		log.Printf("User %d has set session, redirecting.", u.ID)
		c.Redirect("/account", fiber.StatusSeeOther)
	}

	oathType := c.Params("type")
	redirectURI := ""

	switch oathType {
	case "github":
		redirectURI = utils.GithubMakeURL(c.BaseURL())
	}

	return c.Redirect(redirectURI, fiber.StatusSeeOther)
}

func AuthRegisterPost(c *fiber.Ctx) error {
	state := c.Params("state")
	u := models.User{
		Username: c.FormValue("username"),
		Password: c.FormValue("password"),
		Email:    c.FormValue("email"),
	}

	if state == "" || u.Email == "" {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
		})
	}

	decodedState, err := utils.DecodePreparedText(state, utils.AEAD_OAUTH)
	if err != nil {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
		})
	}
	var verified string
	var stateEmail string
	if splitted := strings.Split(decodedState, "+"); len(splitted) == 2 {
		verified, stateEmail = splitted[0], splitted[1]
	} else {
		return c.Next()
	}

	if stateEmail != u.Email {
		c.SendStatus(fiber.StatusInternalServerError)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error",
		})
	}

	if verified == "verified" {

		err := utils.Validate().StructPartial(u, "Username", "Email", "Password")
		if err != nil {
			errors := err.(validator.ValidationErrors)
			log.Println("Validation errors:", errors)

			c.SendStatus(fiber.StatusInternalServerError)
			return c.Render("more_info", fiber.Map{
				"Title": "Register failed",
				"Error": "Failed to register. Make sure your input is correct.",
			})
		}
		u.Password = utils.GenerateHashedPassword(u.Password)
		fmt.Println(u)
		regErr := database.DB.Create(&u)

		if regErr.Error != nil {
			log.Printf("Failed to register %s, error: %s", u.Email, regErr.Error)

			c.SendStatus(fiber.StatusInternalServerError)
			return c.Render("err", fiber.Map{
				"Title": "Register failed",
				"Error": "Internal server error.",
			})
		}

		return c.Render("verification", fiber.Map{
			"Title":        "Successful register",
			"Verification": "Successful registered",
			"Reason":       "You've successfully made a acccount",
		})

	} else {
		c.Locals("Email", u.Email)
		return RegisterPost(c)
	}

}
