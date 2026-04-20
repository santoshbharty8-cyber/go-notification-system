package bootstrap

import (
	"log"
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

	// 1️⃣ Load config
	config.LoadConfig()

	// 2️⃣ Init logger
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

	// 3️⃣ Init Redis
	redisclient.InitRedis(config.AppConfig.RedisURL)

	// 4️⃣ Init DLQ
	queue.InitDLQ(config.AppConfig.DLQSize)

	// 5️⃣ Start workers
	startRedisWorkerPool(config.AppConfig.WorkerCount)

	// 6️⃣ Rate limiter
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

	// ✅ FIXED: health handler (errcheck + headers)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("OK")); err != nil {
			log.Println("failed to write health response:", err)
		}
	})

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  config.AppConfig.HTTPTimeout,
		WriteTimeout: config.AppConfig.HTTPTimeout,
	}

	logger.Info("Server starting",
		zap.String("addr", server.Addr),
	)

	
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server failed to start",
			zap.Error(err),
		)
	}
}

func startRedisWorkerPool(count int) {
	for i := 1; i <= count; i++ {
		go workers.StartRedisWorker(i)
	}
}