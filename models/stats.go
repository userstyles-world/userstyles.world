package models

import (
	"errors"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Stats struct {
	gorm.Model
	Hash    string `gorm:"unique"`
	StyleID int
	Style   Style
}

func AddStatsForStyle(db *gorm.DB, id, ip string) (Stats, error) {
	s := new(Stats)

	styleID, err := strconv.Atoi(id)
	if err != nil {
		return *s, errors.New("failed to convert StyleID to int")
	}

	s.Hash = ip + " " + id
	s.StyleID = styleID

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

	if err != nil {
		log.Fatal("Got error:", err)
	}

	return *s, nil
}
