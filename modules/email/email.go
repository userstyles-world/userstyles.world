// Email package provides helper utilities for rendering email templates.
package email

import (
	"bytes"

	"github.com/gofiber/fiber/v2"

	"userstyles.world/modules/config"
)

var views fiber.Views

// SetRenderer sets render engine for email templates.
func SetRenderer(app *fiber.App) {
	views = app.Config().Views
}

// Send renders templates with provided data, prepares an email, and sends it.
func Send(tmpl, address, title string, args any) error {
	var text bytes.Buffer
	err := views.Render(&text, "email/"+tmpl+".text", args)
	if err != nil {
		return err
	}

	var html bytes.Buffer
	err = views.Render(&html, "email/"+tmpl+".html", args)
	if err != nil {
		return err
	}

	// if !config.Production {
	// 	divider := strings.Repeat("-", 80)
	// 	log.Info.Printf("\n%s%s\n\n%s", text.String(), divider, html.String())
	// }

	return NewEmail().
		SetTo(address).
		SetSubject(title).
		AddPart(*NewPart().SetBody(text.String())).
		AddPart(*NewPart().SetBody(html.String()).HTML()).
		SendEmail(config.Secrets.EmailServer)
}
