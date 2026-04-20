package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-notification-system/internal/config"
	"go-notification-system/internal/logger"
	"go-notification-system/internal/middleware"
	"go-notification-system/internal/models"
	"go-notification-system/internal/queue"
)


func init() {
	config.AppConfig = &config.Config{
		AppEnv:   "test",
		LogLevel: "error", // minimal logs
	}
	logger.InitLogger("test", "error")
}


func resetMock() {
	enqueueFunc = queue.PushToRedisQueue
}


func TestEventHandlerSuccess(t *testing.T) {

	// 🔥 mock Redis enqueue
	enqueueFunc = func(e models.Event) error {
		return nil
	}
	defer resetMock()

	body := `{
		"id": "1",
		"type": "order_created",
		"payload": {"order_id": "ORD1"}
	}`

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString(body))

	// inject request_id
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-req-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	EventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp models.APIResponse
	_ = json.Unmarshal(w.Body.Bytes(), &resp)

	if resp.Status != "success" {
		t.Errorf("expected success, got %s", resp.Status)
	}

	data := resp.Data.(map[string]interface{})

	if data["request_id"] != "test-req-123" {
		t.Errorf("request_id not propagated")
	}
}


func TestEventHandlerInvalidJSON(t *testing.T) {

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString("{invalid"))
	w := httptest.NewRecorder()

	EventHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}


func TestEventHandlerValidationError(t *testing.T) {

	body := `{
		"id": "",
		"type": "order_created"
	}`

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	EventHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}



func TestEventHandlerQueueFailure(t *testing.T) {

	enqueueFunc = func(e models.Event) error {
		return assertError() // simulate failure
	}
	defer resetMock()

	body := `{
		"id": "2",
		"type": "order_created",
		"payload": {"order_id": "ORD2"}
	}`

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	EventHandler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func assertError() error {
	return &customError{}
}

type customError struct{}

func (e *customError) Error() string {
	return "mock error"
}



func TestEventHandlerNoRequestID(t *testing.T) {

	enqueueFunc = func(e models.Event) error {
		return nil
	}
	defer resetMock()

	body := `{
		"id": "3",
		"type": "order_created",
		"payload": {"order_id": "ORD3"}
	}`

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	EventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}