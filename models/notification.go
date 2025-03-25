package models

import (
	"gorm.io/gorm"
)

type Kind int

const (
	KindReview Kind = iota
	KindStylePromotion
	KindBannedStyle
	KindRemovedReview
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

// CreateNotification inserts a new notification.
func CreateNotification(db *gorm.DB, n *Notification) error {
	return db.Create(&n).Error
}

// DeleteNotificationsForStyle tries to delete all notifications for a style.
func DeleteNotificationsForStyle(db *gorm.DB, id int) (err error) {
	return db.Delete(&Notification{}, "style_id = ?", id).Error
}
