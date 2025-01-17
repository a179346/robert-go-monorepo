package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type LoggerConfig struct {
	Enable     bool
	Filename   string
	MaxSizeMBs int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
}

var loggerConfig LoggerConfig

func initLoggerConfig() {
	loggerConfig.Enable = env_helper.GetBoolEnv("LOGGER_ENABLE", true)
	loggerConfig.Filename = env_helper.GetStringEnv("LOGGER_FILEANME", "./logs/post_board/app/requests.log")
	loggerConfig.MaxSizeMBs = env_helper.GetIntEnv("LOGGER_MAX_SIZE_MBS", 50)
	loggerConfig.MaxBackups = env_helper.GetIntEnv("LOGGER_MAX_BACKUPS", 3)
	loggerConfig.MaxAgeDays = env_helper.GetIntEnv("LOGGER_MAX_AGE_DAYS", 30)
	loggerConfig.Compress = env_helper.GetBoolEnv("LOGGER_COMPRESS", false)
}

func GetLoggerConfig() LoggerConfig {
	initAll()
	return loggerConfig
}
