package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type ServerConfig struct {
	Port int
}

var serverConfig ServerConfig

func init() {
	serverConfig.Port = env_helper.GetInt("SERVER_PORT", 8082)
}

func GetServerConfig() ServerConfig {
	return serverConfig
}