package message

import (
	"context"
	"insider-case/internal/pkg/logger"
	"time"
)

// Service handles message-related business logic
type Service struct {
	repo             Repository
	cacheRepo        CacheRepository
	webhookClient    WebhookClient
	messagesPerBatch int
	maxMessageLength int
	maxRetryAttempts int
	retryBaseDelay   time.Duration
}

// NewService creates a new MessageService
func NewService(
	repo Repository,
	cacheRepo CacheRepository,
	webhookClient WebhookClient,
	messagesPerBatch int,
	maxMessageLength int,
	maxRetryAttempts int,
	retryBaseDelay time.Duration,
) *Service {
	return &Service{
		repo:             repo,
		cacheRepo:        cacheRepo,
		webhookClient:    webhookClient,
		messagesPerBatch: messagesPerBatch,
		maxMessageLength: maxMessageLength,
		maxRetryAttempts: maxRetryAttempts,
		retryBaseDelay:   retryBaseDelay,
	}
}

// SendPendingMessages processes and sends queued messages
func (s *Service) SendPendingMessages(ctx context.Context) error {
	messages, err := s.repo.GetUnsentMessages(ctx, s.messagesPerBatch, s.maxRetryAttempts)
	if err != nil {
		logger.Error("Failed to get unsent messages", "error", err)
		return &ErrRepository{Operation: "get unsent messages", Err: err}
	}

	if len(messages) == 0 {
		return nil
	}

	for _, msg := range messages {
		if err := s.processMessage(ctx, msg); err != nil {
			if err := s.handleFailedMessage(ctx, msg, err); err != nil {
				logger.Error("Failed to handle failed message",
					"message_id", msg.ID,
					"error", err,
				)
			}
			continue
		}
	}

	return nil
}

// processMessage processes a single message
func (s *Service) processMessage(ctx context.Context, msg *Message) error {
	if !msg.IsValidContent(s.maxMessageLength) {
		return &ErrContentLengthExceeded{
			Length:    len(msg.Content),
			MaxLength: s.maxMessageLength,
		}
	}

	webhookReq := &WebhookRequest{
		To:      msg.To,
		Content: msg.Content,
	}

	resp, err := s.webhookClient.SendMessage(ctx, webhookReq)
	if err != nil {
		return &ErrWebhook{Err: err}
	}

	if err := s.repo.UpdateMessageStatus(ctx, msg.ID, MessageStatusSent, resp.MessageID); err != nil {
		return &ErrRepository{Operation: "update message status", Err: err}
	}

	if s.cacheRepo != nil {
		sentAt := time.Now()
		if err := s.cacheRepo.SetMessageID(ctx, resp.MessageID, sentAt); err != nil {
			logger.Warn("Failed to cache messageId",
				"message_id", resp.MessageID,
				"error", err,
			)
		}
	}

	return nil
}

// handleFailedMessage handles retry logic for failed messages
func (s *Service) handleFailedMessage(ctx context.Context, msg *Message, err error) error {
	newRetryCount := msg.RetryCount + 1

	logger.Error("Error processing message",
		"message_id", msg.ID,
		"retry_count", newRetryCount,
		"max_retry_attempts", s.maxRetryAttempts,
		"error", err,
	)

	if newRetryCount >= s.maxRetryAttempts {
		logger.Warn("Message exceeded max retry attempts, marking as permanently failed",
			"message_id", msg.ID,
			"retry_count", newRetryCount,
			"max_retry_attempts", s.maxRetryAttempts,
		)
		return s.repo.UpdateMessageStatusAndRetry(ctx, msg.ID, MessageStatusFailed, newRetryCount)
	}

	if err := s.repo.UpdateMessageRetry(ctx, msg.ID, newRetryCount); err != nil {
		return err
	}

	logger.Info("Message scheduled for retry",
		"message_id", msg.ID,
		"retry_count", newRetryCount,
		"max_retry_attempts", s.maxRetryAttempts,
	)

	return nil
}

// GetSentMessages retrieves sent messages with pagination
func (s *Service) GetSentMessages(ctx context.Context, limit, offset int) ([]*Message, int64, error) {
	messages, err := s.repo.GetSentMessages(ctx, limit, offset)
	if err != nil {
		return nil, 0, &ErrRepository{Operation: "get sent messages", Err: err}
	}

	total, err := s.repo.CountSentMessages(ctx)
	if err != nil {
		return nil, 0, &ErrRepository{Operation: "count sent messages", Err: err}
	}

	return messages, total, nil
}
