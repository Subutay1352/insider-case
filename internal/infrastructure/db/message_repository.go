package db

import (
	"context"
	"insider-case/internal/domain/message"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB, dbType string) message.Repository {
	return &Repository{db: db}
}

func NewPostgresRepository(db *gorm.DB) message.Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUnsentMessages(ctx context.Context, limit int, maxRetryAttempts int) ([]*message.Message, error) {
	var messages []*message.Message
	err := r.db.WithContext(ctx).
		Where("status = ? OR (status = ? AND retry_count < ?)",
			message.MessageStatusQueued,
			message.MessageStatusFailed,
			maxRetryAttempts).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error
	return messages, err
}

func (r *Repository) UpdateMessageStatus(ctx context.Context, id uint, status message.MessageStatus, messageID string) error {
	return r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"message_id": messageID,
		}).Error
}

func (r *Repository) UpdateMessageStatusOnly(ctx context.Context, id uint, status message.MessageStatus) error {
	return r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *Repository) UpdateMessageStatusAndRetry(ctx context.Context, id uint, status message.MessageStatus, retryCount int) error {
	return r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"retry_count": retryCount,
		}).Error
}

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

func (r *Repository) UpdateMessageRetry(ctx context.Context, id uint, retryCount int) error {
	return r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"retry_count": retryCount,
			"status":      message.MessageStatusQueued,
		}).Error
}

func (r *Repository) CountSentMessages(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("status = ?", message.MessageStatusSent).
		Count(&count).Error
	return count, err
}
