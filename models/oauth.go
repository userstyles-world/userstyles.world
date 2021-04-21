package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type OAuth struct {
	gorm.Model
	UserID      uint
	User        User `gorm:"foreignKey:ID"`
	Name        string
	Description string
	Scopes      Scopes `gorm:"type:varchar(255);" db:"scopes"`
	RedirectURI string
}

type APIOAuth struct {
	ID          uint
	Name        string
	Description string
	Scopes      Scopes `gorm:"type:varchar(255);" db:"scopes"`
	RedirectURI string
	UserID      uint
	Username    string
}

type Scopes []string

func (s Scopes) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return fmt.Sprintf(`["%s"]`, strings.Join(s, `","`)), nil
}

func (s *Scopes) Scan(src interface{}) (err error) {
	var scopes []string
	switch src := src.(type) {
	case string:
		err = json.Unmarshal([]byte(src), &scopes)
	case []byte:
		err = json.Unmarshal(src, &scopes)
	default:
		return errors.New("incompatible type for Scopes")
	}
	if err != nil {
		return
	}
	*s = scopes
	return nil
}

func ListOAuthsOfUser(db *gorm.DB, username string) (*[]APIOAuth, error) {
	t, q := new(OAuth), new([]APIOAuth)
	err := getDBSession(db).
		Model(t).
		Select("o_auths.id, o_auths.name, u.username").
		Joins("join users u on u.id = o_auths.user_id").
		Find(q, "u.username = ?", username).
		Error

	if err != nil {
		return nil, errors.New("oauth not found")
	}

	return q, nil
}

// Using ID as a string is fine in this case.
func GetOAuthByID(db *gorm.DB, id string) (*APIOAuth, error) {
	t, q := new(OAuth), new(APIOAuth)
	err := getDBSession(db).
		Debug().
		Model(t).
		Select("o_auths.*,  u.username").
		Joins("join users u on u.id = o_auths.user_id").
		First(q, "o_auths.id = ?", id).
		Error

	if err != nil || q.ID == 0 {
		return nil, errors.New("oauth not found")
	}

	return q, nil
}
