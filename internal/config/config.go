package config

import (
	"fmt"
	"insider-case/internal/constants"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Env         string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Webhook     WebhookConfig
	Scheduler   SchedulerConfig
	Message     MessageConfig
	AccessToken string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port            string
	Host            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Type            string
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int           // Connection pool: max open connections
	MaxIdleConns    int           // Connection pool: max idle connections
	ConnMaxLifetime time.Duration // Connection pool: max connection lifetime
	ConnMaxIdleTime time.Duration // Connection pool: max idle time
	LogLevel        string        // GORM log level: silent, error, warn, info
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host           string
	Port           string
	Password       string
	DB             int
	TTL            time.Duration
	ConnectTimeout time.Duration
}

// WebhookConfig holds webhook configuration
type WebhookConfig struct {
	URL              string
	AuthKey          string // X-Ins-Auth-Key header value
	Timeout          time.Duration
	MaxRetryAttempts int // Maximum retry attempts for webhook calls
	RetryDelay       time.Duration
}

// MessageConfig holds message-related configuration
type MessageConfig struct {
	MaxLength     int // Maximum character limit for message content
	DefaultLimit  int // Default pagination limit
	DefaultOffset int // Default pagination offset
}

// SchedulerConfig holds scheduler configuration
type SchedulerConfig struct {
	Interval          time.Duration
	AutoStart         bool
	MessagesPerBatch  int           // Number of messages to process per batch
	ProcessingTimeout time.Duration // Timeout for processing messages in each batch
	ShutdownTimeout   time.Duration // Timeout for graceful shutdown
	RetryBaseDelay    time.Duration // Base delay for exponential backoff (e.g., 3s)
}

// LoadEnvFile loads .env file if ENV is "local"
func LoadEnvFile() error {
	if os.Getenv("ENV") != "local" {
		return nil
	}

	if err := godotenv.Load(); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

// Load loads configuration from environment variables
// Automatically loads .env file if ENV=local
func Load() *Config {
	// Load .env file if ENV is local
	if os.Getenv("ENV") == "local" {
		_ = LoadEnvFile()
	}
	return &Config{
		Env: getEnv("ENV", "production"),
		Server: ServerConfig{
			Port:            getEnv("SERVER_PORT", "8080"),
			Host:            getEnv("SERVER_HOST", "0.0.0.0"),
			ReadTimeout:     10 * time.Second,
			WriteTimeout:    10 * time.Second,
			IdleTimeout:     60 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Database: DatabaseConfig{
			Type:            getEnv("DB_TYPE", constants.DBTypePostgres),
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", constants.DefaultDBUser),
			Password:        getEnv("DB_PASSWORD", constants.DefaultDBPassword),
			Name:            getEnv("DB_NAME", constants.DefaultDBName),
			MaxOpenConns:    25,
			MaxIdleConns:    5,
			ConnMaxLifetime: 5 * time.Minute,
			ConnMaxIdleTime: 10 * time.Minute,
			LogLevel:        getEnv("DB_LOG_LEVEL", "info"), // silent, error, warn, info
		},
		Redis: RedisConfig{
			Host:           getEnv("REDIS_HOST", "localhost"),
			Port:           getEnv("REDIS_PORT", "6379"),
			Password:       getEnv("REDIS_PASSWORD", ""),
			DB:             getEnvAsInt("REDIS_DB", 0),
			TTL:            24 * time.Hour,
			ConnectTimeout: 5 * time.Second,
		},
		Webhook: WebhookConfig{
			URL:              getEnv("WEBHOOK_URL", "https://webhook.site/your-unique-id"),
			AuthKey:          getEnv("WEBHOOK_AUTH_KEY", "your-secret-key"),
			Timeout:          getEnvAsDuration("WEBHOOK_TIMEOUT", 30*time.Second),
			MaxRetryAttempts: 3,
			RetryDelay:       1 * time.Second,
		},
		Scheduler: SchedulerConfig{
			Interval:          getEnvAsDuration("SCHEDULER_INTERVAL", 2*time.Minute),
			AutoStart:         getEnvAsBool("SCHEDULER_AUTO_START", true),
			MessagesPerBatch:  getEnvAsInt("SCHEDULER_MESSAGES_PER_BATCH", 2),
			ProcessingTimeout: 30 * time.Second,
			ShutdownTimeout:   10 * time.Second,
			RetryBaseDelay:    3 * time.Second,
		},
		Message: MessageConfig{
			MaxLength:     getEnvAsInt("MESSAGE_MAX_LENGTH", 1000),
			DefaultLimit:  10,
			DefaultOffset: 0,
		},
		AccessToken: getEnv("ACCESS_TOKEN", "your-access-token"),
	}
}

// GetDSN returns database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host, c.User, c.Password, c.Name, c.Port,
	)
}
