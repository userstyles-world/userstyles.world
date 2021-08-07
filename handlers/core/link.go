package core

import (
	"github.com/gofiber/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/modules/config"
)

func GetLinkedSite(c *fiber.Ctx) error {
	switch c.Params("site") {
	case "discord":
		return c.Redirect(config.AppLinkChatDiscord, fiber.StatusSeeOther)
	case "matrix":
		return c.Redirect(config.AppLinkChatMatrix, fiber.StatusSeeOther)
	case "source":
		return c.Redirect(config.AppLinkSource, fiber.StatusSeeOther)
	default:
		u, _ := jwt.User(c)
		return c.Render("err", fiber.Map{
			"Title": "Invalid link",
			"User":  u,
		})
	}
}
