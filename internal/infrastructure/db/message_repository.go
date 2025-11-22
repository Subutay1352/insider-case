package db

import (
	"context"
	"insider-case/internal/domain/message"

	"gorm.io/gorm"
)

// Repository implements message.Repository using GORM (works with any database)
type Repository struct {
	db *gorm.DB
}

// NewRepository creates a new Repository based on database type
func NewRepository(db *gorm.DB, dbType string) message.Repository {
	return &Repository{db: db}
}

// NewPostgresRepository creates a new Repository (kept for backward compatibility)
func NewPostgresRepository(db *gorm.DB) message.Repository {
	return &Repository{db: db}
}

// GetUnsentMessages retrieves unsent messages (status = queued OR failed with retry_count < maxRetries) up to the specified limit
func (r *Repository) GetUnsentMessages(ctx context.Context, limit int, maxRetries int) ([]*message.Message, error) {
	var messages []*message.Message
	err := r.db.WithContext(ctx).
		Where("status = ? OR (status = ? AND retry_count < ?)",
			message.MessageStatusQueued,
			message.MessageStatusFailed,
			maxRetries).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

// UpdateMessageStatus updates the status of a message
func (r *Repository) UpdateMessageStatus(ctx context.Context, id uint, status message.MessageStatus, messageID string) error {
	return r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"message_id": messageID,
		}).Error
}

// GetSentMessages retrieves all sent messages with pagination
func (r *Repository) GetSentMessages(ctx context.Context, limit, offset int) ([]*message.Message, error) {
	var messages []*message.Message
	err := r.db.WithContext(ctx).
		Where("status = ?", message.MessageStatusSent).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

// UpdateMessageRetry updates retry information for a failed message
func (r *Repository) UpdateMessageRetry(ctx context.Context, id uint, retryCount int) error {
	return r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"retry_count": retryCount,
			"status":      message.MessageStatusQueued, // Retry için tekrar queued yap
		}).Error
}

// CountSentMessages returns the total count of sent messages
func (r *Repository) CountSentMessages(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("status = ?", message.MessageStatusSent).
		Count(&count).Error
	return count, err
}

