package delay_app_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type AppConfig struct {
	ID      string
	Version string
}

var appConfig AppConfig

func init() {
	appConfig.ID = env_helper.GetString("APP_ID", "delay_app")
	appConfig.Version = env_helper.GetString("APP_VERSION", "0.0.0")
}

func GetAppConfig() AppConfig {
	return appConfig
}
