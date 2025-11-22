package db

import (
	"fmt"
	"insider-case/internal/config"
	"insider-case/internal/pkg/logger"

	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type Connector interface {
	Connect() (*gorm.DB, error)
}

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	connector := NewConnector(&cfg.Database)

	db, err := connector.Connect()
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.Database.ConnMaxIdleTime)

	logger.Info("Database connection pool configured",
		"max_open_conns", cfg.Database.MaxOpenConns,
		"max_idle_conns", cfg.Database.MaxIdleConns,
		"conn_max_lifetime", cfg.Database.ConnMaxLifetime,
		"conn_max_idle_time", cfg.Database.ConnMaxIdleTime,
	)

	return db, nil
}

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := ConnectDB(cfg)
	if err != nil {
		return nil, err
	}

	if err := RunMigrations(db); err != nil {
		logger.Error("Failed to run database migrations", "error", err)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func NewConnector(dbCfg *config.DatabaseConfig) Connector {
	dbType := dbCfg.Type
	if dbType == "" {
		dbType = "postgres"
	}

	switch dbType {
	case "postgres":
		return NewPostgresConnector(dbCfg)
	case "sqlite":
		return NewSQLiteConnector(dbCfg)
	default:
		return &unsupportedConnector{dbType: dbType}
	}
}

// gormConfig returns GORM configuration based on config
func gormConfig(cfg *config.DatabaseConfig) *gorm.Config {
	var logLevel gormLogger.LogLevel
	switch cfg.LogLevel {
	case "silent":
		logLevel = gormLogger.Silent
	case "error":
		logLevel = gormLogger.Error
	case "warn":
		logLevel = gormLogger.Warn
	case "info":
		logLevel = gormLogger.Info
	default:
		logLevel = gormLogger.Info
	}

	return &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLevel),
	}
}
