package queue

import (
	"go-notification-system/internal/models"
)

var EventQueue chan models.Event

// Initialize queue
func InitQueue(bufferSize int) {
	EventQueue = make(chan models.Event, bufferSize)
}