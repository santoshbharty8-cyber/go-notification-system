package workers

import (
	"errors"
	"testing"
	"time"

	"go-notification-system/internal/config"
	"go-notification-system/internal/idempotency"
	"go-notification-system/internal/logger"
	"go-notification-system/internal/models"
	"go-notification-system/internal/queue"
	"go-notification-system/internal/services"
)

//
// 🔥 Test Setup
//

func init() {
	config.AppConfig = &config.Config{
		AppEnv:     "test",
		LogLevel:   "error",
		MaxRetries: 3,
	}
	logger.InitLogger("test", "error")
}

//
// 🔧 Reset all mocks
//

func resetWorkerMocks() {
	processFunc = services.ProcessEvent
	sleepFunc = time.Sleep
	pushDLQFunc = queue.PushToDLQ
	claimFunc = idempotency.TryMarkProcessing
}

//
// ✅ 1. Retry → Success
//

func TestWorkerRetryThenSuccess(t *testing.T) {
	defer resetWorkerMocks()

	callCount := 0

	processFunc = func(e models.Event) error {
		callCount++
		if callCount == 1 {
			return errors.New("fail once")
		}
		return nil
	}

	sleepFunc = func(d time.Duration) {}

	claimFunc = func(id string) (bool, error) {
		return true, nil
	}

	event := models.Event{ID: "1"}

	processWithRetry(1, event)

	if callCount != 2 {
		t.Errorf("expected 2 attempts, got %d", callCount)
	}
}

//
// ✅ 2. Max Retries → DLQ
//

func TestWorkerMaxRetriesToDLQ(t *testing.T) {
	defer resetWorkerMocks()

	retryCount := 0

	processFunc = func(e models.Event) error {
		retryCount++
		return errors.New("always fail")
	}

	sleepFunc = func(d time.Duration) {}

	claimFunc = func(id string) (bool, error) {
		return true, nil
	}

	var dlqEvent models.Event
	pushDLQFunc = func(e models.Event) {
		dlqEvent = e
	}

	event := models.Event{ID: "2"}

	processWithRetry(1, event)

	if retryCount != config.AppConfig.MaxRetries {
		t.Errorf("expected %d retries, got %d",
			config.AppConfig.MaxRetries, retryCount)
	}

	if dlqEvent.ID != "2" {
		t.Errorf("expected event to be sent to DLQ")
	}
}

//
// ✅ 3. Duplicate Event → Skipped
//

func TestWorkerDuplicateEventSkipped(t *testing.T) {
	defer resetWorkerMocks()

	claimFunc = func(id string) (bool, error) {
		return false, nil // duplicate
	}

	called := false
	processFunc = func(e models.Event) error {
		called = true
		return nil
	}

	event := models.Event{ID: "3"}

	processWithRetry(1, event)

	if called {
		t.Errorf("expected duplicate event to be skipped")
	}
}

//
// ✅ 4. Idempotency Error (should not panic)
//

func TestWorkerIdempotencyError(t *testing.T) {
	defer resetWorkerMocks()

	claimFunc = func(id string) (bool, error) {
		return false, errors.New("redis error")
	}

	event := models.Event{ID: "4"}

	// should not panic
	processWithRetry(1, event)
}

//
// ✅ 5. Backoff Called Correct Number of Times
//

func TestWorkerBackoffCalled(t *testing.T) {
	defer resetWorkerMocks()

	processFunc = func(e models.Event) error {
		return errors.New("fail")
	}

	sleepCalls := 0
	sleepFunc = func(d time.Duration) {
		sleepCalls++
	}

	claimFunc = func(id string) (bool, error) {
		return true, nil
	}

	event := models.Event{ID: "5"}

	processWithRetry(1, event)

	expectedSleeps := config.AppConfig.MaxRetries - 1

	if sleepCalls != expectedSleeps {
		t.Errorf("expected %d sleep calls, got %d",
			expectedSleeps, sleepCalls)
	}
}

//
// 🔥 6. Retry Flow with Multiple Failures → Success
//

func TestWorkerClaimThenRetryFlow(t *testing.T) {
	defer resetWorkerMocks()

	claimFunc = func(id string) (bool, error) {
		return true, nil
	}

	callCount := 0

	processFunc = func(e models.Event) error {
		callCount++
		if callCount < 3 {
			return errors.New("fail")
		}
		return nil
	}

	sleepFunc = func(d time.Duration) {}

	event := models.Event{ID: "6"}

	processWithRetry(1, event)

	if callCount != 3 {
		t.Errorf("expected 3 attempts, got %d", callCount)
	}
}