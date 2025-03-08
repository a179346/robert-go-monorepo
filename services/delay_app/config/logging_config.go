package delay_app_config

import "github.com/a179346/robert-go-monorepo/pkg/envhelper"

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
	loggingConfig.Enable = envhelper.GetBool("LOGGING_ENABLE", true)
	loggingConfig.Filename = envhelper.GetString("LOGGING_FILEANME", "./logs/delay_app/app/api.log")
	loggingConfig.MaxSizeMBs = envhelper.GetInt("LOGGING_MAX_SIZE_MBS", 50)
	loggingConfig.MaxBackups = envhelper.GetInt("LOGGING_MAX_BACKUPS", 3)
	loggingConfig.MaxAgeDays = envhelper.GetInt("LOGGING_MAX_AGE_DAYS", 30)
	loggingConfig.Compress = envhelper.GetBool("LOGGING_COMPRESS", false)
}

func GetLoggingConfig() LoggingConfig {
	return loggingConfig
}
