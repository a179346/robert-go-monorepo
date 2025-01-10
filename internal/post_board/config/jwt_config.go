package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type JwtConfig struct {
	Secret        string
	ExpireSeconds int
}

func newJwtConfig() JwtConfig {
	secret := env_helper.GetStringEnv("JWT_SECRET", "my1-jwt2-3secret")
	expireSeconds := env_helper.GetIntEnv("JWT_EXPIRE_SECONDS", 3600)

	return JwtConfig{
		Secret:        secret,
		ExpireSeconds: expireSeconds,
	}
}
