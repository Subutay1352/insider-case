package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"insider-case/internal/config"
	"insider-case/internal/domain/message"
	"insider-case/internal/pkg/logger"
	"io"
	"net/http"
	"time"
)

// WebhookClient implements message.WebhookClient using HTTP
type WebhookClient struct {
	client        *http.Client
	webhookURL    string
	authKey       string // X-Ins-Auth-Key header value
	retryAttempts int
	retryDelay    time.Duration
}

// NewWebhookClient creates a new WebhookClient
func NewWebhookClient(cfg *config.Config) message.WebhookClient {
	if cfg.Webhook.URL == "" {
		panic("webhook URL cannot be empty")
	}

	return &WebhookClient{
		client: &http.Client{
			Timeout: cfg.Webhook.Timeout,
		},
	webhookURL:    cfg.Webhook.URL,
	authKey:       cfg.Webhook.AuthKey,
	retryAttempts: cfg.Webhook.MaxRetryAttempts,
	retryDelay:    cfg.Webhook.RetryDelay,
	}
}

// SendMessage sends a message to the webhook with retry logic
func (c *WebhookClient) SendMessage(ctx context.Context, req *message.WebhookRequest) (*message.WebhookResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retryAttempts; attempt++ {
		if attempt > 0 {
			logger.Warn("Retrying webhook request",
				"attempt", attempt,
				"max_attempts", c.retryAttempts+1,
				"url", c.webhookURL,
			)

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.retryDelay):
			}
		}

		resp, err := c.sendRequest(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err

		if httpErr, ok := err.(*HTTPError); ok {
			if httpErr.StatusCode >= 400 && httpErr.StatusCode < 500 {
				logger.Error("Client error, not retrying", "status_code", httpErr.StatusCode, "error", err)
				return nil, err
			}
		}
	}

	logger.Error("All retry attempts failed",
		"attempts", c.retryAttempts+1,
		"url", c.webhookURL,
		"error", lastErr,
	)
	return nil, fmt.Errorf("failed after %d attempts: %w", c.retryAttempts+1, lastErr)
}

// HTTPError represents an HTTP error
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}

// sendRequest performs a single HTTP request
func (c *WebhookClient) sendRequest(ctx context.Context, req *message.WebhookRequest) (*message.WebhookResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		logger.Error("Failed to marshal webhook request", "error", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Failed to create webhook request", "error", err, "url", c.webhookURL)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if c.authKey != "" {
		httpReq.Header.Set("X-Ins-Auth-Key", c.authKey)
	}

	resp, err := c.client.Do(httpReq)
	if err != nil {
		logger.Error("Failed to send webhook request", "error", err, "url", c.webhookURL)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read webhook response", "error", err, "status_code", resp.StatusCode)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusAccepted {
		logger.Error("Unexpected webhook status code", "status_code", resp.StatusCode, "body", string(body))
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("unexpected status code: %d, body: %s", resp.StatusCode, string(body)),
		}
	}

	if len(body) == 0 {
		logger.Error("Empty response body from webhook", "url", c.webhookURL)
		return nil, fmt.Errorf("empty response body from webhook")
	}

	var webhookResp message.WebhookResponse
	if err := json.Unmarshal(body, &webhookResp); err != nil {
		logger.Error("Failed to unmarshal webhook response", "error", err, "body", string(body))
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if webhookResp.MessageID == "" {
		logger.Error("Empty messageId in webhook response", "url", c.webhookURL, "body", string(body))
		return nil, fmt.Errorf("messageId is required in webhook response but was empty")
	}

	return &webhookResp, nil
}
