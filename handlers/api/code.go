package api

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/cache"
	"userstyles.world/modules/storage"
)

func statsMiddleware(c *fiber.Ctx) error {
	if !strings.HasSuffix(c.Path(), ".user.css") {
		c.Next()
	}

	id := strings.TrimPrefix(c.Path(), "/api/style/")
	id = strings.TrimSuffix(id, ".user.css")

	i, err := strconv.Atoi(id)
	if err != nil && i < 1 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid userstyle ID",
		})
	}

	cache.InstallStats.Add(c.IP() + " " + id)

	return c.Next()
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
