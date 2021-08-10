package models

import (
	"gorm.io/gorm"
)

type Kind int

const (
	KindReview Kind = iota
	KindStylePromotion
)

type Notification struct {
	gorm.Model
	Seen     bool
	Kind     Kind
	TargetID int

	User   User
	UserID int

	Style   Style
	StyleID int

	Review   Review
	ReviewID int `gorm:"default:null"`
}

func (n Notification) Create() error {
	return db().Create(&n).Error
}
