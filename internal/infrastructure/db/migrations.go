package db

import (
	"fmt"
	"insider-case/internal/domain/message"
	"insider-case/internal/pkg/logger"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	logger.Info("Running database migrations...")

	if err := runSQLMigrations(db); err != nil {
		logger.Warn("SQL migration failed, continuing with AutoMigrate", "error", err)
	}

	if err := db.AutoMigrate(&message.Message{}); err != nil {
		return fmt.Errorf("failed to run AutoMigrate: %w", err)
	}

	logger.Info("Database migrations completed successfully")
	return nil
}

func runSQLMigrations(db *gorm.DB) error {
	migrationDirs := []string{
		"migrations",
		"./migrations",
		"/app/migrations",
		filepath.Join(".", "migrations"),
	}

	var migrationsDir string
	for _, dir := range migrationDirs {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			migrationsDir = dir
			logger.Info("Found migrations directory", "path", migrationsDir)
			break
		}
	}

	if migrationsDir == "" {
		logger.Warn("Migrations directory not found, using AutoMigrate only", "tried_paths", migrationDirs)
		return nil
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		logger.Warn("Could not read migrations directory, using AutoMigrate only", "error", err)
		return nil
	}

	var sqlFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".sql") {
			sqlFiles = append(sqlFiles, entry.Name())
		}
	}

	if len(sqlFiles) == 0 {
		logger.Info("No SQL migration files found, using AutoMigrate only")
		return nil
	}

	sort.Strings(sqlFiles)

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	for _, fileName := range sqlFiles {
		migrationPath := filepath.Join(migrationsDir, fileName)

		sqlContent, err := os.ReadFile(migrationPath)
		if err != nil {
			logger.Warn("Could not read migration file", "file", migrationPath, "error", err)
			continue
		}

		logger.Info("Executing SQL migration", "file", fileName)

		statements := splitSQLStatements(string(sqlContent))
		for i, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			if _, err := sqlDB.Exec(stmt); err != nil {
				errMsg := strings.ToLower(err.Error())
				if strings.Contains(errMsg, "already exists") ||
					strings.Contains(errMsg, "duplicate") ||
					(strings.Contains(errMsg, "relation") && strings.Contains(errMsg, "already exists")) {
					logger.Info("Statement already applied, skipping", "file", fileName, "statement", i+1)
					continue
				}
				logger.Warn("SQL migration execution warning", "file", fileName, "statement", i+1, "error", err)
				return fmt.Errorf("failed to execute migration %s (statement %d): %w", fileName, i+1, err)
			}
		}

		logger.Info("SQL migration executed successfully", "file", fileName)
	}

	return nil
}

func splitSQLStatements(sql string) []string {
	statements := []string{}
	current := strings.Builder{}

	lines := strings.Split(sql, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") {
			continue
		}

		current.WriteString(line)
		current.WriteString("\n")

		if strings.HasSuffix(trimmed, ";") {
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
		}
	}

	if current.Len() > 0 {
		stmt := strings.TrimSpace(current.String())
		if stmt != "" {
			statements = append(statements, stmt)
		}
	}

	return statements
}
