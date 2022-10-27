package core

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
)

var client = http.Client{
	Timeout: time.Second * 30,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		// Max of three redirections.
		if len(via) >= 3 {
			return errors.New("*giggles* Mikey Wikey hates you")
		}

		// Make sure it doesn't redirect to a loopback thingy.
		if config.Production && utils.IsLoopback(string(req.Host)) {
			return errors.New("*giggles* Mikey Wikey hates you")
		}
		return nil
	},
}

func Proxy(c *fiber.Ctx) error {
	link, id, t := c.Query("link"), c.Query("id"), c.Query("type")

	// Don't render this page.
	if link == "" || id == "" || t == "" || strings.Contains(link, "..") {
		return c.Redirect("/", fiber.StatusSeeOther)
	}

	// Set resource location and name.
	dir := path.Join(config.ProxyDir, path.Clean(t), path.Clean(id))
	name := path.Join(dir, url.PathEscape(link))

	// Check if image exists.
	if _, err := os.Stat(name); os.IsNotExist(err) {
		// Create directory.
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Warn.Printf("Failed to create %q: %s\n", dir, err.Error())
			return nil
		}

		// Create new request.
		req, err := http.NewRequest("GET", link, nil)
		if err != nil {
			log.Info.Println("http.NewRequest: %w", err)
			return nil
		}

		// Ensure we're not doing a local host request.
		if utils.IsLoopback(string(req.Host)) {
			log.Info.Println("A local network was requested to be proxied.")
			return nil
		}

		// Make the actual request and get the response.
		res, err := client.Do(req)
		if err != nil {
			log.Info.Printf("Failed to get image %q, err: %v\n", link, err)
			return nil
		}
		defer res.Body.Close()

		// Make sure the response returned 200 OK.
		if res.StatusCode != 200 {
			log.Info.Printf("Failed to get image %q, didn't return 200 OK: %v\n", link, res.Status)
			return nil
		}

		// Limit image size to 8.388608 Megabytes ^-^, btw according to Mikey,
		// you need to learn bit shitting to understand this magic number.
		if res.ContentLength > 1<<23 {
			log.Info.Printf("Big image detected in %s: %q\n", id, link)
			return c.Redirect("/big-image.svg")
		}

		// Create a file(if none exist) and truncate it before writing to, open
		// it in a write-only manner.
		resFile, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
		if err != nil {
			log.Info.Printf("Failed to create file %q: %v\n", name, err)
			return nil
		}

		// Copy the response's body to the file.
		_, err = io.Copy(resFile, res.Body)
		// Close the file.
		resFile.Close()
		if err != nil {
			log.Info.Printf("Failed to copy: %v\n", err)
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
	re := regexp.MustCompile(`(?m).*"(https://.*)".*`)
	return re.ReplaceAllString(s, "$1")
}
