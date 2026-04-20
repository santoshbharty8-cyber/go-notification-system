package middleware

import (
	"go-notification-system/tests/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	helpers.InitTestEnv()
}

func TestLoggingMiddleware(t *testing.T) {

	handler := LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		if _, err := w.Write([]byte("ok")); err != nil {
			t.Fatal(err)
		}
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != 200 {
		t.Errorf("expected 200")
	}
}