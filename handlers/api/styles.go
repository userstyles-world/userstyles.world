package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ohler55/ojg/oj"
	"github.com/vednoc/go-usercss-parser"
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

var JsonParser = &oj.Parser{Reuse: true}

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
				"data": "This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──",
			})
	}
	var postStyle models.Style
	err = JsonParser.Unmarshal(c.Body(), &postStyle)
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

func DeleteStyle(c *fiber.Ctx) error {
	u, _ := APIUser(c)
	styleID := c.Params("id")

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
				"data": "This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──",
			})
	}

	styleModel := new(models.Style)
	err = database.DB.
		Debug().
		Delete(styleModel, "styles.id = ?", styleID).
		Error

	if err != nil {
		fmt.Printf("Failed to delete style, err: %#+v\n", err)
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't delete style",
			})
	}

	return c.JSON(fiber.Map{
		"data": "Succesful removed the style!",
	})
}

func NewStyle(c *fiber.Ctx) error {
	u, _ := APIUser(c)

	if !utils.Contains(u.Scopes, "style") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"style\" scope to do this.",
			})
	}

	var postStyle models.Style
	err := JsonParser.Unmarshal(c.Body(), &postStyle)
	if err != nil {
		fmt.Println(err)
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't convert style to struct.",
			})
	}

	if postStyle.Name == "" || postStyle.Code == "" || postStyle.Description == "" || postStyle.Category == "" {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Make sure to fill out fields.",
			})
	}
	postStyle.Featured = false
	postStyle.UserID = u.ID

	code := usercss.ParseFromString(postStyle.Code)
	if errs := usercss.BasicMetadataValidation(code); errs != nil {
		var errors string
		for i := 0; i < len(errs); i++ {
			errors += errs[i].Code.Error() + ";"
		}
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error:" + errors,
			})
	}

	// Prevent adding multiples of the same style.
	err = models.CheckDuplicateStyle(database.DB, &postStyle)
	if err != nil {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "A duplicate style was found.",
			})
	}

	newStyle, err := models.CreateStyle(database.DB, &postStyle)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't save style",
			})
	}

	return c.JSON(fiber.Map{
		"data": "Succesful added the style! With ID: " + fmt.Sprintf("%d", newStyle.ID),
	})
}
