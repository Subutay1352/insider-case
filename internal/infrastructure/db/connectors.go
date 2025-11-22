package db

import (
	"fmt"
	"insider-case/internal/config"
	"insider-case/internal/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PostgresConnector struct {
	cfg *config.DatabaseConfig
}

func NewPostgresConnector(cfg *config.DatabaseConfig) Connector {
	return &PostgresConnector{cfg: cfg}
}

func (c *PostgresConnector) Connect() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.cfg.GetDSN()), gormConfig(c.cfg))
}

type SQLiteConnector struct {
	cfg *config.DatabaseConfig
}

func NewSQLiteConnector(cfg *config.DatabaseConfig) Connector {
	return &SQLiteConnector{cfg: cfg}
}

func (c *SQLiteConnector) Connect() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(c.cfg.Path), gormConfig(c.cfg))
}

type unsupportedConnector struct {
	dbType string
}

func (c *unsupportedConnector) Connect() (*gorm.DB, error) {
	logger.Error("Unsupported database type", "db_type", c.dbType)
	return nil, fmt.Errorf("unsupported database type: %s", c.dbType)
}
