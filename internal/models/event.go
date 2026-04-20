package models

type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp int64                  `json:"timestamp"`
	Payload   map[string]interface{} `json:"payload"`
	RequestID string                 `json:"request_id"`
}

const (
	EventUserRegistered = "user_registered"
	EventOrderCreated   = "order_created"
)