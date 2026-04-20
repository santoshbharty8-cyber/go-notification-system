package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthCheck(t *testing.T) {

	// 🔥 Set Gin to test mode (important)
	gin.SetMode(gin.TestMode)

	// Create response recorder
	w := httptest.NewRecorder()

	// Create gin context
	c, _ := gin.CreateTestContext(w)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	c.Request = req

	// Call handler
	HealthCheck(c)

	// Assertions
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if w.Body.String() == "" {
		t.Errorf("expected non-empty response body")
	}
}