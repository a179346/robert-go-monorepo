package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/env_helper"

type DBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

var dbConfig DBConfig

func init() {
	dbConfig.Host = env_helper.GetString("DB_HOST", "localhost")
	dbConfig.Port = env_helper.GetInt("DB_PORT", 5432)
	dbConfig.Database = env_helper.GetString("DB_DATABASE", "post-board")
	dbConfig.User = env_helper.GetString("DB_USER", "post-board-user")
	dbConfig.Password = env_helper.GetString("DB_PASSWORD", "mysecretpassword")
}

func GetDBConfig() DBConfig {
	return dbConfig
}
