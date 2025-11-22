package controllers

import (
	"context"
	"encoding/json"
	"insider-case/internal/config"
	"insider-case/internal/domain/message"
	"insider-case/internal/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// MockRepository for testing
type MockRepository struct {
	GetSentMessagesFunc   func(ctx context.Context, limit, offset int) ([]*message.Message, error)
	CountSentMessagesFunc func(ctx context.Context) (int64, error)
}

func (m *MockRepository) GetUnsentMessages(ctx context.Context, limit int, maxRetryAttempts int) ([]*message.Message, error) {
	return nil, nil
}

func (m *MockRepository) UpdateMessageStatus(ctx context.Context, id uint, status message.MessageStatus, messageID string) error {
	return nil
}

func (m *MockRepository) UpdateMessageStatusOnly(ctx context.Context, id uint, status message.MessageStatus) error {
	return nil
}

func (m *MockRepository) UpdateMessageStatusAndRetry(ctx context.Context, id uint, status message.MessageStatus, retryCount int) error {
	return nil
}

func (m *MockRepository) UpdateMessageRetry(ctx context.Context, id uint, retryCount int) error {
	return nil
}

func (m *MockRepository) GetSentMessages(ctx context.Context, limit, offset int) ([]*message.Message, error) {
	if m.GetSentMessagesFunc != nil {
		return m.GetSentMessagesFunc(ctx, limit, offset)
	}
	return []*message.Message{}, nil
}

func (m *MockRepository) CountSentMessages(ctx context.Context) (int64, error) {
	if m.CountSentMessagesFunc != nil {
		return m.CountSentMessagesFunc(ctx)
	}
	return 0, nil
}

// MockWebhookClient for testing
type MockWebhookClient struct{}

func (m *MockWebhookClient) SendMessage(ctx context.Context, req *message.WebhookRequest) (*message.WebhookResponse, error) {
	return &message.WebhookResponse{Message: "Accepted", MessageID: "test-id"}, nil
}

func setupMessageController() (*MessageController, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	logger.Init("local")

	mockRepo := &MockRepository{}
	mockWebhook := &MockWebhookClient{}
	service := message.NewService(mockRepo, nil, mockWebhook, 2, 1000, 3, 3*time.Second)
	msgConfig := &config.MessageConfig{
		MaxLength:     1000,
		DefaultLimit:  10,
		DefaultOffset: 0,
	}
	controller := NewMessageController(service, msgConfig)

	router := gin.New()
	router.GET("/messages/sent", controller.GetSentMessages)

	return controller, router
}

func TestMessageController_GetSentMessages(t *testing.T) {
	controller, router := setupMessageController()

	mockMessages := []*message.Message{
		{
			ID:        1,
			To:        "+905551111111",
			Content:   "Test message",
			Status:    message.MessageStatusSent,
			MessageID: "msg-123",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Setup mock repository
	mockRepo := &MockRepository{
		GetSentMessagesFunc: func(ctx context.Context, limit, offset int) ([]*message.Message, error) {
			return mockMessages, nil
		},
		CountSentMessagesFunc: func(ctx context.Context) (int64, error) {
			return 1, nil
		},
	}
	mockWebhook := &MockWebhookClient{}
	service := message.NewService(mockRepo, nil, mockWebhook, 2, 1000, 3, 3*time.Second)
	msgConfig := &config.MessageConfig{
		MaxLength:     1000,
		DefaultLimit:  10,
		DefaultOffset: 0,
	}
	controller = NewMessageController(service, msgConfig)
	router = gin.New()
	router.GET("/messages/sent", controller.GetSentMessages)

	req := httptest.NewRequest("GET", "/messages/sent?limit=10&offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	if data["total"] != float64(1) {
		t.Errorf("expected total 1, got %v", data["total"])
	}
}

func TestMessageController_GetSentMessages_WithPagination(t *testing.T) {
	mockRepo := &MockRepository{
		GetSentMessagesFunc: func(ctx context.Context, limit, offset int) ([]*message.Message, error) {
			if limit != 5 || offset != 10 {
				t.Errorf("expected limit=5, offset=10, got limit=%d, offset=%d", limit, offset)
			}
			return []*message.Message{}, nil
		},
		CountSentMessagesFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
	}
	service := message.NewService(mockRepo, nil, &MockWebhookClient{}, 2, 1000, 3, 3*time.Second)
	msgConfig := &config.MessageConfig{
		MaxLength:     1000,
		DefaultLimit:  10,
		DefaultOffset: 0,
	}
	controller := NewMessageController(service, msgConfig)

	router := gin.New()
	router.GET("/messages/sent", controller.GetSentMessages)

	req := httptest.NewRequest("GET", "/messages/sent?limit=5&offset=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestMessageController_GetSentMessages_DefaultPagination(t *testing.T) {
	mockRepo := &MockRepository{
		GetSentMessagesFunc: func(ctx context.Context, limit, offset int) ([]*message.Message, error) {
			if limit != 10 || offset != 0 {
				t.Errorf("expected limit=10, offset=0, got limit=%d, offset=%d", limit, offset)
			}
			return []*message.Message{}, nil
		},
		CountSentMessagesFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
	}
	service := message.NewService(mockRepo, nil, &MockWebhookClient{}, 2, 1000, 3, 3*time.Second)
	msgConfig := &config.MessageConfig{
		MaxLength:     1000,
		DefaultLimit:  10,
		DefaultOffset: 0,
	}
	controller := NewMessageController(service, msgConfig)

	router := gin.New()
	router.GET("/messages/sent", controller.GetSentMessages)

	req := httptest.NewRequest("GET", "/messages/sent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestMessageController_GetSentMessages_InvalidLimit(t *testing.T) {
	mockRepo := &MockRepository{
		GetSentMessagesFunc: func(ctx context.Context, limit, offset int) ([]*message.Message, error) {
			// Invalid limit should default to 10
			if limit != 10 {
				t.Errorf("expected limit=10 (default), got %d", limit)
			}
			return []*message.Message{}, nil
		},
		CountSentMessagesFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
	}
	service := message.NewService(mockRepo, nil, &MockWebhookClient{}, 2, 1000, 3, 3*time.Second)
	msgConfig := &config.MessageConfig{
		MaxLength:     1000,
		DefaultLimit:  10,
		DefaultOffset: 0,
	}
	controller := NewMessageController(service, msgConfig)

	router := gin.New()
	router.GET("/messages/sent", controller.GetSentMessages)

	req := httptest.NewRequest("GET", "/messages/sent?limit=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestMessageController_GetSentMessages_ServiceError(t *testing.T) {
	mockRepo := &MockRepository{
		GetSentMessagesFunc: func(ctx context.Context, limit, offset int) ([]*message.Message, error) {
			return nil, context.DeadlineExceeded
		},
		CountSentMessagesFunc: func(ctx context.Context) (int64, error) {
			return 0, nil
		},
	}
	service := message.NewService(mockRepo, nil, &MockWebhookClient{}, 2, 1000, 3, 3*time.Second)
	msgConfig := &config.MessageConfig{
		MaxLength:     1000,
		DefaultLimit:  10,
		DefaultOffset: 0,
	}
	controller := NewMessageController(service, msgConfig)

	router := gin.New()
	router.GET("/messages/sent", controller.GetSentMessages)

	req := httptest.NewRequest("GET", "/messages/sent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}
