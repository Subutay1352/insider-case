package repository

import (
	"context"
	"insider-case/internal/domain/message"

	"gorm.io/gorm"
)

type PostgresExecutor struct{}

func NewPostgresExecutor() QueryExecutor {
	return &PostgresExecutor{}
}

func (e *PostgresExecutor) GetUnsentMessages(ctx context.Context, db interface{}, limit int, maxRetryAttempts int) ([]*message.Message, error) {
	gormDB := db.(*gorm.DB)
	var messages []*message.Message

	query := `
		UPDATE messages
		SET status = $5
		WHERE id IN (
			SELECT id FROM messages
			WHERE (status = $1 OR (status = $2 AND retry_count < $3))
			ORDER BY created_at ASC
			LIMIT $4
			FOR UPDATE SKIP LOCKED
		)
		RETURNING id, "to", content, status, message_id, retry_count, created_at, updated_at
	`

	return messages, gormDB.WithContext(ctx).Raw(query,
		message.MessageStatusQueued,
		message.MessageStatusFailed,
		maxRetryAttempts,
		limit,
		message.MessageStatusProcessing,
	).Scan(&messages).Error
}
