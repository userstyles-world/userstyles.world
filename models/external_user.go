package models

import (
	"time"

	"gorm.io/gorm"
)

// ExternalUser contains external user information from OAuth providers.
type ExternalUser struct {
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ExternalID  string         `gorm:"primaryKey"`
	UserID      uint           `gorm:"primaryKey"`
	User        User
	Provider    string
	Email       string
	Username    string
	ExternalURL string
	AccessToken string
	RawData     string
}

// TableName returns a table to be used with GORM.
func (u *ExternalUser) TableName() string { return "external_users" }

// NormalizeUsername overwrites the username field for services—like GitHub—that
// use "LoginName" instead of "UserName" in their user information responses.
func (u *ExternalUser) NormalizeUsername(username string) {
	if username != "" {
		u.Username = username
	}
}
