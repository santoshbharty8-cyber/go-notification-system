package queue

import (
	"encoding/json"
	"errors"

	"go-notification-system/internal/models"
	"go-notification-system/internal/redisclient"
)

const RedisQueueName = "event_queue"

func PushToRedisQueue(event models.Event) error {

	// 🔥 guard against nil Redis
	if redisclient.Client == nil {
		return errors.New("redis client not initialized")
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return redisclient.Client.
		LPush(redisclient.Ctx, RedisQueueName, data).
		Err()
}