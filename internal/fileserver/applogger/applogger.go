package fileserver_applogger

import (
	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	"github.com/a179346/robert-go-monorepo/pkg/flushlogger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GetFlushLogger() *flushlogger.FlushLogger {
	config := fileserver_config.GetLoggerConfig()
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
