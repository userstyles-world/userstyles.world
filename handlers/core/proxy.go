package core

import (
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
	"userstyles.world/utils"
)

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
			log.Warn.Printf("Failed to create %v: %s\n", dir, err.Error())
			return nil
		}

		// Download image.
		a := fiber.AcquireAgent()
		defer fiber.ReleaseAgent(a)

		var status int
		var data []byte
		var errs []error

	getImage:
		// Set the request URI.
		a.Request().SetRequestURI(link)

		// Parse the request URI.
		if err := a.Parse(); err != nil {
			log.Info.Println("Agent err:", err.Error())
			return nil
		}

		// Ensure we're not doing a local host request.
		if utils.IsLoopback(string(a.Request().URI().Host())) {
			log.Info.Println("A local network was requested to be proxied.")
			return nil
		}

		// If we don't request to github.com set a max redirect of three.
		if !strings.Contains(link, "https://github.com/") {
			a.MaxRedirectsCount(3)
		}

		// Make the actual request and get the status and bytes.
		status, data, errs = a.Bytes()
		if len(errs) > 0 {
			log.Info.Printf("Failed to get image %v, err: %v\n", link, errs)
			return nil
		}

		// Check after all redirections if the host is still valid.
		if utils.IsLoopback(string(a.Request().URI().Host())) {
			log.Info.Println("A local network was requested to be proxied.")
			return nil
		}

		// HACK: GitHub doesn't set "Location" response header.
		if strings.Contains(link, "https://github.com/") && status >= 300 && status < 400 {
			link = extractImage(string(data))
			goto getImage
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
	re := regexp.MustCompile(`(?m).*"(https://.*)".*`)
	return re.ReplaceAllString(s, "$1")
}
