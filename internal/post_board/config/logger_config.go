package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type LoggerConfig struct {
	Enable         bool
	TargetExchange string

	ConsumerSourceQueue string
	ConsumerConcurrency int
}

var loggerConfig LoggerConfig

func initLoggerConfig() {
	loggerConfig.Enable = env_helper.GetBool("LOGGER_ENABLE", true)
	loggerConfig.TargetExchange = env_helper.GetString("LOGGER_TARGET_EXCHANGE", "logging-exchange")

	loggerConfig.ConsumerSourceQueue = env_helper.GetString("LOGGER_CONSUMER_SOURCE_QUEUE", "logging-queue")
	loggerConfig.ConsumerConcurrency = env_helper.GetInt("LOGGER_CONSUMER_CONCURRENCY", 4)
}

func GetLoggerConfig() LoggerConfig {
	initAll()
	return loggerConfig
}
