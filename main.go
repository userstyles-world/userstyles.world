package main

import (
	"userstyles.world/database"
	"userstyles.world/handlers"
	"userstyles.world/utils"
)

func main() {
	utils.InitializeValidator()
	database.Initialize()
	handlers.Initialize()
}
