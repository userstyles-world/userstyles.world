package models

import (
	"gorm.io/gorm"

	"userstyles.world/modules/database"
	"userstyles.world/modules/errors"
)

type History struct {
	gorm.Model
	StyleID       uint
	DailyViews    int64
	DailyInstalls int64
	DailyUpdates  int64
	TotalViews    int64
	TotalInstalls int64
	TotalUpdates  int64
}

func (h History) GetStatsForStyle(id string) (q *[]History, err error) {
	err = database.Conn.
		Debug().
		Model(History{}).
		Where("style_id = ?", id).
		Find(&q).
		Error
	if err != nil {
		return nil, errors.ErrNoStyleStats
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
		return nil, errors.ErrFailedHistoriesSearch
	}

	return q, nil
}
