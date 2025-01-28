package fileserver_apilogger

import (
	"github.com/a179346/robert-go-monorepo/pkg/iologger"
	fileserver_config "github.com/a179346/robert-go-monorepo/services/fileserver/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GetApiLogger() *iologger.IoLogger {
	config := fileserver_config.GetLoggingConfig()
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
