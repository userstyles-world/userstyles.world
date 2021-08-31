package core

import (
	"fmt"
	"net/url"
	"os"

	"github.com/userstyles-world/fiber/v2"

	"userstyles.world/modules/log"
)

func Proxy(c *fiber.Ctx) error {
	path, id, t := c.Query("url"), c.Query("id"), c.Query("type")

	// Set resource location and name.
	loc := fmt.Sprintf("./data/proxy/%s/%s", t, id)
	name := loc + "/" + url.PathEscape(path)

	stat, err := os.Stat(name)
	if os.IsNotExist(err) {
		a := fiber.AcquireAgent()
		req := a.Request()
		req.SetRequestURI(path)
		if err := a.Parse(); err != nil {
			panic(err) // TODO: Handle this error properly.
		}

		// TODO: Add a "not found" image.
		_, data, _ := a.Bytes()

		// Create directory.
		stat, err := os.Stat(loc)
		if os.IsNotExist(err) {
			if err := os.Mkdir(loc, 0o755); err != nil {
				log.Warn.Fatal(err)
			}
		}
		if err := os.WriteFile(name, data, 0o600); err != nil {
			log.Warn.Println("Failed to write image:", err.Error())
			return fmt.Errorf("failed to write image: %v", err)
		}
	}

	// Serve image.
	f, err := os.Open(name)
	if err != nil {
		log.Warn.Fatal(err)
	}

	if stat, err = f.Stat(); err != nil {
		return c.JSON(err)
	}

	c.Response().SetBodyStream(f, int(stat.Size()))

	return nil
}
