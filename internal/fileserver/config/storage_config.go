package fileserver_config

import (
	"os"
	"path/filepath"
)

type StorageConfig struct {
	RootPath      string
	StoreRootPath string
}

func newStorageConfig() StorageConfig {
	rootPath := os.Getenv("STORAGE_ROOT_PATH")
	if rootPath == "" {
		rootPath = "./storage/fileserver"
	}

	storeRootPath := filepath.Join(rootPath, "store")

	return StorageConfig{
		RootPath:      rootPath,
		StoreRootPath: storeRootPath,
	}
}
