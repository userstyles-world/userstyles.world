package core

import (
	"fmt"
	"net/url"
	"os"

	"github.com/userstyles-world/fiber/v2"

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
	stat, err := os.Stat(name)
	if os.IsNotExist(err) {
		// Create directory.
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Warn.Printf("Failed to create %v: %s\n", dir, err.Error())
			return nil
		}

		// Download image.
		a := fiber.AcquireAgent()
		req := a.Request()
		req.SetRequestURI(link)
		if err := a.Parse(); err != nil {
			log.Info.Println("Agent err:", err.Error())
			return nil
		}

		// TODO: Show a fallback image.
		_, data, errs := a.Bytes()
		if len(errs) > 0 {
			log.Info.Printf("Failed to get image: %v\n", errs)
			return nil
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

	if stat, err = f.Stat(); err != nil {
		log.Info.Println("Failed to get stat:", err.Error())
		return nil
	}

	c.Response().SetBodyStream(f, int(stat.Size()))

	return nil
}
