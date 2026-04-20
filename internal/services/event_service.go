package services

import (
	"errors"
	"math/rand"

	"go-notification-system/internal/logger"
	"go-notification-system/internal/models"

	"go.uber.org/zap"
)

var randomFail = func() bool {
	return rand.Intn(3) == 0
}

func ProcessEvent(event models.Event) error {

	// 🔥 Simulate random failure
	if randomFail() {
		err := errors.New("random processing failure")

		logger.Error("Processing failed",
			zap.String("event_id", event.ID),
			zap.String("request_id", event.RequestID),
			zap.Error(err),
		)

		return err
	}

	switch event.Type {

	case models.EventUserRegistered:
		handleUserRegistered(event)

	case models.EventOrderCreated:
		handleOrderCreated(event)

	default:
		logger.Warn("Unknown event type",
			zap.String("event_id", event.ID),
			zap.String("request_id", event.RequestID),
			zap.String("event_type", string(event.Type)),
		)
	}

	return nil
}

func handleUserRegistered(event models.Event) {

	email, _ := event.Payload["email"].(string)

	logger.Info("Processing user registered",
		zap.String("event_id", event.ID),
		zap.String("request_id", event.RequestID),
		zap.String("email", email),
	)

	logger.Info("Sending welcome email",
		zap.String("event_id", event.ID),
		zap.String("request_id", event.RequestID),
		zap.String("email", email),
	)
}

func handleOrderCreated(event models.Event) {

	orderID, _ := event.Payload["order_id"].(string)

	logger.Info("Processing order created",
		zap.String("event_id", event.ID),
		zap.String("request_id", event.RequestID),
		zap.String("order_id", orderID),
	)

	logger.Info("Triggering order workflow",
		zap.String("event_id", event.ID),
		zap.String("request_id", event.RequestID),
		zap.String("order_id", orderID),
	)
}