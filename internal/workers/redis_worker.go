package workers

import (
	"encoding/json"
	"time"

	"go-notification-system/internal/config"
	"go-notification-system/internal/idempotency"
	"go-notification-system/internal/logger"
	"go-notification-system/internal/models"
	"go-notification-system/internal/queue"
	"go-notification-system/internal/redisclient"
	"go-notification-system/internal/services"

	"go.uber.org/zap"
)

var processFunc = services.ProcessEvent
var sleepFunc = time.Sleep
var pushDLQFunc = queue.PushToDLQ
var claimFunc = idempotency.TryMarkProcessing

func StartRedisWorker(workerID int) {

	logger.Info("Redis worker started",
		zap.Int("worker_id", workerID),
	)

	for {
		res, err := redisclient.Client.BRPop(
			redisclient.Ctx,
			0,
			queue.RedisQueueName,
		).Result()

		if err != nil {
			logger.Error("Redis BRPOP error",
				zap.Int("worker_id", workerID),
				zap.Error(err),
			)
			continue
		}

		var event models.Event
		if err := json.Unmarshal([]byte(res[1]), &event); err != nil {
			logger.Error("Unmarshal error",
				zap.Int("worker_id", workerID),
				zap.Error(err),
			)
			continue
		}

		processWithRetry(workerID, event)
	}
}

func processWithRetry(workerID int, event models.Event) {

	// 🔥 Idempotency check
	claimed, err := claimFunc(event.ID)
	if err != nil {
		logger.Error("Idempotency error",
			zap.String("event_id", event.ID),
			zap.String("request_id", event.RequestID),
			zap.Error(err),
		)
		return
	}

	if !claimed {
		logger.Warn("Duplicate event skipped",
			zap.Int("worker_id", workerID),
			zap.String("event_id", event.ID),
			zap.String("request_id", event.RequestID),
		)
		return
	}

	maxRetries := config.AppConfig.MaxRetries

	for attempt := 1; attempt <= maxRetries; attempt++ {

		logger.Info("Processing event",
			zap.Int("worker_id", workerID),
			zap.String("event_id", event.ID),
			zap.String("request_id", event.RequestID),
			zap.Int("attempt", attempt),
		)

		err := processFunc(event)
		if err == nil {
			logger.Info("Event processed successfully",
				zap.Int("worker_id", workerID),
				zap.String("event_id", event.ID),
				zap.String("request_id", event.RequestID),
			)
			return
		}

		logger.Error("Processing failed",
			zap.Int("worker_id", workerID),
			zap.String("event_id", event.ID),
			zap.String("request_id", event.RequestID),
			zap.Int("attempt", attempt),
			zap.Error(err),
		)

		if attempt < maxRetries {
			sleep := backoff(attempt)

			logger.Warn("Retrying event",
				zap.Int("worker_id", workerID),
				zap.String("event_id", event.ID),
				zap.String("request_id", event.RequestID),
				zap.Duration("backoff", sleep),
			)

			sleepFunc(sleep)
			continue
		}

		// 🔥 DLQ
		logger.Error("Max retries reached, sending to DLQ",
			zap.Int("worker_id", workerID),
			zap.String("event_id", event.ID),
			zap.String("request_id", event.RequestID),
		)

		pushDLQFunc(event)
	}
}

func backoff(attempt int) time.Duration {
	base := config.AppConfig.BackoffSeconds
	return time.Duration(base*(1<<uint(attempt-1))) * time.Second
}