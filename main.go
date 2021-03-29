package main

import (
	"fmt"

	"userstyles.world/database"
	"userstyles.world/handlers"
	"userstyles.world/updater"
	"userstyles.world/utils"
)

func main() {
	if err := utils.SendEmail("YOUR_EMAIL_HERE!!! VEDNOC!!!", "dummy message"); err != nil {
		fmt.Printf("Got error %s", err)
		return
	}
	utils.InitializeValidator()
	database.Initialize()
	updater.Initialize()
	handlers.Initialize()
}
