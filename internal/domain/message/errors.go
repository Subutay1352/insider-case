package message

import (
	"errors"
	"fmt"
)

// Domain-specific errors
var (
	ErrMessageNotFound      = errors.New("message not found")
	ErrInvalidMessageStatus = errors.New("invalid message status")
	ErrMessageAlreadySent   = errors.New("message already sent")
	ErrInvalidContent       = errors.New("message content is invalid")
	ErrSchedulerRunning     = errors.New("scheduler is already running")
	ErrSchedulerNotRunning  = errors.New("scheduler is not running")
	ErrSchedulerTimeout     = errors.New("scheduler shutdown timeout")

	// Validation errors
	ErrToFieldRequired      = errors.New("to field is required")
	ErrContentFieldRequired = errors.New("content field is required")
)

// ErrContentLengthExceeded represents content length validation error
type ErrContentLengthExceeded struct {
	Length    int
	MaxLength int
}

func (e *ErrContentLengthExceeded) Error() string {
	return fmt.Sprintf("content length (%d) exceeds maximum allowed length (%d)", e.Length, e.MaxLength)
}

// ErrRepository wraps repository errors
type ErrRepository struct {
	Operation string
	Err       error
}

func (e *ErrRepository) Error() string {
	return fmt.Sprintf("failed to %s: %v", e.Operation, e.Err)
}

func (e *ErrRepository) Unwrap() error {
	return e.Err
}

// ErrWebhook wraps webhook errors
type ErrWebhook struct {
	Err error
}

func (e *ErrWebhook) Error() string {
	return fmt.Sprintf("failed to send message via webhook: %v", e.Err)
}

func (e *ErrWebhook) Unwrap() error {
	return e.Err
}
