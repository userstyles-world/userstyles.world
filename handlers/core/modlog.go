package core

import (
	"github.com/userstyles-world/fiber/v2"

	"userstyles.world/handlers/jwt"
	"userstyles.world/models"
)

// GetModLog renders the modlog view.
// It will pass trough the relevant information from the database.
func GetModLog(c *fiber.Ctx) error {
	u, _ := jwt.User(c)

	bannedUsers, err := models.GetLogOfKind(models.LogBanUser)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Internal Server error",
			"User":  u,
		})
	}

	removedStyles, err := models.GetLogOfKind(models.LogRemoveStyle)
	if err != nil {
		return c.Render("err", fiber.Map{
			"Title": "Internal Server error",
			"User":  u,
		})
	}

	return c.Render("core/modlog", fiber.Map{
		"BannedUsers":   bannedUsers,
		"RemovedStyles": removedStyles,
		"User":          u,
		"Title":         "Mod Log",
		"Canonical":     "modlog",
	})
}
