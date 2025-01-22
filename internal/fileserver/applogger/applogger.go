package fileserver_applogger

import (
	fileserver_config "github.com/a179346/robert-go-monorepo/internal/fileserver/config"
	"github.com/a179346/robert-go-monorepo/pkg/iologger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GetAppLogger() *iologger.IoLogger {
	config := fileserver_config.GetLoggerConfig()
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
