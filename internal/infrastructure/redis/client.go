package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"

	"insider-case/internal/config"
	"insider-case/internal/pkg/logger"
)

// InitRedis initializes the Redis connection
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Redis.ConnectTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err, "addr", addr)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
