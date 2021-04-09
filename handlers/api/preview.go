package api

import (
	"io/fs"
	"os"
	"strings"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/gofiber/fiber/v2"
	"userstyles.world/images"
)

func getFileExtension(path string) string {
	n := strings.LastIndexByte(path, '.')
	if n < 0 {
		return ""
	}
	return path[n:]
}

var notFound = func(c *fiber.Ctx) error {
	c.Status(fiber.StatusNotFound)
	return c.SendString("Screenshot not found")
}

func GetPreviewScreenshot(c *fiber.Ctx) error {
	styleId := c.Params("id")
	format := getFileExtension(styleId)
	styleId = strings.TrimSuffix(styleId, format)

	info, err := images.GetImageFromStyle(styleId)
	if err != nil {
		return notFound(c)
	}

	var stat fs.FileInfo
	var fileName string
	orignialFile := images.CacheFolder + styleId + ".originial"

	switch format[1:] {
	case "jpeg":
		if info.Jpeg == nil {
			return notFound(c)
		}
		stat = info.Originial
		fileName = images.CacheFolder + styleId + ".jpeg"
	case "webp":
		fileName = images.CacheFolder + styleId + ".webp"
		if info.WebP == nil {
			err = images.DecodeImage(orignialFile, fileName, vips.ImageTypeWEBP)
			if err != nil {
				return notFound(c)
			}
			webpStat, err := os.Stat(fileName)
			if err != nil {
				return notFound(c)
			}
			stat = webpStat
			break
		}
		stat = info.WebP
	}

	if stat == nil || fileName == "" {
		return notFound(c)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return notFound(c)
	}
	c.Type(getFileExtension(stat.Name()))
	c.Response().SetBodyStream(file, int(stat.Size()))

	return nil
}
