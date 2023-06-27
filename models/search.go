package models

import "userstyles.world/modules/database"

func InitStyleSearch() error {
	init := `
DROP TABLE IF EXISTS fts_styles;
CREATE VIRTUAL TABLE fts_styles USING FTS5(id, name, description, notes, tokenize="trigram");
INSERT INTO fts_styles(id, name, description, notes) SELECT id, name, description, notes FROM styles;

DROP TRIGGER IF EXISTS fts_styles_insert;
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
	return database.Conn.Exec(init).Error
}
