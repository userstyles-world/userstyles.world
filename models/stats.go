package models

import (
	"gorm.io/gorm"
)

type Stats struct {
	gorm.Model
	IP      string
	StyleID int
	Style   Style
}
