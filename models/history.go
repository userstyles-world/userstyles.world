package models

import (
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	StyleID       uint
	DailyViews    int64
	DailyInstalls int64
	DailyUpdates  int64
}
