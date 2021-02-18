package utils

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"userstyles.world/config"
)

func DatabaseLogLevel(db *gorm.DB) logger.LogLevel {
	switch config.DB_DEBUG {
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	case "silent":
		return logger.Silent
	default:
		return logger.Silent
	}
}
