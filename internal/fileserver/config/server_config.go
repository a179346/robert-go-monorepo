package fileserver_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type ServerConfig struct {
	isInited bool
	Port     int
}

var serverConfig ServerConfig

func initServerConfig() {
	if serverConfig.isInited {
		return
	}
	serverConfig.Port = env_helper.GetIntEnv("SERVER_PORT", 8081)

	serverConfig.isInited = true
}

func GetServerConfig() ServerConfig {
	initServerConfig()
	return serverConfig
}
