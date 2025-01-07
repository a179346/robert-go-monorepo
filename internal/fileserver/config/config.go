package fileserver_config

type Config struct {
	Server  ServerConfig
	Storage StorageConfig
}

func New() Config {
	return Config{
		Server:  newServerConfig(),
		Storage: newStorageConfig(),
	}
}
