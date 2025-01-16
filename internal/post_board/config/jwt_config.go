package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type JwtConfig struct {
	Secret        []byte
	ExpireSeconds int
}

var jwtConfig JwtConfig

func initJwtConfig() {
	jwtConfig.Secret = []byte(env_helper.GetStringEnv("JWT_SECRET", "my1-jwt2-3secret"))
	jwtConfig.ExpireSeconds = env_helper.GetIntEnv("JWT_EXPIRE_SECONDS", 3600)
}

func GetJwtConfig() JwtConfig {
	initAll()
	return jwtConfig
}
