package idempotency

import (
	"testing"
	"time"

	"go-notification-system/internal/config"
	"go-notification-system/internal/redisclient"
	"go-notification-system/tests/helpers"
)

func init() {
	helpers.InitTestEnv()
}

func setup() {
	config.AppConfig = &config.Config{
		RedisURL: "localhost:6379",
	}
	redisclient.InitRedis(config.AppConfig.RedisURL)

	// clean before tests
	redisclient.Client.FlushDB(redisclient.Ctx)
}

func TestIsProcessed(t *testing.T) {
	setup()

	id := "id-1"

	// initially false
	ok, err := IsProcessed(id)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Errorf("expected false initially")
	}

	// mark it
	_, _ = TryMarkProcessing(id)

	ok, err = IsProcessed(id)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Errorf("expected true after marking")
	}
}

func TestDeleteMark(t *testing.T) {
	setup()

	id := "id-2"

	_, _ = TryMarkProcessing(id)

	DeleteMark(id)

	ok, _ := IsProcessed(id)
	if ok {
		t.Errorf("expected false after delete")
	}
}

func TestExtendTTL(t *testing.T) {
	setup()

	id := "id-3"

	_, _ = TryMarkProcessing(id)

	// extend TTL
	err := ExtendTTL(id, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	ttl, err := redisclient.Client.TTL(redisclient.Ctx, buildKey(id)).Result()
	if err != nil {
		t.Fatal(err)
	}

	if ttl <= 0 {
		t.Errorf("expected TTL to be extended")
	}
}

func TestTryMarkProcessingDuplicate(t *testing.T) {
	setup()

	id := "id-4"

	ok, _ := TryMarkProcessing(id)
	if !ok {
		t.Errorf("expected first call true")
	}

	ok, _ = TryMarkProcessing(id)
	if ok {
		t.Errorf("expected duplicate false")
	}
}

func TestTryMarkProcessingRedisError(t *testing.T) {
	setup()

	// simulate Redis failure
	redisclient.Client = nil

	defer redisclient.InitRedis(config.AppConfig.RedisURL)

	_, err := TryMarkProcessing("err-id")

	if err == nil {
		t.Errorf("expected error when Redis is nil")
	}
}

func TestIsProcessedRedisError(t *testing.T) {
	setup()

	redisclient.Client = nil

	defer redisclient.InitRedis(config.AppConfig.RedisURL)

	_, err := IsProcessed("err-id")

	if err == nil {
		t.Errorf("expected error when Redis is nil")
	}
}

func TestDeleteMarkRedisNil(t *testing.T) {
	setup()

	redisclient.Client = nil
	defer redisclient.InitRedis(config.AppConfig.RedisURL)

	err := DeleteMark("id")

	if err == nil {
		t.Errorf("expected error when redis is nil")
	}
}

func TestExtendTTLRedisNil(t *testing.T) {
	setup()

	redisclient.Client = nil
	defer redisclient.InitRedis(config.AppConfig.RedisURL)

	err := ExtendTTL("id", 10*time.Second)

	if err == nil {
		t.Errorf("expected error when redis is nil")
	}
}

func TestIsProcessedRedisClosed(t *testing.T) {
	setup()

	redisclient.Client.Close()

	_, err := IsProcessed("id")

	if err == nil {
		t.Errorf("expected error")
	}

	// restore
	redisclient.InitRedis(config.AppConfig.RedisURL)
}

func TestTryMarkProcessingRedisClosed(t *testing.T) {
	setup()

	redisclient.Client.Close()

	_, err := TryMarkProcessing("id")

	if err == nil {
		t.Errorf("expected error")
	}

	redisclient.InitRedis(config.AppConfig.RedisURL)
}