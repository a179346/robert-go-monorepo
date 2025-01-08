package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type ServerConfig struct {
	Port int
}

func newServerConfig() ServerConfig {
	port := env_helper.GetIntEnv("SERVER_PORT", 8082)

	return ServerConfig{
		Port: port,
	}
}
