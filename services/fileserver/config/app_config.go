package fileserver_config

import "github.com/a179346/robert-go-monorepo/pkg/envhelper"

type AppConfig struct {
	ID      string
	Version string
}

var appConfig AppConfig

func init() {
	appConfig.ID = envhelper.GetString("APP_ID", "fileserver")
	appConfig.Version = envhelper.GetString("APP_VERSION", "0.0.0")
}

func GetAppConfig() AppConfig {
	return appConfig
}
