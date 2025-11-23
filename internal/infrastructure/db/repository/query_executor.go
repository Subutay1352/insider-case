package repository

import (
	"context"
	"insider-case/internal/domain/message"
)

type QueryExecutor interface {
	GetUnsentMessages(ctx context.Context, db interface{}, limit int, maxRetryAttempts int) ([]*message.Message, error)
}
