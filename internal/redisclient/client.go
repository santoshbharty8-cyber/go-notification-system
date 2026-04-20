package redisclient

import (
	"context"
	"time"

	"go-notification-system/internal/logger"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	Ctx    = context.Background()
	Client *redis.Client
)

func InitRedis(redisURL string) {

	Client = redis.NewClient(&redis.Options{
		Addr:         redisURL,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		logger.Error("Redis connection failed",
			zap.String("url", redisURL),
			zap.Error(err),
		)
		panic(err)
	}

	logger.Info("Redis connected",
		zap.String("url", redisURL),
	)
}