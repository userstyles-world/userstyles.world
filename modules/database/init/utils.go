package init

import (
	"gorm.io/gorm/logger"

	"userstyles.world/config"
	"userstyles.world/modules/database"
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

func migrate(t ...interface{}) error {
	return database.Conn.AutoMigrate(t...)
}

func drop(t ...interface{}) error {
	return database.Conn.Migrator().DropTable(t...)
}
