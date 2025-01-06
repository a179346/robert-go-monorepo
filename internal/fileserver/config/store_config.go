package fileserver_config

import (
	"os"
)

type StoreConfig struct {
	RootPath string
}

func newStoreConfig() StoreConfig {
	rootPath := os.Getenv("STORE_ROOT_PATH")
	if rootPath == "" {
		rootPath = "./storage/fileserver/store"
	}

	return StoreConfig{
		RootPath: rootPath,
	}
}
