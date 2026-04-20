package ratelimiter

import (
	"sync"
	"time"
)

type clientData struct {
	timestamps []time.Time
}

type SlidingWindowLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientData
	limit    int
	window   time.Duration
}

func NewLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		clients: make(map[string]*clientData),
		limit:   limit,
		window:  window,
	}
}

func (l *SlidingWindowLimiter) Allow(clientID string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	data, exists := l.clients[clientID]
	if !exists {
		data = &clientData{}
		l.clients[clientID] = data
	}

	// Remove old timestamps
	var valid []time.Time
	for _, t := range data.timestamps {
		if now.Sub(t) <= l.window {
			valid = append(valid, t)
		}
	}

	data.timestamps = valid

	// Check limit
	if len(data.timestamps) >= l.limit {
		return false
	}

	// Add current request
	data.timestamps = append(data.timestamps, now)
	return true
}