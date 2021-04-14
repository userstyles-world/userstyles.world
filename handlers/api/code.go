package api

import (
	"hash/crc32"

	"github.com/gofiber/fiber/v2"
	"github.com/vednoc/go-usercss-parser"

	"userstyles.world/database"
	"userstyles.world/models"
	"userstyles.world/utils"
)

func GetStyleSource(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(database.DB, id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	// Override updateURL field for Stylus integration.
	// TODO: Also override it in the database on demand?
	uc := usercss.ParseFromString(style.Code)
	url := "https://userstyles.world/api/style/" + id + ".user.css"
	uc.OverrideUpdateURL(url)

	// Count installs.
	_, err = models.AddStatsToStyle(database.DB, id, c.IP(), true)
	if err != nil {
		return c.JSON(fiber.Map{
			"data": "Internal server error",
		})
	}

	c.Set("Content-Type", "text/css")
	return c.SendString(uc.SourceCode)
}

var normalizedHeaderETag = []byte("Etag")

func GetStyleEtag(c *fiber.Ctx) error {
	id := c.Params("id")

	style, err := models.GetStyleSourceCodeAPI(database.DB, id)
	if err != nil {
		return c.JSON(fiber.Map{"data": "style not found"})
	}

	// TODO: add a possible update stat?
	// TODO: internal switch to a byte pool to avoid allocations over time.
	etagValue := make([]byte, 0, 40)
	etagValue = appendUint(etagValue, uint32(len(style.Code)))
	etagValue = append(etagValue, '-')
	etagValue = appendUint(etagValue, crc32.ChecksumIEEE(utils.S2b(style.Code)))

	c.Response().Header.SetCanonical(normalizedHeaderETag, etagValue)
	return nil
}

// appendUint appends n to dst and returns the extended dst.
func appendUint(dst []byte, n uint32) []byte {
	var b [20]byte
	buf := b[:]
	i := len(buf)
	var q uint32
	for n >= 10 {
		i--
		q = n / 10
		buf[i] = '0' + byte(n-q*10)
		n = q
	}
	i--
	buf[i] = '0' + byte(n)

	dst = append(dst, buf[i:]...)
	return dst
}
