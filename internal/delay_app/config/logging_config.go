package delay_app_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type LoggingConfig struct {
	Enable     bool
	Filename   string
	MaxSizeMBs int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
}

var loggingConfig LoggingConfig

func init() {
	loggingConfig.Enable = env_helper.GetBool("LOGGING_ENABLE", true)
	loggingConfig.Filename = env_helper.GetString("LOGGING_FILEANME", "./logs/delay_app/app/api.log")
	loggingConfig.MaxSizeMBs = env_helper.GetInt("LOGGING_MAX_SIZE_MBS", 50)
	loggingConfig.MaxBackups = env_helper.GetInt("LOGGING_MAX_BACKUPS", 3)
	loggingConfig.MaxAgeDays = env_helper.GetInt("LOGGING_MAX_AGE_DAYS", 30)
	loggingConfig.Compress = env_helper.GetBool("LOGGING_COMPRESS", false)
}

func GetLoggingConfig() LoggingConfig {
	return loggingConfig
}
