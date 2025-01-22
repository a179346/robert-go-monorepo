package delay_app_applogger

import (
	delay_app_config "github.com/a179346/robert-go-monorepo/internal/delay_app/config"
	"github.com/a179346/robert-go-monorepo/pkg/iologger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GetAppLogger() *iologger.IoLogger {
	config := delay_app_config.GetLoggerConfig()
	if !config.Enable {
		return nil
	}

	return iologger.New(&lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSizeMBs,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAgeDays,
		Compress:   config.Compress,
	})
}
