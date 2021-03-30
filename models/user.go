package models

import (
	"errors"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	Username   string `gorm:"unique;not null" validate:"required,username,min=5,max=20"`
	Email      string `gorm:"unique;not null" validate:"required,email"`
	Password   string `gorm:"not null"        validate:"required,min=8,max=32"`
	Biography  string `validate:"min=0,max=512"`
}

type APIUser struct {
	Username  string
	Email     string
	ID        uint
	Biography string
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
