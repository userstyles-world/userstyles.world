package init

import (
	"userstyles.world/models"
	"userstyles.world/modules/database"
	"userstyles.world/modules/log"
)

func migration() {
	log.Database.Println("Migration started.")

	db := database.Conn

	var eu models.ExternalUser
	if err := db.AutoMigrate(&eu); err != nil {
		log.Database.Fatalf("Failed to create %q table.\n", eu.TableName())
	}

	log.Database.Println("Migration completed.")
}
