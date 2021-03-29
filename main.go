package main

import (
	"userstyles.world/database"
	"userstyles.world/handlers"
	"userstyles.world/updater"
	"userstyles.world/utils"
)

func main() {
	utils.InitalizeCrypto()
	utils.InitializeValidator()
	database.Initialize()
	updater.Initialize()
	handlers.Initialize()
}
