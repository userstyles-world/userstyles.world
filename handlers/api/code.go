package api

import (
	"fmt"
	"hash/crc32"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

func GetStyleSource(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	// Override updateURL field for Stylus integration.
	// TODO: Also override it in the database on demand?
	uc := usercss.ParseFromString(style.Code)
	url := "https://userstyles.world/api/style/" + id + ".user.css"
	uc.OverrideUpdateURL(url)

	// Upsert style installs.
	go func(id, ip string) {
		s := new(models.Stats)
		if err := s.CreateRecord(id, ip); err != nil {
			log.Warn.Printf("Failed to create record for %s: %s\n", id, err.Error())
		}
		if err := s.UpsertInstall(); err != nil {
			log.Warn.Printf("Failed to upsert install for %v: %s\n", s.StyleID, err.Error())
		}
	}(id, c.IP())

	c.Set("Content-Type", "text/css")
	return c.SendString(uc.SourceCode)
}

var normalizedHeaderETag = []byte("Etag")

func GetStyleEtag(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	// Follows the format "source code length - MD5 Checksum of source code"
	etagValue := fmt.Sprintf("\"%v-%v\"", len(style.Code), crc32.ChecksumIEEE([]byte(style.Code)))

	// Set the value for "Etag" header
	c.Response().Header.SetCanonical(normalizedHeaderETag, []byte(etagValue))
	return nil
}
