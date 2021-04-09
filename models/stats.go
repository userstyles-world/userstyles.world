package models

import (
	"errors"
	"log"
	"strconv"

	"gorm.io/gorm"
)

type Stats struct {
	gorm.Model
	IP      string
	StyleID int
	Style   Style
}

func AddStatsForStyle(db *gorm.DB, id, ip string) (Stats, error) {
	s := new(Stats)

	styleID, err := strconv.Atoi(id)
	if err != nil {
		return *s, errors.New("failed to convert StyleID to int")
	}

	s.IP = ip
	s.StyleID = styleID

	err = db.
		Debug().
		Model(s).
		Create(s).
		Error

	if err != nil {
		log.Fatal("Got error:", err)
	}

	return *s, nil
}
