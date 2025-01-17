package delay_app_applogger

import (
	delay_app_config "github.com/a179346/robert-go-monorepo/internal/delay_app/config"
	"github.com/a179346/robert-go-monorepo/pkg/flushlogger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GetFlushLogger() *flushlogger.FlushLogger {
	config := delay_app_config.GetLoggerConfig()
	if !config.Enable {
		return nil
	}

	return flushlogger.New(&lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSizeMBs,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAgeDays,
		Compress:   config.Compress,
	})
}
