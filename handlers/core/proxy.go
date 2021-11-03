package core

import (
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/log"
)

func Proxy(c *fiber.Ctx) error {
	link, id, t := c.Query("link"), c.Query("id"), c.Query("type")

	// Don't render this page.
	if link == "" || id == "" || t == "" {
		return c.Redirect("/", fiber.StatusSeeOther)
	}

	// Set resource location and name.
	dir := fmt.Sprintf("./data/proxy/%s/%s", t, id)
	name := dir + "/" + url.PathEscape(link)

	// Check if image exists.
	if _, err := os.Stat(name); os.IsNotExist(err) {
		// Create directory.
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Warn.Printf("Failed to create %v: %s\n", dir, err.Error())
			return nil
		}

		// Download image.
		a := fiber.AcquireAgent()
		defer fiber.ReleaseAgent(a)

		var status int
		var data []byte
		var errs []error

		// HACK: GitHub doesn't set "Location" response header.
		if strings.Contains(link, "https://github.com/") {
		getImage:
			a.Request().SetRequestURI(link)
			if err := a.Parse(); err != nil {
				log.Info.Println("Agent err:", err.Error())
				return nil
			}

			status, data, errs = a.Bytes()
			if len(errs) > 0 {
				log.Info.Printf("Failed to get image %v, err: %v\n", link, errs)
				return nil
			}

			if status >= 300 && status <= 400 {
				link = extractImage(string(data))
				goto getImage
			}
		} else {
			a.Request().SetRequestURI(link)
			a.MaxRedirectsCount(3)

			if err := a.Parse(); err != nil {
				log.Info.Println("Agent err:", err.Error())
				return nil
			}

			// TODO: Show a fallback image.
			status, data, errs = a.Bytes()
			if len(errs) > 0 {
				log.Info.Printf("Failed to get image %v, err: %v\n", link, errs)
				return nil
			}
		}

		if err := os.WriteFile(name, data, 0o600); err != nil {
			log.Info.Println("Failed to write image:", err.Error())
			return nil
		}
	}

	// Serve image.
	f, err := os.Open(name)
	if err != nil {
		log.Info.Println("Failed to open image:", err.Error())
		return nil
	}

	stat, err := f.Stat()
	if err != nil {
		log.Info.Println("Failed to get stat:", err.Error())
		return nil
	}

	c.Response().SetBodyStream(f, int(stat.Size()))

	return nil
}

func extractImage(s string) string {
	fmt.Printf("%#v\n", s)
	re := regexp.MustCompile(`(?m).*"(https://.*)".*`)
	return re.ReplaceAllString(s, "$1")
}
