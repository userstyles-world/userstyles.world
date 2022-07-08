package storage

import (
	"time"

	"userstyles.world/modules/database"
)

// User is a field-aligned struct optimized for users.
type User struct {
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Username      string    `json:"username"`
	DisplayName   string    `json:"display_name"`
	Email         string    `json:"email"`
	OAuthProvider string    `json:"oauth"`
	ID            int       `json:"id"`
}

// FindUsersCreatedOn returns style cards made on a specific date.
func FindUsersCreatedOn(date time.Time) ([]User, error) {
	var res []User

	cond := notDeleted + " AND DATE(?) == DATE(created_at)"
	err := database.Conn.Find(&res, cond, date).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}
