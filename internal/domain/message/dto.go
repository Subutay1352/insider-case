package message

import (
	"unicode/utf8"
)

// DTOs (Data Transfer Objects) for API requests and responses

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	To      string `json:"to" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// Validate validates the SendMessageRequest
func (r *SendMessageRequest) Validate(maxLength int) error {
	if r.To == "" {
		return ErrToFieldRequired
	}

	if r.Content == "" {
		return ErrContentFieldRequired
	}

	contentLength := utf8.RuneCountInString(r.Content)
	if contentLength > maxLength {
		return &ErrContentLengthExceeded{
			Length:    contentLength,
			MaxLength: maxLength,
		}
	}

	return nil
}

// MessageResponse represents a message in API responses
type MessageResponse struct {
	ID        uint   `json:"id"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Status    string `json:"status"`
	MessageID string `json:"message_id,omitempty"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// SentMessagesResponse represents paginated sent messages response
type SentMessagesResponse struct {
	Messages []*MessageResponse `json:"messages"`
	Total    int64              `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}
