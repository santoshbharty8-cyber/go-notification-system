package helpers

import (
	"go-notification-system/internal/config"
	"go-notification-system/internal/logger"
)

func InitTestEnv() {
	config.AppConfig = &config.Config{
		AppEnv:   "test",
		LogLevel: "error",
	}
	logger.InitLogger("test", "error")
}