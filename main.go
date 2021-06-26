package main

import (
	"userstyles.world/handlers"
	database "userstyles.world/modules/database/init"
	"userstyles.world/modules/images"
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
