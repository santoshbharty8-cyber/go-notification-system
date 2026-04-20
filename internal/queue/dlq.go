package queue

import (
	"log"

	"go-notification-system/internal/models"
)

var DeadLetterQueue chan models.Event

func InitDLQ(bufferSize int) {
	DeadLetterQueue = make(chan models.Event, bufferSize)
}

func PushToDLQ(event models.Event) {
	select {
	case DeadLetterQueue <- event:
		log.Println("Event moved to DLQ:", event.ID)
	default:
		log.Println("DLQ FULL! Event lost:", event.ID)
	}
}