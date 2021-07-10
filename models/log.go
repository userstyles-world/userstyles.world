package models

import (
	"gorm.io/gorm"

	"userstyles.world/modules/database"
	"userstyles.world/modules/errors"
)

type LogKind = uint8

const (
	LogBanUser LogKind = iota + 1
	LogRemoveStyle
)

// Log struct has all the relavant information for a log entry
type Log struct {
	gorm.Model
	UserID         uint
	Username       string
	Reason         string
	Kind           LogKind
	TargetData     string
	TargetUserName string
}

// AddLog adds a new log to the database.
func (l *Log) AddLog(log Log) (err error) {
	err = database.Conn.
		Debug().
		Model(Log{}).
		Create(&log).
		Error
	if err != nil {
		return errors.ErrFailedLogAddition
	}

	return nil
}

// GetLogOfKind returns all the logs of the specified kind and
// select the correct user Author.
func GetLogOfKind(kind LogKind) (q *[]Log, err error) {
	err = database.Conn.
		Debug().
		Model(Log{}).
		Select("logs.*, u.id, u.username").
		Joins("join users u on u.id = logs.user_id").
		Where("kind = ?", kind).
		Find(&q).
		Error
	if err != nil {
		return nil, errors.ErrFailedLogRetrieval
	}
	return q, nil
}
