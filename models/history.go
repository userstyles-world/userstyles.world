package models

import (
	"errors"

	"gorm.io/gorm"

	"userstyles.world/modules/database"
)

type History struct {
	gorm.Model
	StyleID       uint
	DailyViews    int64
	DailyInstalls int64
	DailyUpdates  int64
}

func (h History) GetStatsForStyle(id string) (q *[]History, err error) {
	err = database.Conn.
		Debug().
		Model(History{}).
		Where("style_id = ?", id).
		Find(&q).
		Error
	if err != nil {
		return nil, errors.New("style doesn't have stats yet")
	}

	return q, nil
}

func (h History) GetStatsForAllStyles() (q *[]History, err error) {
	err = database.Conn.
		Debug().
		Model(History{}).
		Find(&q).
		Error
	if err != nil {
		return nil, errors.New("failed to find all style histories")
	}

	return q, nil
}
