package post_board_config

type Config struct {
	Server    ServerConfig
	DB        DBConfig
	Jwt       JwtConfig
	Migration MigrationConfig
}

func New() Config {
	return Config{
		Server:    newServerConfig(),
		DB:        newDBConfig(),
		Jwt:       newJwtConfig(),
		Migration: newMigrationConfig(),
	}
}
