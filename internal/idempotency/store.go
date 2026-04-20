package idempotency

import (
	"errors"
	"fmt"
	"time"

	"go-notification-system/internal/redisclient"

	"github.com/redis/go-redis/v9"
)

const KeyPrefix = "processed_event:"
const DefaultTTL = 24 * time.Hour

var ErrRedisNotInitialized = errors.New("redis client not initialized")

// TryMarkProcessing uses SET NX + TTL
func TryMarkProcessing(eventID string) (bool, error) {

	if redisclient.Client == nil {
		return false, ErrRedisNotInitialized
	}

	key := buildKey(eventID)

	cmd := redisclient.Client.SetArgs(
		redisclient.Ctx,
		key,
		"1",
		redis.SetArgs{
			Mode: "NX",
			TTL:  DefaultTTL,
		},
	)

	res, err := cmd.Result()

	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return res == "OK", nil
}

func IsProcessed(eventID string) (bool, error) {

	if redisclient.Client == nil {
		return false, ErrRedisNotInitialized
	}

	key := buildKey(eventID)

	exists, err := redisclient.Client.Exists(
		redisclient.Ctx,
		key,
	).Result()

	if err != nil {
		return false, err
	}

	return exists == 1, nil
}

func DeleteMark(eventID string) error {

	if redisclient.Client == nil {
		return ErrRedisNotInitialized
	}

	key := buildKey(eventID)

	return redisclient.Client.Del(
		redisclient.Ctx,
		key,
	).Err()
}

func ExtendTTL(eventID string, ttl time.Duration) error {

	if redisclient.Client == nil {
		return ErrRedisNotInitialized
	}

	key := buildKey(eventID)

	return redisclient.Client.Expire(
		redisclient.Ctx,
		key,
		ttl,
	).Err()
}

func buildKey(eventID string) string {
	return fmt.Sprintf("%s%s", KeyPrefix, eventID)
}