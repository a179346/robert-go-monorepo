package delay_app_config

import "github.com/a179346/robert-go-monorepo/pkg/envhelper"

type ServerConfig struct {
	Port int
}

var serverConfig ServerConfig

func init() {
	serverConfig.Port = envhelper.GetInt("SERVER_PORT", 8080)
}

func GetServerConfig() ServerConfig {
	return serverConfig
}
