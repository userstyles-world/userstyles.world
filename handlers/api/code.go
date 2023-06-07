package api

import (
	"fmt"
	"hash/crc32"

	"github.com/gofiber/fiber/v2"
	"userstyles.world/modules/storage"
)

func GetStyleSource(c *fiber.Ctx) error {
	return nil
}

func GetStyleEtag(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid userstyle ID",
		})
	}

	code, err := storage.FindStyleCode(id) // sftodo: use something file-related or etags won't work.
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
