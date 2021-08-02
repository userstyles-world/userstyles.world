package models

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/ohler55/ojg/oj"
	"gorm.io/gorm"

	"userstyles.world/modules/errors"
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

var modelOAuth = OAuth{}

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
		return errors.ErrIncompatibleType
	}
	if err != nil {
		return err
	}
	*s = scopes
	return nil
}

// TableName specify the table name that should be used.
func (OAuth) TableName() string {
	return "oauths"
}

func ListOAuthsOfUser(username string) (*[]APIOAuth, error) {
	q := new([]APIOAuth)
	err := getDBSession().
		Model(modelOAuth).
		Select("oauths.id, oauths.name, u.username").
		Joins("join users u on u.id = oauths.user_id").
		Find(q, "u.username = ?", username).
		Error
	if err != nil {
		return nil, errors.ErrOAuthNotFound
	}

	return q, nil
}

// GetOAuthByID note: Using ID as a string is fine in this case.
func GetOAuthByID(id string) (*APIOAuth, error) {
	q := new(APIOAuth)
	err := getDBSession().
		Debug().
		Model(modelOAuth).
		Select("oauths.*,  u.username").
		Joins("join users u on u.id = oauths.user_id").
		First(q, "oauths.id = ?", id).
		Error

	if err != nil || q.ID == 0 {
		return nil, errors.ErrOAuthNotFound
	}

	return q, nil
}

// GetOAuthByClientID note: Using ID as a string is fine in this case.
func GetOAuthByClientID(clientID string) (*APIOAuth, error) {
	q := new(APIOAuth)
	err := getDBSession().
		Debug().
		Model(modelOAuth).
		Select("oauths.*,  u.username").
		Joins("join users u on u.id = oauths.user_id").
		First(q, "oauths.client_id = ?", clientID).
		Error

	if err != nil || q.ID == 0 {
		return nil, errors.ErrOAuthNotFound
	}

	return q, nil
}

func CreateOAuth(o *OAuth) (*OAuth, error) {
	err := getDBSession().
		Debug().
		Create(&o).
		Error
	if err != nil {
		return o, err
	}

	return o, nil
}

func UpdateOAuth(o *OAuth, id string) error {
	err := getDBSession().
		Debug().
		Model(modelOAuth).
		Where("id = ?", id).
		Updates(o).
		Error
	if err != nil {
		return err
	}

	return nil
}
