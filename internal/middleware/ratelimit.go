package middleware

import (
	"encoding/json"
	"net"
	"net/http"

	"go-notification-system/internal/logger"
	"go-notification-system/internal/models"
	"go-notification-system/internal/ratelimiter"

	"go.uber.org/zap"
)

func RateLimitMiddleware(limiter *ratelimiter.SlidingWindowLimiter, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// 🔥 Extract request_id (for tracing)
		requestID, _ := r.Context().Value(RequestIDKey).(string)

		// 🔥 Safe IP extraction
		ip := extractIP(r)

		if !limiter.Allow(ip) {

			// 🔥 Add useful headers
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Retry-After", "10") // seconds (optional tuning)

			w.WriteHeader(http.StatusTooManyRequests)

			_ = json.NewEncoder(w).Encode(models.APIResponse{
				Status:  "error",
				Message: "rate limit exceeded",
			})

			// 🔥 Structured log (IMPORTANT)
			logger.Warn("Rate limit exceeded",
				zap.String("request_id", requestID),
				zap.String("ip", ip),
				zap.String("path", r.URL.Path),
			)

			return
		}

		next(w, r)
	}
}



func extractIP(r *http.Request) string {

	// If behind proxy (future-ready)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return ip
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // fallback
	}

	return ip
}