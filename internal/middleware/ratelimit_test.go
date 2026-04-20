package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-notification-system/internal/ratelimiter"
)

func TestRateLimitMiddleware(t *testing.T) {

	limiter := ratelimiter.NewLimiter(1, time.Minute)

	handler := RateLimitMiddleware(limiter, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	req := httptest.NewRequest("GET", "/", nil)

	// First request → allowed
	w1 := httptest.NewRecorder()
	handler(w1, req)

	if w1.Code != 200 {
		t.Fatalf("expected 200, got %d", w1.Code)
	}

	// Second request → should be blocked
	w2 := httptest.NewRecorder()
	handler(w2, req)

	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", w2.Code)
	}
}

func TestExtractIPFromHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Forwarded-For", "192.168.1.1")

	ip := extractIP(req)

	if ip != "192.168.1.1" {
		t.Errorf("expected header IP, got %s", ip)
	}
}

func TestExtractIPFromRemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:8080"

	ip := extractIP(req)

	if ip != "127.0.0.1" {
		t.Errorf("expected remote IP, got %s", ip)
	}
}

func TestExtractIPInvalidRemoteAddr(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "invalid"

	ip := extractIP(req)

	if ip == "" {
		t.Errorf("expected fallback IP handling")
	}
}