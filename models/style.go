package models

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

type Style struct {
	gorm.Model
	UserID      uint
	User        User `gorm:"foreignKey:ID"`
	Name        string
	Summary     string
	Description string
	Code        string
	Preview     string
	Archived    bool   `gorm:"default:false"`
	Featured    bool   `gorm:"default:false"`
	Category    string `gorm:"not null"`
}

type APIStyle struct {
	ID          uint
	Name        string
	Summary     string
	Description string
	Code        string
	Preview     string
	Archived    bool
	Featured    bool
	Category    string
	UserID      uint
	Username    string
}
