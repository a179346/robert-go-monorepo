package post_board_applogger

import (
	post_board_config "github.com/a179346/robert-go-monorepo/internal/post_board/config"
	"github.com/a179346/robert-go-monorepo/pkg/flushlogger"
	"gopkg.in/natefinch/lumberjack.v2"
)

func GetFlushLogger() *flushlogger.FlushLogger {
	config := post_board_config.GetLoggerConfig()
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
