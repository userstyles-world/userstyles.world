package models

import (
	"gorm.io/gorm"

	"userstyles.world/modules/errors"
)

type History struct {
	gorm.Model
	StyleID        uint `gorm:"index"`
	DailyViews     int64
	DailyInstalls  int64
	DailyUpdates   int64
	WeeklyViews    int64
	WeeklyInstalls int64
	WeeklyUpdates  int64
	TotalViews     int64
	TotalInstalls  int64
	TotalUpdates   int64
}

func GetStyleHistory(id string) (h []History, err error) {
	err = db().
		Model(modelHistory).
		Where("style_id = ?", id).
		Find(&h).
		Error
	if err != nil {
		return nil, errors.ErrNoStyleStats
	}

	return h, nil
}

func GetAllStyleHistories() (h []History, err error) {
	stmt := "sum(daily_installs) DailyInstalls, sum(daily_views) DailyViews, sum(daily_updates) DailyUpdates, "
	stmt += "sum(total_installs) TotalInstalls, sum(total_views) TotalViews, created_at"

	err = db().
		Select(stmt).
		Group("date(created_at)").
		Find(&h, "created_at > date('now', '-3 months')").
		Error
	if err != nil {
		return nil, errors.ErrFailedHistoriesSearch
	}

	return h, nil
}
