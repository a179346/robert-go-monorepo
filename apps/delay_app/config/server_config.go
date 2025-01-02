package delay_app_config

import (
	"os"
	"strconv"
)

type ServerConfig struct {
	Port int
}

func newServerConfig() ServerConfig {
	port := 8080
	if p, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
		port = p
	}

	return ServerConfig{
		Port: port,
	}
}
