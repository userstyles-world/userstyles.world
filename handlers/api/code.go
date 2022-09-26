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
	"userstyles.world/modules/storage"
)

func GetStyleSource(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}

	key := strconv.Itoa(id)
	val, found := cache.LRU.Get(key)
	if !found {
		code, err := storage.FindStyleCode(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"data": "style not found",
			})
		}

		// Override updateURL field to prevent abuse.
		url := config.BaseURL + "/api/style/" + key + ".user.css"
		code = usercss.OverrideUpdateURL(code, url)

		// Cache the userstyle.
		cache.LRU.Add(key, code)

		// Reassign code var.
		val = code
	}

	// Upsert style installs.
	s := new(models.Stats)
	if err := s.UpsertInstall(key, c.IP()); err != nil {
		log.Database.Printf("Failed to upsert installs for %q: %s\n", key, err)
	}

	c.Type("css", "utf-8")
	return c.SendString(val.(string))
}

func GetStyleEtag(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}

	code, err := storage.FindStyleCode(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "style not found",
		})
	}

	// Follows the format "source code length - MD5 Checksum of source code"
	val := fmt.Sprintf(`%v-%v`, len(code), crc32.ChecksumIEEE([]byte(code)))

	// Set the value for "Etag" header
	c.Set("etag", val)
	return nil
}
