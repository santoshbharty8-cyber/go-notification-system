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

func EventHandler(w http.ResponseWriter, r *http.Request) {

	var event models.Event

	// 🔹 Decode request body
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{
			Status:  "error",
			Message: "invalid request body",
		})
		return
	}

	// 🔥 Extract request_id from context
	requestID, _ := r.Context().Value(middleware.RequestIDKey).(string)

	// 🔥 Attach to event (VERY IMPORTANT)
	event.RequestID = requestID
	event.Timestamp = time.Now().Unix()

	// 🔹 Validate event
	err = validator.ValidateEvent(event)
	if err != nil {

		logger.Warn("Validation failed",
			zap.String("event_id", event.ID),
			zap.String("request_id", requestID),
			zap.Error(err),
		)

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	// 🔥 Push event to Redis queue
	// err = queue.PushToRedisQueue(event)
	err = enqueueFunc(event)
	if err != nil {

		logger.Error("Failed to enqueue event",
			zap.String("event_id", event.ID),
			zap.String("request_id", requestID),
			zap.Error(err),
		)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.APIResponse{
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

	// 🔥 Success response (include request_id)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(models.APIResponse{
		Status:  "success",
		Message: "event queued",
		Data: map[string]string{
			"event_id":  event.ID,
			"request_id": requestID,
		},
	})
}