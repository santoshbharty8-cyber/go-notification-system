package queue

import (
	"testing"
	"time"

	"go-notification-system/internal/models"
)

func TestPushToDLQSuccess(t *testing.T) {

	InitDLQ(1)

	event := models.Event{ID: "1"}

	PushToDLQ(event)

	select {
	case e := <-DeadLetterQueue:
		if e.ID != "1" {
			t.Errorf("expected event ID 1, got %s", e.ID)
		}
	default:
		t.Errorf("expected event in DLQ")
	}
}

func TestPushToDLQFull(t *testing.T) {

	InitDLQ(1)

	event1 := models.Event{ID: "1"}
	event2 := models.Event{ID: "2"}

	PushToDLQ(event1)
	PushToDLQ(event2) // should NOT block, should drop

	count := 0

	select {
	case <-DeadLetterQueue:
		count++
	default:
	}

	select {
	case <-DeadLetterQueue:
		count++
	default:
	}

	if count != 1 {
		t.Errorf("expected only 1 event in DLQ, got %d", count)
	}
}

func TestPushToDLQNonBlocking(t *testing.T) {

	InitDLQ(1)

	PushToDLQ(models.Event{ID: "1"})

	done := make(chan bool)

	go func() {
		PushToDLQ(models.Event{ID: "2"})
		done <- true
	}()

	select {
	case <-done:
		// success
	case <-time.After(100 * time.Millisecond):
		t.Errorf("PushToDLQ blocked on full queue")
	}
}
