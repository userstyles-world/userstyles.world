package models

import (
	"gorm.io/gorm"
)

type History struct {
	gorm.Model
	StyleID       uint
	DailyViews    int
	DailyInstalls int
	DailyUpdates  int
}
