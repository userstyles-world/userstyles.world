package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	"github.com/ohler55/ojg/oj"
	"gorm.io/gorm"
)

type OAuth struct {
	gorm.Model
	UserID       uint
	User         User
	Name         string     `gorm:"unique;not null" validate:"required,min=1,max=256"`
	Description  string     `validate:"required,min=0,max=1028"`
	Scopes       StringList `gorm:"type:varchar(255);"`
	RedirectURI  string
	ClientID     string
	ClientSecret string
}

type APIOAuth struct {
	ID           uint
	Name         string `gorm:"unique"`
	Description  string
	Scopes       StringList `gorm:"type:varchar(255);"`
	RedirectURI  string
	UserID       uint
	Username     string
	ClientID     string
	ClientSecret string
}

// Custom []string time for the GORM.
// As gorm highly dislike slices, we have to implement, this ourself.
type StringList []string

func (s StringList) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return fmt.Sprintf(`["%s"]`, strings.Join(s, `","`)), nil
}

func (s *StringList) Scan(src interface{}) (err error) {
	var scopes []string
	switch src := src.(type) {
	case string:
		err = oj.Unmarshal([]byte(src), &scopes)
	case []byte:
		err = oj.Unmarshal(src, &scopes)
	default:
		return errors.New("incompatible type for Scopes")
	}
	if err != nil {
		return
	}
	*s = scopes
	return nil
}

// Override table name.
func (OAuth) TableName() string {
	return "oauths"
}

func ListOAuthsOfUser(db *gorm.DB, username string) (*[]APIOAuth, error) {
	t, q := new(OAuth), new([]APIOAuth)
	err := getDBSession(db).
		Model(t).
		Select("oauths.id, oauths.name, u.username").
		Joins("join users u on u.id = oauths.user_id").
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
		Select("oauths.*,  u.username").
		Joins("join users u on u.id = oauths.user_id").
		First(q, "oauths.id = ?", id).
		Error

	if err != nil || q.ID == 0 {
		return nil, errors.New("oauth not found")
	}

	return q, nil
}

// Using ID as a string is fine in this case.
func GetOAuthByClientID(db *gorm.DB, clientID string) (*APIOAuth, error) {
	t, q := new(OAuth), new(APIOAuth)
	err := getDBSession(db).
		Debug().
		Model(t).
		Select("oauths.*,  u.username").
		Joins("join users u on u.id = oauths.user_id").
		First(q, "oauths.client_id = ?", clientID).
		Error

	if err != nil || q.ID == 0 {
		return nil, errors.New("oauth not found")
	}

	return q, nil
}

func CreateOAuth(db *gorm.DB, o *OAuth) (*OAuth, error) {
	err := getDBSession(db).
		Create(&o).
		Error
	if err != nil {
		return o, err
	}

	return o, nil
}

func UpdateOAuth(db *gorm.DB, o *OAuth) error {
	err := getDBSession(db).
		Debug().
		Model(OAuth{}).
		Where("id", o.ID).
		Updates(o).
		Error
	if err != nil {
		return err
	}

	return nil
}
