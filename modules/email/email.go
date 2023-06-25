// Email package provides helper utilities for rendering email templates.
package email

import (
	"io"

	"github.com/gofiber/fiber/v2"
)

var views fiber.Views

// SetRenderer sets render engine for email templates.
func SetRenderer(app *fiber.App) {
	views = app.Config().Views
}

// Render rendres both kinds of email templates with provided data.
func Render(bufText, bufHTML io.Writer, name string, args any) error {
	err := views.Render(bufText, "email/"+name+".text", args)
	if err != nil {
		return err
	}
	err = views.Render(bufHTML, "email/"+name+".html", args)
	if err != nil {
		return err
	}

	return nil
}
