package db

import (
	"context"
	"insider-case/internal/domain/message"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db     *gorm.DB
	dbType string
}

func NewRepository(db *gorm.DB, dbType string) message.Repository {
	return &Repository{db: db, dbType: dbType}
}

func NewPostgresRepository(db *gorm.DB) message.Repository {
	return &Repository{db: db, dbType: DBTypePostgres}
}

func (r *Repository) GetUnsentMessages(ctx context.Context, limit int, maxRetryAttempts int) ([]*message.Message, error) {
	if r.dbType == DBTypePostgres {
		return r.getUnsentMessagesPostgres(ctx, limit, maxRetryAttempts)
	}
	return r.getUnsentMessagesSQLite(ctx, limit, maxRetryAttempts)
}

func (r *Repository) getUnsentMessagesPostgres(ctx context.Context, limit int, maxRetryAttempts int) ([]*message.Message, error) {
	var messages []*message.Message

	query := `
		UPDATE messages
		SET status = 'processing'
		WHERE id IN (
			SELECT id FROM messages
			WHERE (status = $1 OR (status = $2 AND retry_count < $3))
			ORDER BY created_at ASC
			LIMIT $4
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, "to", content, status, message_id, retry_count, created_at, updated_at
	`

	return messages, r.db.WithContext(ctx).Raw(query,
		message.MessageStatusQueued,
		message.MessageStatusFailed,
		maxRetryAttempts,
		limit,
	).Scan(&messages).Error
}

func (r *Repository) getUnsentMessagesSQLite(ctx context.Context, limit int, maxRetryAttempts int) ([]*message.Message, error) {
	var messages []*message.Message

	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("status = ? OR (status = ? AND retry_count < ?)",
		message.MessageStatusQueued,
		message.MessageStatusFailed,
		maxRetryAttempts).
		Order("created_at ASC").
		Limit(limit).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Find(&messages).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if len(messages) == 0 {
		tx.Rollback()
		return nil, nil
	}

	ids := make([]uint, len(messages))
	for i, msg := range messages {
		ids[i] = msg.ID
		msg.Status = message.MessageStatusProcessing
	}

	if err := tx.Model(&message.Message{}).
		Where("id IN ?", ids).
		Update("status", message.MessageStatusProcessing).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return messages, tx.Commit().Error
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
