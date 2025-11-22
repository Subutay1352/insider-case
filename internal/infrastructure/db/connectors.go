package db

import (
	"fmt"
	"insider-case/internal/config"
	"insider-case/internal/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// PostgresConnector implements Connector for PostgreSQL
type PostgresConnector struct {
	cfg *config.DatabaseConfig
}

// NewPostgresConnector creates a new PostgreSQL connector
func NewPostgresConnector(cfg *config.DatabaseConfig) Connector {
	return &PostgresConnector{cfg: cfg}
}

// Connect establishes PostgreSQL connection
func (c *PostgresConnector) Connect() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(c.cfg.GetDSN()), gormConfig(c.cfg))
}

// SQLiteConnector implements Connector for SQLite
type SQLiteConnector struct {
	cfg *config.DatabaseConfig
}

// NewSQLiteConnector creates a new SQLite connector
func NewSQLiteConnector(cfg *config.DatabaseConfig) Connector {
	return &SQLiteConnector{cfg: cfg}
}

// Connect establishes SQLite connection
func (c *SQLiteConnector) Connect() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(c.cfg.Path), gormConfig(c.cfg))
}

// unsupportedConnector handles unsupported database types
type unsupportedConnector struct {
	dbType string
}

// Connect returns error for unsupported database types
func (c *unsupportedConnector) Connect() (*gorm.DB, error) {
	logger.Error("Unsupported database type", "db_type", c.dbType)
	return nil, fmt.Errorf("unsupported database type: %s", c.dbType)
}
