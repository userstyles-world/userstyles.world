package models

import (
	"gorm.io/gorm"

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
	err = db().
		Model(modelHistory).
		Where("style_id = ?", id).
		Find(&q).
		Error
	if err != nil {
		return nil, errors.ErrNoStyleStats
	}

	return q, nil
}

func (h History) GetStatsForAllStyles() (q *[]History, err error) {
	stmt := "sum(daily_installs) DailyInstalls, sum(daily_views) DailyViews, sum(daily_updates) DailyUpdates, "
	stmt += "sum(total_installs) TotalInstalls, sum(total_views) TotalViews, created_at"

	err = db().
		Select(stmt).
		Group("date(histories.created_at)").
		Find(&q).
		Error
	if err != nil {
		return nil, errors.ErrFailedHistoriesSearch
	}

	return q, nil
}
