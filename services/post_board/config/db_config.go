package post_board_config

import "github.com/a179346/robert-go-monorepo/pkg/envhelper"

type DBConfig struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

var dbConfig DBConfig

func init() {
	dbConfig.Host = envhelper.GetString("DB_HOST", "localhost")
	dbConfig.Port = envhelper.GetInt("DB_PORT", 5432)
	dbConfig.Database = envhelper.GetString("DB_DATABASE", "post-board")
	dbConfig.User = envhelper.GetString("DB_USER", "post-board-user")
	dbConfig.Password = envhelper.GetString("DB_PASSWORD", "mysecretpassword")
}

func GetDBConfig() DBConfig {
	return dbConfig
}
