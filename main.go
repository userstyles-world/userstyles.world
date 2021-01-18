package main

import (
	"userstyles.world/database"
	"userstyles.world/handlers"
)

func main() {
	database.Connect()
	database.Prepare()
	handlers.Initialize()

}
