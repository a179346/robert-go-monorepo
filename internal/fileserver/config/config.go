package fileserver_config

type Config struct {
	Server ServerConfig
	Store  StoreConfig
}

func New() Config {
	return Config{
		Server: newServerConfig(),
		Store:  newStoreConfig(),
	}
}
