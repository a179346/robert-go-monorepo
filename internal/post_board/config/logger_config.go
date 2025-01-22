package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type LoggerConfig struct {
	Enable         bool
	TargetExchange string
}

var loggerConfig LoggerConfig

func initLoggerConfig() {
	loggerConfig.Enable = env_helper.GetBool("LOGGER_ENABLE", true)
	loggerConfig.TargetExchange = env_helper.GetString("LOGGER_TARGET_EXCHANGE", "logging-exchange")
}

func GetLoggerConfig() LoggerConfig {
	initAll()
	return loggerConfig
}
