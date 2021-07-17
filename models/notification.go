package models

import (
	"gorm.io/gorm"

	"userstyles.world/modules/database"
)

type Kind int

const (
	KindReview Kind = iota
)

type Notification struct {
	gorm.Model
	Seen bool
	Kind Kind

	User   User
	UserID int

	Style   Style
	StyleID int

	Review   Review
	ReviewID int
}

func (n Notification) Create() error {
	return database.Conn.Debug().Create(&n).Error
}
