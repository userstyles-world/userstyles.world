package main

import (
	"log"

	"userstyles.world/database"
	"userstyles.world/handlers"
	"userstyles.world/updater"
	"userstyles.world/utils"
)

func main() {
	PlainPart := utils.NewPart().
		SetBody("Simple plain email")

	HTMLPart := utils.NewPart().
		SetContentType("text/html; charset=\"utf-8\"").
		SetBody("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n" +
			"<h3>Simple HTML body It's an simple</h3>\n" +
			"<p style=\"color: green\">Check and make sure that everything works :D</p>")

	err := utils.NewEmail().
		SetTo("YOUR EMAIL ADDRESS HERE").
		SetSubject("Verify your email address").
		AddPart(*PlainPart).
		AddPart(*HTMLPart).
		SendEmail()

	if err != nil {
		log.Fatalf("Couldn't send email due to %s", err)
	}

	utils.InitializeValidator()
	database.Initialize()
	updater.Initialize()
	handlers.Initialize()
}
