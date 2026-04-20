package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-notification-system/internal/config"
	"go-notification-system/internal/handlers"
	"go-notification-system/internal/logger"
	"go-notification-system/internal/models"
	"go-notification-system/internal/queue"
	"go-notification-system/internal/redisclient"
	"go-notification-system/internal/workers"
)

func setup() {
	config.AppConfig = &config.Config{
		AppEnv:      "test",
		LogLevel:    "error",
		WorkerCount: 1,
		MaxRetries:  2,
		RedisURL:    "localhost:6379", 
	}

	logger.InitLogger("test", "error")

	redisclient.InitRedis(config.AppConfig.RedisURL)

	// clear queue
	redisclient.Client.Del(redisclient.Ctx, queue.RedisQueueName)
}

func TestEventFlowEndToEnd(t *testing.T) {

	setup()

	// 🔥 Start worker (real)
	go workers.StartRedisWorker(1)

	time.Sleep(500 * time.Millisecond) // allow worker to start

	// 🔥 Create request
	body := models.Event{
		ID:   "100",
		Type: models.EventOrderCreated,
		Payload: map[string]interface{}{
			"order_id": "ORD100",
		},
	}

	b, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/event", bytes.NewBuffer(b))
	w := httptest.NewRecorder()

	// Add request_id
	ctx := context.WithValue(req.Context(), "request_id", "test-req-100")
	req = req.WithContext(ctx)

	// 🔥 Call handler
	handlers.EventHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// 🔥 Wait for worker to process
	time.Sleep(2 * time.Second)

	// No assertion on logs here, but:
	// if system crashes, test will fail
}