package main

import (
	"userstyles.world/database"
	"userstyles.world/handlers"
	"userstyles.world/images"
	"userstyles.world/search"
	"userstyles.world/services/cron"
	"userstyles.world/utils"
)

func main() {
	utils.InitalizeCrypto()
	utils.InitializeValidator()
	database.Initialize()
	cron.Initialize()
	search.Initialize()
	images.Initialize()
	handlers.Initialize()
}
