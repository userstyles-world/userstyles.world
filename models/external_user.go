package models

import (
	"time"

	"gorm.io/gorm"
)

// ExternalUser contains external user information from OAuth providers.
type ExternalUser struct {
	ID          uint `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	ExternalID  string         `gorm:"index:idx_external_users_user_id,unique"`
	UserID      uint           `gorm:"index:idx_external_users_user_id,unique"`
	User        User
	Provider    string
	Email       string `gorm:"email;type:TEXT COLLATE NOCASE"`
	Username    string `gorm:"username;type:TEXT COLLATE NOCASE"`
	ExternalURL string
	AccessToken string
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
