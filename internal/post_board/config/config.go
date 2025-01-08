package post_board_config

type Config struct {
	Server    ServerConfig
	DB        DBConfig
	Migration MigrationConfig
}

func New() Config {
	return Config{
		Server:    newServerConfig(),
		DB:        newDBConfig(),
		Migration: newMigrationConfig(),
	}
}
