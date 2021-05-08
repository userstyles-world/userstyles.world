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

func setSocials(u *models.User, k, v string) models.User {
	switch k {
	case "github":
		if v != "" {
			u.Socials.Github = v
		}
	case "gitlab":
		if v != "" {
			u.Socials.Gitlab = v
		}
	case "codeberg":
		if v != "" {
			u.Socials.Codeberg = v
		}
	}

	return *u
}

func Account(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

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
	u, _ := jwt.User(c)

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

	name := c.FormValue("name")
	if name != "" {
		prev := user.DisplayName
		user.DisplayName = name

		if err := utils.Validate().StructPartial(user, "DisplayName"); err != nil {
			errors := err.(validator.ValidationErrors)
			log.Println("Validation errors:", errors)
			user.DisplayName = prev

			return c.Render("account", fiber.Map{
				"Title":  "Validation Error",
				"User":   u,
				"Params": user,
				"Styles": styles,
				"Error":  "Display name must be longer than 5 and shorter than 20 characters.",
			})
		}
	}

	bio := c.FormValue("bio")
	if bio != "" {
		prev := user.Biography
		user.Biography = bio

		if err := utils.Validate().StructPartial(user, "Biography"); err != nil {
			errors := err.(validator.ValidationErrors)
			log.Println("Validation errors:", errors)
			user.Biography = prev

			return c.Render("account", fiber.Map{
				"Title":  "Validation Error",
				"User":   u,
				"Params": user,
				"Styles": styles,
				"Error":  "Biography must be shorter than 512 characters.",
			})
		}
	}

	setSocials(user, "github", c.FormValue("github"))
	setSocials(user, "gitlab", c.FormValue("gitlab"))
	setSocials(user, "codeberg", c.FormValue("codeberg"))

	dbErr := database.DB.
		Model(models.User{}).
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
