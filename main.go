package main

import (
	"userstyles.world/database"
	"userstyles.world/handlers"
)

func main() {
	database.Initialize()
	handlers.Initialize()
}
