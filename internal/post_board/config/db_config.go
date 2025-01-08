package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type DBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func newDBConfig() DBConfig {
	host := env_helper.GetStringEnv("DB_HOST", "localhost")
	port := env_helper.GetIntEnv("DB_PORT", 5432)
	database := env_helper.GetStringEnv("DB_DATABASE", "post-board")
	user := env_helper.GetStringEnv("DB_USER", "post-board-user")
	password := env_helper.GetStringEnv("DB_PASSWORD", "mysecretpassword")

	return DBConfig{
		Host:     host,
		Port:     port,
		Database: database,
		User:     user,
		Password: password,
	}
}
