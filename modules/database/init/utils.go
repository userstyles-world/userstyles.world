package init

import (
	"gorm.io/gorm/logger"

	"userstyles.world/config"
)

func logLevel() logger.LogLevel {
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

func colorful() bool {
	switch config.DB_COLOR {
	case "true", "yes", "1":
		return true
	default:
		return false
	}
}

func dropTables() bool {
	switch config.DB_DROP {
	case "true", "yes", "1":
		return true
	default:
		return false
	}
}
