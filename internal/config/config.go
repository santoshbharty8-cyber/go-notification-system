package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	AppEnv   string
	LogLevel string

	RedisURL string

	WorkerCount int

	MaxRetries     int
	BackoffSeconds int

	RateLimit         int
	RateWindowSeconds int

	HTTPTimeout  time.Duration
	RedisTimeout time.Duration

	DLQSize int
}

// Global config instance
var AppConfig *Config

// LoadConfig initializes configuration
func LoadConfig() {

	env := getEnv("APP_ENV", "dev")
	envFile := ".env." + env

	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			log.Printf("⚠️ Failed to load %s\n", envFile)
		} else {
			log.Printf("✅ Loaded config from %s\n", envFile)
		}
	} else {
		log.Printf("ℹ️ %s not found, using system env\n", envFile)
	}

	AppConfig = &Config{
		AppEnv:   env,
		LogLevel: getEnv("LOG_LEVEL", "info"),

		RedisURL: getEnv("REDIS_URL", "localhost:6379"),

		WorkerCount: getEnvAsInt("WORKER_COUNT", 3),

		MaxRetries:     getEnvAsInt("MAX_RETRIES", 3),
		BackoffSeconds: getEnvAsInt("BACKOFF_SECONDS", 1),

		RateLimit:         getEnvAsInt("RATE_LIMIT", 5),
		RateWindowSeconds: getEnvAsInt("RATE_WINDOW_SECONDS", 10),

		HTTPTimeout:  time.Duration(getEnvAsInt("HTTP_TIMEOUT", 5)) * time.Second,
		RedisTimeout: time.Duration(getEnvAsInt("REDIS_TIMEOUT", 2)) * time.Second,

		DLQSize: getEnvAsInt("DLQ_SIZE", 100),
	}

	// 🔥 VALIDATE (now testable)
	if err := validateConfig(AppConfig); err != nil {
		log.Fatal(err)
	}

	log.Printf(
		"🚀 Config Loaded: ENV=%s | Workers=%d | Retries=%d | LogLevel=%s",
		AppConfig.AppEnv,
		AppConfig.WorkerCount,
		AppConfig.MaxRetries,
		AppConfig.LogLevel,
	)
}

func getEnv(key string, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("⚠️ Invalid value for %s, using default\n", key)
		return defaultVal
	}
	return val
}

// 🔥 NOW RETURNS ERROR (testable)
func validateConfig(c *Config) error {

	if c.RedisURL == "" {
		return errors.New("REDIS_URL is required")
	}

	if c.WorkerCount <= 0 {
		return errors.New("WORKER_COUNT must be > 0")
	}

	if c.MaxRetries < 0 {
		return errors.New("MAX_RETRIES cannot be negative")
	}

	if c.BackoffSeconds <= 0 {
		return errors.New("BACKOFF_SECONDS must be > 0")
	}

	if c.RateLimit <= 0 {
		return errors.New("RATE_LIMIT must be > 0")
	}

	if c.RateWindowSeconds <= 0 {
		return errors.New("RATE_WINDOW_SECONDS must be > 0")
	}

	if c.DLQSize <= 0 {
		return errors.New("DLQ_SIZE must be > 0")
	}

	return nil
}