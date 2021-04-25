package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func roleToString(role models.Role) string {
	switch role {
	case 1:
		return "Moderator"
	case 2:
		return "Admin"
	default:
		return "Regular"
	}
}

func UserGet(c *fiber.Ctx) error {
	u, _ := APIUser(c)

	if !utils.Contains(u.Scopes, "user") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"user\" scope to do this.",
			})
	}

	user, err := models.FindUserByName(database.DB, u.Username)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't find user.",
			})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"username":  user.Username,
			"email":     user.Email,
			"biography": user.Biography,
			"role":      roleToString(user.Role),
			"socials": fiber.Map{
				"github":   user.Socials.Github,
				"gitlab":   user.Socials.Gitlab,
				"codeberg": user.Socials.Codeberg,
			},
		},
	})
}

func SpecificUserGet(c *fiber.Ctx) error {
	identifier := c.Params("identifier")

	var user *models.User
	var err error
	if _, intErr := strconv.Atoi(identifier); intErr == nil {
		user, err = models.FindUserByID(database.DB, identifier)
	} else {
		user, err = models.FindUserByName(database.DB, identifier)
	}
	if err != nil {
		return c.Status(404).
			JSON(fiber.Map{
				"data": "User not found.",
			})
	}

	return c.JSON(fiber.Map{
		"data": fiber.Map{
			"username":  user.Username,
			"biography": user.Biography,
			"role":      roleToString(user.Role),
			"socials": fiber.Map{
				"github":   user.Socials.Github,
				"gitlab":   user.Socials.Gitlab,
				"codeberg": user.Socials.Codeberg,
			},
		},
	})
}
