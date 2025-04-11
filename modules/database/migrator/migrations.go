package migrator

import (
	"strings"

	"gorm.io/gorm"

	"userstyles.world/models"
	"userstyles.world/modules/config"
	"userstyles.world/modules/log"
)

// deindent removes unneeded control characters from string literals.
func deindent(s string) string {
	s = strings.TrimPrefix(s, "\n")
	s = strings.ReplaceAll(s, "\t\t", "    ")
	s = strings.ReplaceAll(s, "\t", "")
	return s
}

func m1(db *gorm.DB) error {
	const q = `
	CREATE TABLE migrations(
		version    INTEGER,
		name       TEXT,
		applied_at DATETIME
	);`
	return db.Exec(deindent(q)).Error
}

func m2(db *gorm.DB) error {
	cx := []string{"slug", "code_size", "code_checksum"}

	var s models.Style
	for _, c := range cx {
		if !db.Migrator().HasColumn(s, c) {
			if err := db.Migrator().AddColumn(s, c); err != nil {
				return err
			}
			log.Database.Printf("Added %q column to styles table.\n", c)
		}
	}

	var sx []models.Style
	return db.FindInBatches(&sx, 500, func(tx *gorm.DB, size int) error {
		for i, s := range sx {
			s.Prepare()
			if err := tx.Select(cx).UpdateColumns(s).Error; err != nil {
				return err
			}

			// Print everything in debug mode or every 100th message.
			if config.App.Debug || i%100 == 0 {
				const f = "id=%-8d cs=%-8d cc=%-8s s=%.40s\n"
				log.Database.Printf(f, s.ID, s.CodeSize, s.CodeChecksum, s.Slug)
			}
		}

		return nil
	}).Error
}

func m3(db *gorm.DB) error {
	const q = `
	CREATE TABLE changelogs(
		id          INTEGER,
		user_id     INTEGER,
		created_at  DATETIME,
		updated_at  DATETIME,
		deleted_at  DATETIME,
		title       TEXT,
		description TEXT,
		PRIMARY KEY (id),
		CONSTRAINT fk_changelogs_user FOREIGN KEY (user_id) REFERENCES users(id)
	);
	CREATE INDEX idx_changelogs_deleted_at ON changelogs(deleted_at);`
	return db.Exec(deindent(q)).Error
}

func m4(db *gorm.DB) error {
	const q = `
	UPDATE reviews
	SET deleted_at = 'm4 ' || DATETIME('now')
	WHERE style_id IN(
		SELECT id FROM styles WHERE deleted_at IS NOT NULL
	);`
	return db.Exec(deindent(q)).Error
}

func m5(db *gorm.DB) error {
	const q = `
	UPDATE notifications
	SET deleted_at = 'm5 ' || DATETIME('now')
	WHERE kind = 0 AND style_id IN(
		SELECT id FROM styles WHERE deleted_at IS NOT NULL
	);`
	return db.Exec(deindent(q)).Error
}

func m6(db *gorm.DB) error {
	var err error
	const q = `
	ALTER TABLE logs RENAME TO logs_old;
	DROP INDEX idx_logs_deleted_at;

	CREATE TABLE logs(
		id         INTEGER,
		by_user_id INTEGER,
		to_user_id INTEGER,
		style_id   INTEGER,
		review_id  INTEGER,
		kind       INTEGER,
		created_at DATETIME,
		updated_at DATETIME,
		deleted_at DATETIME,
		reason     TEXT,
		message    TEXT,
		censor     NUMERIC,
		PRIMARY KEY (id),
		CONSTRAINT fk_logs_by_user FOREIGN KEY (by_user_id) REFERENCES users(id),
		CONSTRAINT fk_logs_to_user FOREIGN KEY (to_user_id) REFERENCES users(id),
		CONSTRAINT fk_logs_style   FOREIGN KEY (style_id)   REFERENCES styles(id),
		CONSTRAINT fk_logs_review  FOREIGN KEY (review_id)  REFERENCES reviews(id)
	);

	INSERT INTO logs(
		id, created_at, updated_at, deleted_at,
		by_user_id, reason, message, censor, kind
	)
	SELECT
		id, created_at, updated_at, deleted_at,
		user_id, reason, message, censor, kind
	FROM logs_old;

	CREATE INDEX idx_logs_deleted_at ON logs(deleted_at);`
	if err = db.Exec(deindent(q)).Error; err != nil {
		return err
	}

	{
		var users []struct {
			LogID    int    `gorm:"column:id"`
			Username string `gorm:"column:target_user_name"`
		}
		if err = db.Table("logs_old").Find(&users, "kind = 1").Error; err != nil {
			return err
		}

		for _, u := range users {
			var uid int
			const q = "SELECT id FROM users WHERE username = ?"
			if err = db.Raw(q, u.Username).Scan(&uid).Error; err != nil {
				return err
			}

			err = db.
				Table("logs").
				Where("id = ?", u.LogID).
				UpdateColumn("to_user_id", uid).
				Error
			if err != nil {
				return err
			}
		}
	}

	{
		var styles []struct {
			LogID    int    `gorm:"column:id"`
			Name     string `gorm:"column:target_data"`
			Username string `gorm:"column:target_user_name"`
		}
		if err = db.Table("logs_old").Find(&styles, "kind = 2").Error; err != nil {
			return err
		}

		for _, s := range styles {
			var uid int
			q := "SELECT id FROM users WHERE username = ?"
			if err = db.Raw(q, s.Username).Scan(&uid).Error; err != nil {
				return err
			}

			var sid int
			q = "SELECT id FROM styles WHERE name = ? AND user_id = ?"
			if err = db.Raw(q, s.Name, uid).Scan(&sid).Error; err != nil {
				return err
			}

			err = db.
				Table("logs").
				Where("id = ?", s.LogID).
				UpdateColumns(map[string]any{
					"to_user_id": uid,
					"style_id":   sid,
				}).
				Error
			if err != nil {
				return err
			}
		}
	}

	{
		var reviews []struct {
			LogID    int    `gorm:"column:id"`
			ReviewID int    `gorm:"column:target_data"`
			Username string `gorm:"column:target_user_name"`
		}
		if err = db.Table("logs_old").Find(&reviews, "kind = 3").Error; err != nil {
			return err
		}

		for _, r := range reviews {
			var uid int
			q := "SELECT id FROM users WHERE username = ?"
			if err = db.Raw(q, r.Username).Scan(&uid).Error; err != nil {
				return err
			}

			err = db.
				Table("logs").
				Where("id = ?", r.LogID).
				UpdateColumns(map[string]any{
					"to_user_id": uid,
					"review_id":  r.ReviewID,
				}).
				Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
