package init

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"userstyles.world/models"
	"userstyles.world/modules/log"
)

func runMigration(db *gorm.DB) {
	log.Database.Println("Migration started.")
	t := time.Now()

	db.Config.Logger = db.Config.Logger.LogMode(logger.Info)

	// Wrap in a transaction to allow rollbacks.
	db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("PRAGMA foreign_keys = OFF")

		var u models.User
		var err error

		if err = tx.Migrator().DropIndex(u, "idx_users_deleted_at"); err != nil {
			log.Database.Fatalf("Failed to drop index for users: %s\n", err)
		}

		if err = tx.Migrator().RenameTable("users", "users_temp"); err != nil {
			log.Database.Fatalf("Failed to rename table users: %s\n", err)
		}

		if err = tx.Migrator().CreateTable(u); err != nil {
			log.Database.Fatalf("Failed to create table users: %s\n", err)
		}

		q := `INSERT INTO users SELECT id, created_at, updated_at, deleted_at, `
		q += `username, email, o_auth_provider, password, biography, `
		q += `display_name, role, last_login, last_password_reset, `
		q += `authorized_o_auth, github, gitlab, codeberg FROM users_temp`
		if err = tx.Exec(q).Error; err != nil {
			log.Database.Fatalf("Failed to insert users: %s\n", err)
		}

		if err = tx.Migrator().DropTable("users_temp"); err != nil {
			log.Database.Fatalf("Failed to drop users_temp table: %s\n", err)
		}

		// HACK: Fix up foreign key references in tables that have user_id.
		if err = tx.Migrator().RenameTable("users", "users_temp"); err != nil {
			log.Database.Fatalf("Failed to rename table users: %s\n", err)
		}
		if err = tx.Migrator().RenameTable("users_temp", "users"); err != nil {
			log.Database.Fatalf("Failed to rename table users: %s\n", err)
		}

		var eu models.ExternalUser
		if err = tx.AutoMigrate(&eu); err != nil {
			log.Database.Fatalf("Failed to create external_users: %s\n", err)
		}

		return nil
	})

	log.Database.Printf("Migration completed in %s.\n", time.Since(t))
}
