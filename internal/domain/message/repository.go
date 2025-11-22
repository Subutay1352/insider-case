package message

import (
	"context"
	"time"
)

// Repository defines the interface for message persistence operations
type Repository interface {
	// GetUnsentMessages retrieves unsent messages (status = queued OR failed with retry_count < maxRetries) up to the specified limit
	GetUnsentMessages(ctx context.Context, limit int, maxRetries int) ([]*Message, error)

	// UpdateMessageStatus updates the status of a message
	UpdateMessageStatus(ctx context.Context, id uint, status MessageStatus, messageID string) error

	// UpdateMessageRetry updates retry information for a failed message
	UpdateMessageRetry(ctx context.Context, id uint, retryCount int) error

	// GetSentMessages retrieves all sent messages
	GetSentMessages(ctx context.Context, limit, offset int) ([]*Message, error)

	// CountSentMessages returns the total count of sent messages
	CountSentMessages(ctx context.Context) (int64, error)
}

// CacheRepository defines the interface for cache operations
type CacheRepository interface {
	// SetMessageID caches the messageId with sending time
	SetMessageID(ctx context.Context, messageID string, sentAt time.Time) error

	// GetMessageID retrieves cached messageId information
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
