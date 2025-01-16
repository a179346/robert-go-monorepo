package delay_app_config

import "sync"

var once sync.Once

func initAll() {
	once.Do(func() {
		initServerConfig()
	})
}
