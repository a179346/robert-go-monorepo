package post_board_config

import (
	"sync"
)

var once sync.Once

func initAll() {
	once.Do(func() {
		initDBConfig()
		initDebugConfig()
		initJwtConfig()
		initLoggerConfig()
		initMigrationConfig()
		initServerConfig()
	})
}
