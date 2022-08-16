package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ohler55/ojg/oj"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/database"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/modules/search"
	"userstyles.world/modules/storage"
	"userstyles.world/utils"
)

func StylesGet(c *fiber.Ctx) error {
	u, _ := User(c)

	if !utils.ContainsString(u.Scopes, "style") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"style\" scope to do this.",
			})
	}

	styles, err := storage.FindStyleCardsForUsername(u.Username)
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
	u, _ := User(c)
	sStyleID := c.Params("id")

	styleID, err := strconv.Atoi(sStyleID)
	if err != nil {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: Couldn't parse param \"id\"",
			})
	}

	if u.StyleID == 0 && !utils.ContainsString(u.Scopes, "style") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: You need the \"style\" scope to do this.",
			})
	}

	if u.StyleID != 0 && uint(styleID) != u.StyleID {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──",
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
				"data": "Error: This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──",
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

	uc := new(usercss.UserCSS)
	if err := uc.Parse(postStyle.Code); err != nil {
		return c.Status(403).JSON(fiber.Map{"data": "Error: " + err.Error()})
	}
	if errs := uc.Validate(); errs != nil {
		var errors string
		for i := 0; i < len(errs); i++ {
			errors += errs[i].Code.Error() + ";"
		}
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: " + errors,
			})
	}

	// Prevent broken traditional userstyles.
	// TODO: Remove a week or two after Stylus v1.5.20 is released.
	if len(uc.MozDocument) == 0 {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: Bad style format (visit https://userstyles.world/docs/faq#bad-style-format-error)",
			})
	}

	err = models.UpdateStyle(&postStyle)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't save style",
			})
	}

	// Evict from LRU cache when source code is updated.
	if style.Code != postStyle.Code {
		cache.LRU.Remove(sStyleID)
	}

	go func(style *models.Style) {
		if err = search.IndexStyle(style.ID); err != nil {
			log.Warn.Printf("Failed to re-index style %d: %s\n", style.ID, err)
		}
	}(&postStyle)

	return c.JSON(fiber.Map{
		"data": "Successfully edited style.",
	})
}

func DeleteStyle(c *fiber.Ctx) error {
	u, _ := User(c)
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
		Delete(styleModel, "styles.id = ?", sStyleID).
		Error

	if err != nil {
		log.Warn.Println("Failed to delete style from database:", err.Error())
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't delete style",
			})
	}

	if err = search.DeleteStyle(style.ID); err != nil {
		log.Warn.Printf("Failed to delete style %d from index: %s", style.ID, err)
	}

	return c.JSON(fiber.Map{
		"data": "Successful removed the style!",
	})
}

func NewStyle(c *fiber.Ctx) error {
	u, _ := User(c)

	if !utils.ContainsString(u.Scopes, "style") {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "You need the \"style\" scope to do this.",
			})
	}

	var postStyle models.Style
	err := JSONParser.Unmarshal(c.Body(), &postStyle)
	if err != nil {
		log.Warn.Println("Failed to convert new style to a struct")
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't convert style to struct.",
			})
	}

	if postStyle.Name == "" || postStyle.Code == "" || postStyle.Description == "" || postStyle.Category == "" {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: Make sure to fill out fields.",
			})
	}
	postStyle.Featured = false
	postStyle.UserID = u.ID

	uc := new(usercss.UserCSS)
	if err := uc.Parse(postStyle.Code); err != nil {
		return c.Status(403).JSON(fiber.Map{"data": "Error: " + err.Error()})
	}
	if errs := uc.Validate(); errs != nil {
		var errors string
		for i := 0; i < len(errs); i++ {
			errors += errs[i].Code.Error() + ";"
		}
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error:" + errors,
			})
	}

	// Prevent broken traditional userstyles.
	// TODO: Remove a week or two after Stylus v1.5.20 is released.
	if len(uc.MozDocument) == 0 {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: Bad style format (visit https://userstyles.world/docs/faq#bad-style-format-error)",
			})
	}

	// Prevent adding multiples of the same style.
	err = models.CheckDuplicateStyle(&postStyle)
	if err != nil {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: A duplicate style was found.",
			})
	}

	s, err := models.CreateStyle(&postStyle)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't save style",
			})
	}

	// Check preview image.
	file, _ := c.FormFile("preview")
	styleID := strconv.FormatUint(uint64(s.ID), 10)
	if file != nil || s.Preview != "" {
		if err = images.Generate(file, styleID, "0", "", s.Preview); err != nil {
			log.Warn.Println("Error:", err)
			s.Preview = ""
		} else {
			s.SetPreview()
			if err = s.UpdateColumn("preview", s.Preview); err != nil {
				log.Warn.Printf("Failed to update preview for %d: %s\n", s.ID, err)
			}
		}
	}

	if err = search.IndexStyle(s.ID); err != nil {
		log.Warn.Printf("Failed to index style %d: %s", s.ID, err)
	}

	return c.JSON(fiber.Map{
		"data": "Successfully added the style. ID: " + styleID,
	})
}
