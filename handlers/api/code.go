package api

import (
	"fmt"
	"hash/crc32"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/cache"
	"userstyles.world/modules/storage"
)

func GetStyleSource(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}

	code, err := storage.FindStyleCode(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"data": "style not found",
		})
	}

	cache.InstallStats.Add(c.IP() + " " + strconv.Itoa(id))

	c.Type("css", "utf-8")
	return c.SendString(code)
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
