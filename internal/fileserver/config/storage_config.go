package fileserver_config

import (
	"path/filepath"

	"github.com/a179346/robert-go-monorepo/pkg/env_helper"
)

type StorageConfig struct {
	RootPath      string
	StoreRootPath string
}

func newStorageConfig() StorageConfig {
	rootPath := env_helper.GetStringEnv("STORAGE_ROOT_PATH", "./storage/fileserver")

	storeRootPath := filepath.Join(rootPath, "store")

	return StorageConfig{
		RootPath:      rootPath,
		StoreRootPath: storeRootPath,
	}
}
