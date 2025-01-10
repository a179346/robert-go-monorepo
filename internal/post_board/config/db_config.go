package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type DBConfig struct {
	isInited bool
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

var dbConfig DBConfig

func initDBConfig() {
	if dbConfig.isInited {
		return
	}
	dbConfig.Host = env_helper.GetStringEnv("DB_HOST", "localhost")
	dbConfig.Port = env_helper.GetIntEnv("DB_PORT", 5432)
	dbConfig.Database = env_helper.GetStringEnv("DB_DATABASE", "post-board")
	dbConfig.User = env_helper.GetStringEnv("DB_USER", "post-board-user")
	dbConfig.Password = env_helper.GetStringEnv("DB_PASSWORD", "mysecretpassword")

	dbConfig.isInited = true
}

func GetDBConfig() DBConfig {
	initDBConfig()
	return dbConfig
}
