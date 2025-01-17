package fileserver_config

import "sync"

var once sync.Once

func initAll() {
	once.Do(func() {
		initLoggerConfig()
		initServerConfig()
		initStorageConfig()
	})
}
