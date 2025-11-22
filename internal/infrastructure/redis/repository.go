package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"insider-case/internal/domain/message"
	"insider-case/internal/pkg/logger"
	"time"

	"github.com/go-redis/redis/v8"
)

// CacheRepository implements message.CacheRepository using Redis
type CacheRepository struct {
	client *redis.Client
	ttl    time.Duration
}

// NewCacheRepository creates a new CacheRepository
func NewCacheRepository(client *redis.Client, ttl time.Duration) message.CacheRepository {
	return &CacheRepository{
		client: client,
		ttl:    ttl,
	}
}

// SetMessageID caches the messageId with sending time
func (r *CacheRepository) SetMessageID(ctx context.Context, messageID string, sentAt time.Time) error {
	key := fmt.Sprintf("message:%s", messageID)

	data := map[string]interface{}{
		"message_id": messageID,
		"sent_at":    sentAt.Unix(),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal cache data", "error", err, "message_id", messageID)
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	return r.client.Set(ctx, key, jsonData, r.ttl).Err()
}

// GetMessageID retrieves cached messageId information
func (r *CacheRepository) GetMessageID(ctx context.Context, messageID string) (*time.Time, error) {
	key := fmt.Sprintf("message:%s", messageID)

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Not found
	}
	if err != nil {
		logger.Error("Failed to get from cache", "error", err, "message_id", messageID)
		return nil, fmt.Errorf("failed to get from cache: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		logger.Error("Failed to unmarshal cache data", "error", err, "message_id", messageID)
		return nil, fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	sentAtUnix, ok := data["sent_at"].(float64)
	if !ok {
		logger.Error("Invalid sent_at format in cache", "message_id", messageID, "data", data)
		return nil, fmt.Errorf("invalid sent_at format")
	}

	sentAt := time.Unix(int64(sentAtUnix), 0)
	return &sentAt, nil
}
