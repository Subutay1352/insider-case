package db

import (
	"insider-case/internal/domain/message"
	"insider-case/internal/pkg/logger"

	"gorm.io/gorm"
)

// RunMigrations runs database migrations
func RunMigrations(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	// Auto-migrate all models
	if err := db.AutoMigrate(&message.Message{}); err != nil {
		return err
	}

	logger.Info("Database migrations completed successfully")
	return nil
}
