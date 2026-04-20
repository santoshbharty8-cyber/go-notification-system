package bootstrap

import (
	"net/http"
	"time"

	"go-notification-system/internal/config"
	"go-notification-system/internal/handlers"
	"go-notification-system/internal/logger"
	"go-notification-system/internal/middleware"
	"go-notification-system/internal/queue"
	"go-notification-system/internal/ratelimiter"
	"go-notification-system/internal/redisclient"
	"go-notification-system/internal/workers"

	"go.uber.org/zap"
)

func StartServer() {

	config.LoadConfig()

	logger.InitLogger(
		config.AppConfig.AppEnv,
		config.AppConfig.LogLevel,
	)
	logger.Info("Config loaded",
		zap.String("env", config.AppConfig.AppEnv),
		zap.Int("workers", config.AppConfig.WorkerCount),
		zap.Int("retries", config.AppConfig.MaxRetries),
		zap.String("log_level", config.AppConfig.LogLevel),
	)

	redisclient.InitRedis(config.AppConfig.RedisURL)

	queue.InitDLQ(config.AppConfig.DLQSize)

	startRedisWorkerPool(config.AppConfig.WorkerCount)

	limiter := ratelimiter.NewLimiter(
		config.AppConfig.RateLimit,
		time.Duration(config.AppConfig.RateWindowSeconds)*time.Second,
	)

	// 7️⃣ Routes
	http.HandleFunc("/event",
		middleware.LoggingMiddleware(
			middleware.RateLimitMiddleware(limiter, handlers.EventHandler),
		),
	)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  config.AppConfig.HTTPTimeout,
		WriteTimeout: config.AppConfig.HTTPTimeout,
	}

	logger.Info("Server starting",
		zap.String("addr", server.Addr),
	)

	http.ListenAndServe(server.Addr, nil)
}

func startRedisWorkerPool(count int) {
	for i := 1; i <= count; i++ {
		go workers.StartRedisWorker(i)
	}
}