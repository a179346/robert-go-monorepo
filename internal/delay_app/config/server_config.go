package delay_app_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type ServerConfig struct {
	Port int
}

var serverConfig ServerConfig

func initServerConfig() {
	serverConfig.Port = env_helper.GetIntEnv("SERVER_PORT", 8080)
}

func GetServerConfig() ServerConfig {
	initAll()
	return serverConfig
}
