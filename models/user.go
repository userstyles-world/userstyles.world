package models

import (
	"strings"

	"gorm.io/gorm"
	"userstyles.world/errors_helper"
)

type Role int

const (
	Regular Role = iota
	Moderator
	Admin
)

type SocialMedia struct {
	Github   string
	Gitlab   string
	Codeberg string
}

type User struct {
	gorm.Model    `json:"-"`
	Socials       SocialMedia `gorm:"embedded"`
	Username      string      `gorm:"unique;not null" validate:"required,username,min=5,max=20"`
	Email         string      `gorm:"unique" validate:"required,email"`
	OAuthProvider string      `gorm:"default:none"`
	Password      string      `validate:"required,min=8,max=32"`
	Biography     string      `validate:"min=0,max=512"`
	DisplayName   string      `validate:"displayName,min=5,max=20"`
	Role          Role        `gorm:"default=0"` // The values within SocialMedia struct
	// Will be saved under the user struct
}

type APIUser struct {
	Username    string
	DisplayName string
	Email       string
	Biography   string
	ID          uint
	Role        Role
}

// Check if user set any social media.
func (u User) HasSocials() bool {
	return u.Socials.Codeberg != "" ||
		u.Socials.Gitlab != "" ||
		u.Socials.Github != ""
}

// Return display name if it is set.
func (u User) Name() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}

	return u.Username
}

// Return user's role in string format.
func (u User) RoleString() (s string) {
	switch u.Role {
	case Regular:
		s = "Regular"
	case Moderator:
		s = "Moderator"
	case Admin:
		s = "Admin"
	}
	return s
}

func FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	user := new(User)

	if res := db.Where("email = ?", email).First(&user); res.Error != nil {
		return nil, res.Error
	}

	if user.ID == 0 {
		return nil, errors_helper.ErrUserNotFound
	}

	return user, nil
}

func FindUserByName(db *gorm.DB, name string) (*User, error) {
	user := new(User)

	err := getDBSession(db).
		Where("username = ?", strings.ToLower(name)).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors_helper.ErrUserNotFound
	}

	return user, nil
}
