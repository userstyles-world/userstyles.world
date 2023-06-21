package init

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"userstyles.world/modules/log"
)

func runMigration(db *gorm.DB) {
	log.Database.Println("Migration started.")
	t := time.Now()

	db.Config.Logger = db.Config.Logger.LogMode(logger.Info)

	// Wrap in a transaction to allow rollbacks.
	db.Transaction(func(tx *gorm.DB) error {
		bootstrap := `DROP TABLE IF EXISTS fts_styles;
CREATE VIRTUAL TABLE fts_styles USING FTS5(id, name, description, notes, tokenize="trigram");
INSERT INTO fts_styles(id, name, description, notes) SELECT id, name, description, notes FROM styles;
`
		if err := tx.Exec(bootstrap).Error; err != nil {
			log.Database.Fatalf("Failed to run bootstrap: %s\n", err)
		}

		triggers := `DROP TRIGGER IF EXISTS fts_styles_insert;
CREATE TRIGGER fts_styles_insert AFTER INSERT ON styles
BEGIN
	INSERT INTO fts_styles(id, name, description, notes)
	VALUES (new.id, new.name, new.description, new.notes);
END;

DROP TRIGGER IF EXISTS fts_styles_update;
CREATE TRIGGER fts_styles_update AFTER UPDATE ON styles
BEGIN
	UPDATE fts_styles
	SET name = new.name, description = new.description, notes = new.notes
	WHERE id = old.id;
END;

DROP TRIGGER IF EXISTS fts_styles_delete;
CREATE TRIGGER fts_styles_delete AFTER DELETE ON styles
BEGIN
	DELETE FROM fts_styles WHERE id = old.id;
END;
`
		if err := tx.Exec(triggers).Error; err != nil {
			log.Database.Fatalf("Failed to add triggers: %s\n", err)
		}

		return nil
	})

	log.Database.Printf("Done in %s.\n", time.Since(t).Round(time.Microsecond))
}
