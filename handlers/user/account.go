package user

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"userstyles.world/database"
	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func Account(c *fiber.Ctx) error {
	u := jwt.User(c)

	styles, err := models.GetStylesByUser(database.DB, u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Server error",
			"User":  u,
		})
	}

	user, err := models.FindUserByName(database.DB, u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	return c.Render("account", fiber.Map{
		"Title":  "Account",
		"User":   u,
		"Params": user,
		"Styles": styles,
	})
}

func EditAccount(c *fiber.Ctx) error {
	u := jwt.User(c)

	styles, err := models.GetStylesByUser(database.DB, u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"User":  u,
			"Title": "Server error",
		})
	}

	user, err := models.FindUserByName(database.DB, u.Username)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "User not found",
			"User":  u,
		})
	}

	prevBio := user.Biography
	user.Biography = c.FormValue("bio")

	if err := utils.Validate().StructPartial(user, "Biography"); err != nil {
		errors := err.(validator.ValidationErrors)
		log.Println("Validation errors:", errors)
		user.Biography = prevBio

		return c.Render("account", fiber.Map{
			"Title":  "Validation Error",
			"User":   u,
			"Params": user,
			"Styles": styles,
			"Error":  "Validation error",
		})
	}

	t := new(models.User)

	dbErr := database.DB.
		Model(t).
		Where("id", user.ID).
		Updates(user).
		Error

	if dbErr != nil {
		log.Println("Updating style failed, err:", err)
		return c.Render("err", fiber.Map{
			"Title": "Internal server error.",
			"User":  u,
		})
	}

	return c.Render("account", fiber.Map{
		"Title":  "Account",
		"User":   u,
		"Params": user,
		"Styles": styles,
	})
}
