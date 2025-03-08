package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/envhelper"

type DebugConfig struct {
	ResponseErrorDetail bool
}

var debugConfig DebugConfig

func init() {
	debugConfig.ResponseErrorDetail = envhelper.GetBool("DEBUG_RESPONSE_ERROR_DETAIL", true)
}

func GetDebugConfig() DebugConfig {
	return debugConfig
}
