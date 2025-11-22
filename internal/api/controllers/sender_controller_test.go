package controllers

import (
	"context"
	"encoding/json"
	"insider-case/internal/domain/message"
	"insider-case/internal/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

type mockRepo struct{}

func (m *mockRepo) GetUnsentMessages(ctx context.Context, limit int, maxRetryAttempts int) ([]*message.Message, error) {
	return nil, nil
}

func (m *mockRepo) UpdateMessageStatus(ctx context.Context, id uint, status message.MessageStatus, messageID string) error {
	return nil
}

func (m *mockRepo) UpdateMessageStatusOnly(ctx context.Context, id uint, status message.MessageStatus) error {
	return nil
}

func (m *mockRepo) UpdateMessageStatusAndRetry(ctx context.Context, id uint, status message.MessageStatus, retryCount int) error {
	return nil
}

func (m *mockRepo) UpdateMessageRetry(ctx context.Context, id uint, retryCount int) error {
	return nil
}

func (m *mockRepo) GetSentMessages(ctx context.Context, limit, offset int) ([]*message.Message, error) {
	return nil, nil
}

func (m *mockRepo) CountSentMessages(ctx context.Context) (int64, error) {
	return 0, nil
}

type mockWebhook struct{}

func (m *mockWebhook) SendMessage(ctx context.Context, req *message.WebhookRequest) (*message.WebhookResponse, error) {
	return &message.WebhookResponse{Message: "Accepted", MessageID: "test-id"}, nil
}

func setupSenderController() (*SenderController, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	logger.Init("local")

	// Create real scheduler with mock service
	mockRepo := &mockRepo{}
	mockWebhook := &mockWebhook{}
	service := message.NewService(mockRepo, nil, mockWebhook, 2, 1000, 3, 3*time.Second)
	scheduler := message.NewScheduler(service, 1*time.Minute, 30*time.Second)
	controller := NewSenderController(scheduler)

	router := gin.New()
	router.POST("/sender/start", controller.Start)
	router.POST("/sender/stop", controller.Stop)
	router.GET("/sender/status", controller.Status)

	return controller, router
}

func TestSenderController_Start(t *testing.T) {
	controller, router := setupSenderController()

	req := httptest.NewRequest("POST", "/sender/start", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if !controller.scheduler.IsRunning() {
		t.Error("scheduler should be running after start")
	}
}

func TestSenderController_Start_AlreadyRunning(t *testing.T) {
	controller, router := setupSenderController()

	// Start scheduler first
	_ = controller.scheduler.Start()

	req := httptest.NewRequest("POST", "/sender/start", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSenderController_Stop(t *testing.T) {
	controller, router := setupSenderController()

	// Start scheduler first
	_ = controller.scheduler.Start()

	req := httptest.NewRequest("POST", "/sender/stop", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if controller.scheduler.IsRunning() {
		t.Error("scheduler should not be running after stop")
	}
}

func TestSenderController_Stop_NotRunning(t *testing.T) {
	_, router := setupSenderController()

	req := httptest.NewRequest("POST", "/sender/stop", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSenderController_Status(t *testing.T) {
	controller, router := setupSenderController()

	// Start scheduler first
	_ = controller.scheduler.Start()

	req := httptest.NewRequest("GET", "/sender/status", nil)
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
	if data["is_running"] != true {
		t.Errorf("expected is_running true, got %v", data["is_running"])
	}
}
