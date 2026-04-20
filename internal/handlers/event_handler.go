package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go-notification-system/internal/logger"
	"go-notification-system/internal/middleware"
	"go-notification-system/internal/models"
	"go-notification-system/internal/queue"
	"go-notification-system/internal/validator"

	"go.uber.org/zap"
)

var enqueueFunc = queue.PushToRedisQueue

// 🔥 helper (centralized response handling)
func writeJSON(w http.ResponseWriter, status int, resp models.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error("failed to write response", zap.Error(err))
	}
}

func EventHandler(w http.ResponseWriter, r *http.Request) {

	var event models.Event

	// 🔹 Decode request body
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeJSON(w, http.StatusBadRequest, models.APIResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	// 🔥 Extract request_id
	requestID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	// 🔥 Attach metadata
	event.RequestID = requestID
	event.Timestamp = time.Now().Unix()

	// 🔹 Validate
	if err := validator.ValidateEvent(event); err != nil {

		logger.Warn("Validation failed",
			zap.String("event_id", event.ID),
			zap.String("request_id", requestID),
			zap.Error(err),
		)

		writeJSON(w, http.StatusBadRequest, models.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// 🔥 Queue push
	if err := enqueueFunc(event); err != nil {

		logger.Error("Failed to enqueue event",
			zap.String("event_id", event.ID),
			zap.String("request_id", requestID),
			zap.Error(err),
		)

		writeJSON(w, http.StatusInternalServerError, models.APIResponse{
			Status:  "error",
			Message: "failed to enqueue event",
		})
		return
	}

	// 🔥 Success log
	logger.Info("Event queued successfully",
		zap.String("event_id", event.ID),
		zap.String("request_id", requestID),
	)

	// 🔥 Success response
	writeJSON(w, http.StatusOK, models.APIResponse{
		Status:  "success",
		Message: "event queued",
		Data: map[string]string{
			"event_id":  event.ID,
			"request_id": requestID,
		},
	})
}