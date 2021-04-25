package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ohler55/ojg/oj"
	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func StylesGet(c *fiber.Ctx) error {
	u, _ := APIUser(c)

	if !utils.Contains(u.Scopes, "style") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"style\" scope to do this.",
			})
	}

	styles, err := models.GetStylesByUser(database.DB, u.Username)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't find styles.",
			})
	}

	return c.JSON(fiber.Map{
		"data": styles,
	})

}

var p = &oj.Parser{Reuse: true}

func StylePost(c *fiber.Ctx) error {
	u, _ := APIUser(c)
	styleID := c.Params("id")

	if !utils.Contains(u.Scopes, "style") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"style\" scope to do this.",
			})
	}
	style, err := models.GetStyleByID(database.DB, styleID)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't find style with ID.",
			})
	}
	if style.UserID != u.ID {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "This style doesn't belong to this user.",
			})
	}
	var postStyle models.Style
	err = p.Unmarshal(c.Body(), &postStyle)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't convert style to struct.",
			})
	}

	// Just to prevent from weird peeps doing this shit.
	postStyle.ID = style.ID
	postStyle.UserID = u.ID
	postStyle.Featured = style.Featured

	err = models.UpdateStyle(database.DB, &postStyle)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't save style",
			})
	}

	return c.JSON(fiber.Map{
		"data": "Succesful edited style!",
	})

}
