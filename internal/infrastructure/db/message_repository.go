package db

import (
	"context"
	"insider-case/internal/constants"
	"insider-case/internal/domain/message"
	"insider-case/internal/infrastructure/db/repository"
	"insider-case/internal/pkg/logger"

	"gorm.io/gorm"
)

type Repository struct {
	db            *gorm.DB
	queryExecutor repository.QueryExecutor
}

func NewRepository(db *gorm.DB, dbType string) message.Repository {
	if dbType == "" {
		dbType = constants.DBTypePostgres
	}

	return &Repository{
		db:            db,
		queryExecutor: newQueryExecutor(dbType),
	}
}

func newQueryExecutor(dbType string) repository.QueryExecutor {
	switch dbType {
	case constants.DBTypePostgres:
		return repository.NewPostgresExecutor()
	default:
		logger.Warn("Unsupported database type, falling back to PostgreSQL", "db_type", dbType)
		return repository.NewPostgresExecutor()
	}
}

func (r *Repository) GetUnsentMessages(ctx context.Context, limit int, maxRetryAttempts int) ([]*message.Message, error) {
	return r.queryExecutor.GetUnsentMessages(ctx, r.db, limit, maxRetryAttempts)
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
	result := r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("id = ? AND status = ?", id, message.MessageStatusProcessing).
		Updates(map[string]interface{}{
			"retry_count": retryCount,
			"status":      message.MessageStatusQueued,
		})

	return result.Error
}

func (r *Repository) CountSentMessages(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&message.Message{}).
		Where("status = ?", message.MessageStatusSent).
		Count(&count).Error
	return count, err
}
