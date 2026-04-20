package config

import "testing"

func TestValidateConfigSuccess(t *testing.T) {
	cfg := &Config{
		RedisURL:          "localhost:6379",
		WorkerCount:       3,
		MaxRetries:        3,
		BackoffSeconds:    1,
		RateLimit:         5,
		RateWindowSeconds: 10,
		DLQSize:           100,
	}

	if err := validateConfig(cfg); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateConfigFailures(t *testing.T) {

	tests := []Config{
		{RedisURL: ""},
		{RedisURL: "x", WorkerCount: 0},
		{RedisURL: "x", WorkerCount: 1, MaxRetries: -1},
		{RedisURL: "x", WorkerCount: 1, MaxRetries: 1, BackoffSeconds: 0},
	}

	for _, cfg := range tests {
		if err := validateConfig(&cfg); err == nil {
			t.Errorf("expected error for %+v", cfg)
		}
	}
}