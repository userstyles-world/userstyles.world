package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/models"
	"userstyles.world/utils"
)

func roleToString(role models.Role) string {
	switch role {
	case models.Moderator:
		return "Moderator"
	case models.Admin:
		return "Admin"
	case models.Regular:
		return "Regular"
	}
	return "Invalid User"
}

func UserGet(c *fiber.Ctx) error {
	u, _ := User(c)

	if !utils.ContainsString(u.Scopes, "user") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"user\" scope to do this.",
			})
	}

	user, err := models.FindUserByName(u.Username)
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
			"socials":   user.Socials,
		},
	})
}

func SpecificUserGet(c *fiber.Ctx) error {
	identifier := c.Params("identifier")

	var user *models.User
	var err error
	if _, intErr := strconv.Atoi(identifier); intErr == nil {
		user, err = models.FindUserByID(identifier, "HACK")
	} else {
		user, err = models.FindUserByName(identifier)
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
			"socials":   user.Socials,
		},
	})
}
