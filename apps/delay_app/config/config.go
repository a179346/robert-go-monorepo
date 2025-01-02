package delay_app_config

type Config struct {
	Server ServerConfig
}

func New() Config {
	return Config{
		Server: newServerConfig(),
	}
}
