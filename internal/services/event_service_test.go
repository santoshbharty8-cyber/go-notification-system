package services

import (
	"testing"

	"go-notification-system/internal/config"
	"go-notification-system/internal/logger"
	"go-notification-system/internal/models"
)

func resetRandom() {
	randomFail = func() bool {
		return false
	}
}

func init() {
	config.AppConfig = &config.Config{
		AppEnv:   "test",
		LogLevel: "error", // keep logs silent
	}
	logger.InitLogger("test", "error")
}

func TestProcessEventOrderCreatedSuccess(t *testing.T) {

	randomFail = func() bool { return false }
	defer resetRandom()

	event := models.Event{
		ID:   "1",
		Type: models.EventOrderCreated,
		Payload: map[string]interface{}{
			"order_id": "ORD1",
		},
	}

	err := ProcessEvent(event)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestProcessEventUserRegisteredSuccess(t *testing.T) {

	randomFail = func() bool { return false }
	defer resetRandom()

	event := models.Event{
		ID:   "2",
		Type: models.EventUserRegistered,
		Payload: map[string]interface{}{
			"email": "test@test.com",
		},
	}

	err := ProcessEvent(event)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestProcessEventFailure(t *testing.T) {

	randomFail = func() bool { return true }
	defer resetRandom()

	event := models.Event{
		ID:   "3",
		Type: models.EventOrderCreated,
		Payload: map[string]interface{}{
			"order_id": "ORD3",
		},
	}

	err := ProcessEvent(event)

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestProcessEventUnknownType(t *testing.T) {

	randomFail = func() bool { return false }
	defer resetRandom()

	event := models.Event{
		ID:   "4",
		Type: "unknown",
	}

	err := ProcessEvent(event)

	// Should NOT fail (based on your implementation)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestProcessEventTableDriven(t *testing.T) {

	randomFail = func() bool { return false }
	defer resetRandom()

	tests := []struct {
		name  string
		event models.Event
	}{
		{
			name: "order created",
			event: models.Event{
				ID:   "5",
				Type: models.EventOrderCreated,
				Payload: map[string]interface{}{
					"order_id": "ORD5",
				},
			},
		},
		{
			name: "user registered",
			event: models.Event{
				ID:   "6",
				Type: models.EventUserRegistered,
				Payload: map[string]interface{}{
					"email": "a@test.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := ProcessEvent(tt.event)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestProcessEventUnknownTypeCoverage(t *testing.T) {
	randomFail = func() bool { return false }

	event := models.Event{
		ID:   "x",
		Type: "invalid_type",
	}

	_ = ProcessEvent(event)
}