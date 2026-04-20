package queue

import (
	"testing"

	"go-notification-system/internal/config"
	"go-notification-system/internal/logger"
	"go-notification-system/internal/models"
	"go-notification-system/internal/redisclient"
)

func setup() {
	config.AppConfig = &config.Config{
		AppEnv:   "test",
		LogLevel: "error",
		RedisURL: "localhost:6379",
	}

	logger.InitLogger("test", "error") // 🔥 IMPORTANT

	redisclient.InitRedis(config.AppConfig.RedisURL)
}

func TestInitQueue(t *testing.T) {
	InitQueue(10)
	// If it sets globals, optionally assert non-nil.
}

func TestPushToRedisQueue(t *testing.T) {
	setup()

	event := models.Event{
		ID: "q1",
	}

	err := PushToRedisQueue(event)

	if err != nil {
		t.Errorf("failed to push event")
	}
}

func TestPushToRedisQueueSuccess(t *testing.T) {
	setup()

	event := models.Event{ID: "q1"}

	err := PushToRedisQueue(event)
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestPushToRedisQueueRedisNil(t *testing.T) {
	setup()

	redisclient.Client = nil
	defer redisclient.InitRedis(config.AppConfig.RedisURL)

	err := PushToRedisQueue(models.Event{ID: "q2"})
	if err == nil {
		t.Errorf("expected error when redis is nil")
	}
}

func TestPushToRedisQueueRedisClosed(t *testing.T) {
	setup()

	if err := redisclient.Client.Close(); err != nil {
		t.Log(err)
	}

	err := PushToRedisQueue(models.Event{ID: "q3"})
	if err == nil {
		t.Errorf("expected error when redis is closed")
	}

	// restore
	redisclient.InitRedis(config.AppConfig.RedisURL)
}

func TestPushToRedisQueueMarshalError(t *testing.T) {
	setup()

	// JSON cannot marshal channels → guaranteed error
	event := models.Event{
		ID: "bad",
		Payload: map[string]interface{}{
			"invalid": make(chan int),
		},
	}

	err := PushToRedisQueue(event)

	if err == nil {
		t.Errorf("expected marshal error")
	}
}