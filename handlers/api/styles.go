package api

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ohler55/ojg/oj"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/modules/database"
	"userstyles.world/models"
	"userstyles.world/search"
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

	styles, err := models.GetStylesByUser(u.Username)
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

// JSONParser defined options.
var JSONParser = &oj.Parser{Reuse: true}

func StylePost(c *fiber.Ctx) error {
	u, _ := APIUser(c)
	sStyleID := c.Params("id")

	styleID, err := strconv.Atoi(sStyleID)
	if err != nil {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Couldn't parse param \"id\"",
			})
	}

	if u.StyleID == 0 && !utils.Contains(u.Scopes, "style") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"style\" scope to do this.",
			})
	}

	if u.StyleID != 0 && uint(styleID) != u.StyleID {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──",
			})
	}

	style, err := models.GetStyleByID(sStyleID)
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
	err = JSONParser.Unmarshal(c.Body(), &postStyle)
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

	err = models.UpdateStyle(&postStyle)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't save style",
			})
	}

	if err = search.IndexStyle(postStyle.ID); err != nil {
		log.Printf("Re-indexing style %d failed, err: %s", postStyle.ID, err.Error())
	}

	return c.JSON(fiber.Map{
		"data": "Successful edited style!",
	})
}

func DeleteStyle(c *fiber.Ctx) error {
	u, _ := APIUser(c)
	sStyleID := c.Params("id")

	styleID, err := strconv.Atoi(sStyleID)
	if err != nil {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Couldn't parse param \"id\"",
			})
	}

	style, err := models.GetStyleByID(sStyleID)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't find style with ID.",
			})
	}

	if u.StyleID != 0 && uint(styleID) != u.StyleID {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──",
			})
	}

	if style.UserID != u.ID {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──",
			})
	}

	styleModel := new(models.Style)
	err = database.Conn.
		Debug().
		Delete(styleModel, "styles.id = ?", sStyleID).
		Error

	if err != nil {
		log.Printf("Failed to delete style, err: %#+v\n", err)
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't delete style",
			})
	}

	if err = search.DeleteStyle(style.ID); err != nil {
		log.Printf("Couldn't delete style %d failed, err: %s", style.ID, err.Error())
	}

	return c.JSON(fiber.Map{
		"data": "Successful removed the style!",
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
	err := JSONParser.Unmarshal(c.Body(), &postStyle)
	if err != nil {
		log.Println(err)
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
	err = models.CheckDuplicateStyle(&postStyle)
	if err != nil {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "A duplicate style was found.",
			})
	}

	newStyle, err := models.CreateStyle(&postStyle)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't save style",
			})
	}

	if err = search.IndexStyle(postStyle.ID); err != nil {
		log.Printf("Re-indexing style %d failed, err: %s", postStyle.ID, err.Error())
	}

	return c.JSON(fiber.Map{
		"data": "Successful added the style! With ID: " + fmt.Sprintf("%d", newStyle.ID),
	})
}
