package models

import (
	"strings"
	"time"

	"gorm.io/gorm"

	"userstyles.world/modules/errors"
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
	gorm.Model `json:"-"`
	// The values within SocialMedia struct
	Socials SocialMedia `gorm:"embedded"`
	// Will be saved under the user struct
	AuthorizedOAuth   StringList `gorm:"type:varchar(255)"`
	Username          string     `gorm:"unique;not null" validate:"required,username,min=3,max=32"`
	Email             string     `gorm:"unique" validate:"required,email"`
	OAuthProvider     string     `gorm:"default:none"`
	Password          string     `validate:"required,min=8,max=32"`
	Biography         string     `validate:"min=0,max=512"`
	DisplayName       string     `validate:"displayName,min=3,max=32"`
	Role              Role       `gorm:"default=0"`
	LastLogin         time.Time
	LastPasswordReset time.Time
}

type APIUser struct {
	Username    string
	DisplayName string
	Email       string
	Biography   string
	ID          uint
	Role        Role
	Scopes      StringList
}

// HasSocials checks if user set any social media.
func (u User) HasSocials() bool {
	return u.Socials.Codeberg != "" ||
		u.Socials.Gitlab != "" ||
		u.Socials.Github != ""
}

// Name Return display name if it is set.
func (u User) Name() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}

	return u.Username
}

// RoleString Return user's role in string format.
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

// TODO: Remove `APIUser` when JWT is swapped to `User` model.
func (u APIUser) IsAdmin() bool {
	return u.Role == Admin
}

// TODO: Remove `APIUser` when JWT is swapped to `User` model.
func (u APIUser) IsModOrAdmin() bool {
	return u.Role == Moderator || u.Role == Admin
}

func FindUserByEmail(email string) (*User, error) {
	user := new(User)

	if res := db().Where("email = ?", email).First(&user); res.Error != nil {
		return nil, res.Error
	}

	if user.ID == 0 {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func FindUserByName(name string) (*User, error) {
	user := new(User)

	err := db().
		Where("username = ?", strings.ToLower(name)).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func FindUserByNameOrEmail(name, email string) (*User, error) {
	user := new(User)

	err := db().
		Where("username = ? or email = ?", strings.ToLower(name), email).
		First(&user).
		Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func FindUserByID(id string, forbid ...string) (*User, error) {
	user := new(User)

	// HACK: Forbid scraping public user info via API if they have no userstyles.
	var hack string
	if len(forbid) > 0 {
		hack = "(SELECT user_id FROM styles WHERE user_id = users.id) AND "
	}

	err := db().Model(modelUser).Where(hack+"id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

func UpdateUser(u *User) error {
	err := db().
		Debug().
		Model(modelUser).
		Where("id", u.ID).
		Updates(u).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (*User) DeleteWhereID(id any) error {
	return db().Delete(&User{}, "id = ?", id).Error
}

func (u *User) UpdateLastLogin() error {
	return db().Model(&u).Where("id", u.ID).
		UpdateColumn("last_login", time.Now()).Error
}

func (u *User) UpdateLastPasswordRequest() error {
	return db().Model(&u).Where("id", u.ID).
		UpdateColumn("last_password_reset", time.Now()).Error
}
