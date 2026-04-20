package validator

import (
	"testing"

	"go-notification-system/internal/models"
)

func TestValidateEvent(t *testing.T) {

	tests := []struct {
		name    string
		event   models.Event
		wantErr bool
	}{
		{
			name: "valid order created",
			event: models.Event{
				ID:   "1",
				Type: models.EventOrderCreated,
				Timestamp: 1234567890,
				Payload: map[string]interface{}{
					"order_id": "ORD1",
				},
			},
			wantErr: false,
		},
		{
			name: "missing id",
			event: models.Event{
				ID: "",
				Type: models.EventOrderCreated,
			},
			wantErr: true,
		},
		{
			name: "unknown event type",
			event: models.Event{
				ID:   "2",
				Type: "unknown",
			},
			wantErr: true,
		},
		{
			name: "missing payload",
			event: models.Event{
				ID:      "3",
				Type:    models.EventOrderCreated,
				Payload: nil,
			},
			wantErr: true,
		},
		{
			name: "missing timestamp",
			event: models.Event{
				ID:   "5",
				Type: models.EventOrderCreated,
				Payload: map[string]interface{}{
					"order_id": "ORD5",
				},
			},
			wantErr: true,
		},
		{
			name: "missing order_id in payload",
			event: models.Event{
				ID:        "6",
				Type:      models.EventOrderCreated,
				Timestamp: 123,
				Payload:   map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "missing email in user_registered",
			event: models.Event{
				ID:        "7",
				Type:      models.EventUserRegistered,
				Timestamp: 123,
				Payload:   map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "valid user_registered",
			event: models.Event{
				ID:        "8",
				Type:      models.EventUserRegistered,
				Timestamp: 123,
				Payload: map[string]interface{}{
					"email": "test@test.com",
				},
			},
			wantErr: false,
		},
		{
			name: "empty event type",
			event: models.Event{
				ID:        "9",
				Type:      "",
				Timestamp: 123,
				Payload:   map[string]interface{}{},
			},
			wantErr: true,
		},
		{
			name: "nil payload",
			event: models.Event{
				ID:        "10",
				Type:      models.EventOrderCreated,
				Timestamp: 123,
				Payload:   nil,
			},
			wantErr: true,
		},
		{
			name: "order_id wrong type",
			event: models.Event{
				ID:        "11",
				Type:      models.EventOrderCreated,
				Timestamp: 123,
				Payload: map[string]interface{}{
					"order_id": 123, // ❌ wrong type
				},
			},
			wantErr: true,
		},
		{
			name: "email wrong type",
			event: models.Event{
				ID:        "12",
				Type:      models.EventUserRegistered,
				Timestamp: 123,
				Payload: map[string]interface{}{
					"email": 123, // ❌ wrong type
				},
			},
			wantErr: true,
		},
		{
			name: "invalid event type with payload",
			event: models.Event{
				ID:        "13",
				Type:      "invalid_type",
				Timestamp: 123,
				Payload: map[string]interface{}{
					"some_key": "value",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := ValidateEvent(tt.event)

			if tt.wantErr && err == nil {
				t.Errorf("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}