package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type JwtConfig struct {
	Secret        []byte
	ExpireSeconds int
}

var jwtConfig JwtConfig

func init() {
	jwtConfig.Secret = []byte(env_helper.GetString("JWT_SECRET", "my1-jwt2-3secret"))
	jwtConfig.ExpireSeconds = env_helper.GetInt("JWT_EXPIRE_SECONDS", 3600)
}

func GetJwtConfig() JwtConfig {
	return jwtConfig
}
