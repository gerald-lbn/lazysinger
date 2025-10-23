package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gerald-lbn/lazysinger/config"
	"github.com/gerald-lbn/lazysinger/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a new database connection
func Connect() *gorm.DB {
	// Create database directory if it doesn't exist
	dbDir := filepath.Join(config.Server.General.DataPath)
	if err := ensureDir(dbDir); err != nil {
		log.Fatal().Err(err).Msg("Failed to create database directory")
		return nil
	}

	dbPath := filepath.Join(dbDir, config.Server.General.DatabaseName)
	log.Info().Str("path", dbPath).Msg("Opening database")

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel()),
	}

	db, err := gorm.Open(sqlite.Open(dbPath), gormConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open database")
		return nil
	}

	// Enable foreign key constraints
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		log.Fatal().Err(err).Msg("Failed to enable foreign key constraints")
		return nil
	}

	// Enable WAL mode for better concurrent access
	if err := db.Exec("PRAGMA journal_mode = WAL").Error; err != nil {
		log.Fatal().Err(err).Msg("Failed to enable WAL mode")
		return nil
	}

	// Set busy timeout to handle concurrent access
	if err := db.Exec("PRAGMA busy_timeout = 5000").Error; err != nil {
		log.Fatal().Err(err).Msg("Failed to set busy timeout")
		return nil
	}

	// Auto-migrate database schema
	if err := AutoMigrate(db); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database schema")
		return nil
	}

	return db
}

// Close closes the database connection
func Close(db *gorm.DB) error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying database instance: %w", err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close database connection: %w", err)
		}
	}
	return nil
}

// ensureDir creates a directory if it doesn't exist
func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// logLevel converts the application log level to GORM logger level
func logLevel() logger.LogLevel {
	switch config.Server.Logger.Level {
	case "debug":
		return logger.Info
	case "info":
		return logger.Warn
	case "warn":
		return logger.Error
	case "error":
		return logger.Silent
	default:
		return logger.Warn
	}
}
