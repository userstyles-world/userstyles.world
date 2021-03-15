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
	default:
		return logger.Silent
	}
}

func DatabaseColorful(db *gorm.DB) bool {
	switch config.DB_COLOR {
	case "true", "yes", "1":
		return true
	default:
		return false
	}
}

func DatabaseDropTables(db *gorm.DB) bool {
	switch config.DB_DROP {
	case "true", "yes", "1":
		return true
	default:
		return false
	}
}
