package database

import (
	"gorm.io/gorm"
)

// Conn holds the pointer to database connection.
var Conn *gorm.DB
