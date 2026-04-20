package validator

import (
	"errors"

	"go-notification-system/internal/models"
)

func ValidateEvent(e models.Event) error {

	// 🔹 Basic validation
	if e.ID == "" {
		return errors.New("event id is required")
	}

	if e.Type == "" {
		return errors.New("event type is required")
	}

	if e.Timestamp == 0 {
		return errors.New("timestamp is required")
	}

	if e.Payload == nil {
		return errors.New("payload cannot be empty")
	}

	// 🔥 Type validation
	switch e.Type {

	case models.EventOrderCreated:
		return validateOrderCreated(e)

	case models.EventUserRegistered:
		return validateUserRegistered(e)

	default:
		return errors.New("invalid event type")
	}
}


func validateOrderCreated(e models.Event) error {

	orderID, ok := e.Payload["order_id"]
	if !ok {
		return errors.New("order_id is required")
	}

	str, ok := orderID.(string)
	if !ok || str == "" {
		return errors.New("order_id must be a non-empty string")
	}

	return nil
}

func validateUserRegistered(e models.Event) error {

	email, ok := e.Payload["email"]
	if !ok {
		return errors.New("email is required")
	}

	str, ok := email.(string)
	if !ok || str == "" {
		return errors.New("email must be a non-empty string")
	}

	return nil
}