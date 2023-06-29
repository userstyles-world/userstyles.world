package api

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ohler55/ojg/oj"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/database"
	"userstyles.world/modules/images"
	"userstyles.world/modules/log"
	"userstyles.world/modules/search"
	"userstyles.world/modules/storage"
	"userstyles.world/modules/util"
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

	id, err := c.ParamsInt("id")
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

	if u.StyleID != 0 && uint(id) != u.StyleID {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Error: This style doesn't belong to you! ╰༼⇀︿⇀༽つ-]═──",
			})
	}

	style, err := models.GetStyleByID(c.Params("id"))
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
	postStyle.Code = util.RemoveUpdateURL(postStyle.Code)

	msg, err := postStyle.ValidateCode(utils.Validate(), true)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"data": msg})
	}

	err = models.UpdateStyle(&postStyle)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't save style",
			})
	}

	// TODO: Benchmark this approach.
	if style.Code != postStyle.Code {
		err = models.SaveStyleCode(strconv.Itoa(int(postStyle.ID)), postStyle.Code)
		if err != nil {
			log.Warn.Printf("kind=code id=%v err=%q\n", postStyle.ID, err)
		}
		cache.Code.Remove(id)
	}

	if err = search.IndexStyle(postStyle.ID); err != nil {
		log.Warn.Printf("Failed to re-index style %d: %s\n", postStyle.ID, err)
	}

	return c.JSON(fiber.Map{
		"data": "Successfully edited style.",
	})
}

func DeleteStyle(c *fiber.Ctx) error {
	u, _ := User(c)
	id := c.Params("id")

	i, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(403).
			JSON(fiber.Map{
				"data": "Couldn't parse param \"id\"",
			})
	}

	style, err := models.GetStyleByID(id)
	if err != nil {
		return c.Status(500).
			JSON(fiber.Map{
				"data": "Error: Couldn't find style with ID.",
			})
	}

	if u.StyleID != 0 && uint(i) != u.StyleID {
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
		Delete(styleModel, "styles.id = ?", id).
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

	err = models.RemoveStyleCode(strconv.Itoa(int(style.ID)))
	if err != nil {
		log.Warn.Printf("kind=removecode id=%v err=%q\n", style.ID, err)
	}

	cache.Code.Remove(i)

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

	postStyle.UserID = u.ID
	postStyle.Featured = false
	postStyle.Code = util.RemoveUpdateURL(postStyle.Code)

	msg, err := postStyle.ValidateCode(utils.Validate(), true)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"data": msg})
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

	err = models.SaveStyleCode(strconv.Itoa(int(s.ID)), s.Code)
	if err != nil {
		log.Warn.Printf("kind=code id=%v err=%q\n", s.ID, err)
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
