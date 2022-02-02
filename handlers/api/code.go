package api

import (
	"fmt"
	"hash/crc32"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/cache"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

func GetStyleSource(c *fiber.Ctx) error {
	i, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}
	id := strconv.Itoa(i)

	code, found := cache.LRU.Get(id)
	if !found {
		style, err := models.GetStyleSourceCodeAPI(id)
		if err != nil {
			return c.JSON(fiber.Map{"data": "style not found"})
		}

		// Override updateURL field to prevent abuse.
		url := config.BaseURL + "/api/style/" + id + ".user.css"
		src := usercss.OverrideUpdateURL(style.Code, url)

		// Cache the userstyle.
		cache.LRU.Add(id, src)

		// Reassign code var.
		code = src
	}

	// Upsert style installs.
	go func(id, ip string) {
		s := new(models.Stats)
		if err := s.UpsertInstall(id, ip); err != nil {
			log.Warn.Printf("Failed to upsert install for %s, %s\n", id, err)
		}
	}(id, c.IP())

	c.Type("text/css", "utf-8")
	return c.SendString(code.(string))
}

func GetStyleEtag(c *fiber.Ctx) error {
	i, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}
	id := strconv.Itoa(i)

	style, err := models.GetStyleSourceCodeAPI(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "style not found",
		})
	}

	// Follows the format "source code length - MD5 Checksum of source code"
	etagValue := fmt.Sprintf("\"%v-%v\"", len(style.Code), crc32.ChecksumIEEE([]byte(style.Code)))

	// Set the value for "Etag" header
	c.Set("etag", etagValue)
	return nil
}
