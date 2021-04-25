package models

import (
	"errors"

	"gorm.io/gorm"
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
	Username      string `gorm:"unique;not null" validate:"required,username,min=5,max=20"`
	Email         string `gorm:"unique" validate:"required,email"`
	OAuthProvider string `gorm:"default:none"`
	Password      string `validate:"required,min=8,max=32"`
	Biography     string `validate:"min=0,max=512"`
	Role          Role   `gorm:"default=0"`

	// The values within SocialMedia struct
	// Will be saved under the user struct
	Socials         SocialMedia `gorm:"embedded"`
	AuthorizedOAuth StringList  `gorm:"default=[];type:varchar(255)"`
}

type APIUser struct {
	Username  string
	Email     string
	ID        uint
	Biography string
	Role      Role
	Scopes    StringList
}

func FindUserByEmail(db *gorm.DB, email string) (*User, error) {
	user := new(User)

	if res := db.Where("email = ?", email).First(&user); res.Error != nil {
		return nil, res.Error
	}

	if user.ID == 0 {
		return nil, errors.New("User not found.")
	}

	return user, nil
}

func FindUserByName(db *gorm.DB, name string) (*User, error) {
	user := new(User)

	err := getDBSession(db).
		Where("username = ?", name).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("User not found.")
	}

	return user, nil
}

func UpdateUser(db *gorm.DB, u *User) error {
	err := getDBSession(db).
		Debug().
		Model(User{}).
		Where("id", u.ID).
		Updates(u).
		Error

	if err != nil {
		return err
	}

	return nil

}
