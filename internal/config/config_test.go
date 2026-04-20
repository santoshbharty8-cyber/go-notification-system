package config

import (
	"os"
	"testing"
)

// 🔥 helper for safe env setup
func mustSetEnv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatal(err)
	}
}

func TestLoadConfigDefaults(t *testing.T) {

	// clear env (important)
	os.Clearenv()

	LoadConfig()

	if AppConfig == nil {
		t.Fatalf("config not loaded")
	}

	if AppConfig.RedisURL == "" {
		t.Errorf("expected default RedisURL")
	}

	if AppConfig.WorkerCount <= 0 {
		t.Errorf("invalid worker count")
	}

	if AppConfig.MaxRetries <= 0 {
		t.Errorf("invalid retries")
	}
}

func TestLoadConfigWithEnv(t *testing.T) {

	mustSetEnv(t, "APP_ENV", "test")
	mustSetEnv(t, "REDIS_URL", "localhost:6379")
	mustSetEnv(t, "WORKER_COUNT", "5")
	mustSetEnv(t, "MAX_RETRIES", "4")
	mustSetEnv(t, "DLQ_SIZE", "50")

	LoadConfig()

	if AppConfig.WorkerCount != 5 {
		t.Errorf("expected worker count 5")
	}

	if AppConfig.MaxRetries != 4 {
		t.Errorf("expected retries 4")
	}

	if AppConfig.DLQSize != 50 {
		t.Errorf("expected DLQ size 50")
	}
}

func TestGetEnv(t *testing.T) {

	mustSetEnv(t, "TEST_KEY", "value")

	val := getEnv("TEST_KEY", "default")

	if val != "value" {
		t.Errorf("expected value, got %s", val)
	}
}

func TestGetEnvDefault(t *testing.T) {

	val := getEnv("UNKNOWN_KEY", "default")

	if val != "default" {
		t.Errorf("expected default")
	}
}

func TestGetEnvAsInt(t *testing.T) {

	mustSetEnv(t, "INT_KEY", "5")

	val := getEnvAsInt("INT_KEY", 1)

	if val != 5 {
		t.Errorf("expected 5, got %d", val)
	}
}

func TestGetEnvAsIntInvalid(t *testing.T) {

	mustSetEnv(t, "INT_KEY", "abc")

	val := getEnvAsInt("INT_KEY", 2)

	if val != 2 {
		t.Errorf("expected fallback value")
	}
}

func TestLoadConfig_Defaults(t *testing.T) {

	os.Clearenv()

	LoadConfig()

	if AppConfig == nil {
		t.Fatalf("config not loaded")
	}

	if AppConfig.RedisURL == "" {
		t.Errorf("expected default RedisURL")
	}

	if AppConfig.WorkerCount <= 0 {
		t.Errorf("invalid worker count")
	}

	if AppConfig.MaxRetries <= 0 {
		t.Errorf("invalid retries")
	}
}

func TestLoadConfigWithEnvFile(t *testing.T) {

	// create temp env file
	file := ".env.test"
	content := "WORKER_COUNT=7\n"

	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// cleanup (safe ignore)
	defer func() {
		_ = os.Remove(file)
	}()

	mustSetEnv(t, "APP_ENV", "test")

	LoadConfig()

	if AppConfig.WorkerCount != 7 {
		t.Errorf("expected worker count from env file")
	}
}

func TestValidateConfigAllFailures(t *testing.T) {

	tests := []Config{
		{
			RedisURL:          "x",
			WorkerCount:       1,
			MaxRetries:        1,
			BackoffSeconds:    1,
			RateLimit:         0, // ❌
			RateWindowSeconds: 10,
			DLQSize:           100,
		},
		{
			RedisURL:          "x",
			WorkerCount:       1,
			MaxRetries:        1,
			BackoffSeconds:    1,
			RateLimit:         5,
			RateWindowSeconds: 0, // ❌
			DLQSize:           100,
		},
		{
			RedisURL:          "x",
			WorkerCount:       1,
			MaxRetries:        1,
			BackoffSeconds:    1,
			RateLimit:         5,
			RateWindowSeconds: 10,
			DLQSize:           0, // ❌
		},
	}

	for _, cfg := range tests {
		if err := validateConfig(&cfg); err == nil {
			t.Errorf("expected validation error")
		}
	}
}