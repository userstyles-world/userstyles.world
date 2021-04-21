package models

import (
	"errors"

	"gorm.io/gorm"
)

type OAuth struct {
	gorm.Model
	UserID      uint
	User        User `gorm:"foreignKey:ID"`
	Name        string
	Description string
	Scopes      []string
	RedirectURI string
}

type APIOAuth struct {
	ID          uint
	Name        string
	Description string
	Scopes      []string
	RedirectURI []string
	UserID      uint
	Username    string
}

func GetOAuthByUser(db *gorm.DB, username string) (*[]APIOAuth, error) {
	t, q := new(OAuth), new([]APIOAuth)
	err := getDBSession(db).
		Model(t).
		Select("oauth.id, oauth.name, u.username").
		Joins("join users u on u.id = oauth.user_id").
		Find(q, "u.username = ?", username).
		Error

	if err != nil {
		return nil, errors.New("OAuth not found.")
	}

	return q, nil
}

// Using ID as a string is fine in this case.
func GetOAuthByID(db *gorm.DB, id string) (*APIOAuth, error) {
	t, q := new(OAuth), new(APIOAuth)
	err := getDBSession(db).
		Model(t).
		Select("oauth.*,  u.username").
		Joins("join users u on u.id = oauth.user_id").
		Find(q, "oauth.id = ?", id).
		Error

	if err != nil || q.ID == 0 {
		return nil, errors.New("OAuth not found.")
	}

	return q, nil
}
