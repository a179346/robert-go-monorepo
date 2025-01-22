package delay_app_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type DebugConfig struct {
	ResponseErrorDetail bool
}

var debugConfig DebugConfig

func initDebugConfig() {
	debugConfig.ResponseErrorDetail = env_helper.GetBool("DEBUG_RESPONSE_ERROR_DETAIL", true)
}

func GetDebugConfig() DebugConfig {
	initAll()
	return debugConfig
}
