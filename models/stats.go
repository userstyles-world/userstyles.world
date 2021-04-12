package models

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"userstyles.world/config"
)

var (
	lastWeek = time.Now().AddDate(0, 0, -7)
)

type Stats struct {
	gorm.Model
	Hash    string `gorm:"unique"`
	Install bool   `gorm:"default:false"`
	StyleID int
	Style   Style
}

func AddStatsToStyle(db *gorm.DB, id, ip string, install bool) (Stats, error) {
	s := new(Stats)

	styleID, err := strconv.Atoi(id)
	if err != nil {
		return *s, err
	}

	// TODO: Refactor as GenerateHashedRecord; we have a circular dependency now.
	record := ip + " " + id

	// Generate unique hash.
	h := hmac.New(sha512.New, []byte(config.STATS_KEY))
	h.Write([]byte(record))
	sha := hex.EncodeToString(h.Sum(nil))

	// Set values.
	s.Hash = sha
	s.StyleID = styleID

	if install {
		err = db.
			Debug().
			Model(s).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "hash"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"updated_at": time.Now(),
					"install":    true,
				}),
			}).
			Create(s).
			Error
	} else {
		err = db.
			Debug().
			Model(s).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "hash"}},
				DoUpdates: clause.Assignments(map[string]interface{}{
					"updated_at": time.Now(),
				}),
			}).
			Create(s).
			Error
	}

	if err != nil {
		return *s, err
	}

	return *s, nil
}

func GetWeeklyInstallsForStyle(db *gorm.DB, id string) int64 {
	var weekly int64
	db.
		Model(Stats{}).
		Where("style_id = ? and install = ? and updated_at > ?", id, true, lastWeek).
		Count(&weekly)

	return weekly
}

func GetTotalInstallsForStyle(db *gorm.DB, id string) int64 {
	var total int64
	db.
		Model(Stats{}).
		Where("style_id = ? and install = ?", id, true).
		Count(&total)

	return total
}

func GetTotalViewsForStyle(db *gorm.DB, id string) int64 {
	var total int64
	db.
		Model(Stats{}).
		Where("style_id = ?", id).
		Count(&total)

	return total
}
