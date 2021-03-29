package main

import (
	"log"

	"userstyles.world/database"
	"userstyles.world/handlers"
	"userstyles.world/updater"
	"userstyles.world/utils"
)

func main() {
	err := utils.NewEmail().
		SetDefaultFrom().
		SetTo("HERE YOUR EMAIL").
		SetSubject("Verify your email address").
		SetPlainBody("Simple plain email").
		SetHTMLBody("<h3>Simple HTML body</h3>").
		SendEmail()

	if err != nil {
		log.Fatalf("Couldn't send email due to %s", err)
	}

	utils.InitializeValidator()
	database.Initialize()
	updater.Initialize()
	handlers.Initialize()
}
