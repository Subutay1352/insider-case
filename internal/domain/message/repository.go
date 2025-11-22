package message

import (
	"context"
	"time"
)

// Repository defines the interface for message persistence operations
type Repository interface {
	GetUnsentMessages(ctx context.Context, limit int, maxRetryAttempts int) ([]*Message, error)
	UpdateMessageStatus(ctx context.Context, id uint, status MessageStatus, messageID string) error
	UpdateMessageStatusOnly(ctx context.Context, id uint, status MessageStatus) error
	UpdateMessageStatusAndRetry(ctx context.Context, id uint, status MessageStatus, retryCount int) error
	UpdateMessageRetry(ctx context.Context, id uint, retryCount int) error
	GetSentMessages(ctx context.Context, limit, offset int) ([]*Message, error)
	CountSentMessages(ctx context.Context) (int64, error)
}

// CacheRepository defines the interface for cache operations
type CacheRepository interface {
	SetMessageID(ctx context.Context, messageID string, sentAt time.Time) error
	GetMessageID(ctx context.Context, messageID string) (*time.Time, error)
}

// WebhookRequest represents the request payload for webhook
type WebhookRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

// WebhookResponse represents the response from webhook
type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}

// WebhookClient defines the interface for webhook operations
type WebhookClient interface {
	SendMessage(ctx context.Context, req *WebhookRequest) (*WebhookResponse, error)
}

// MessageProcessor defines the interface for processing messages (used by scheduler)
type MessageProcessor interface {
	SendPendingMessages(ctx context.Context) error
}
